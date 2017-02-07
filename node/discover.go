package node

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const DefaultHeartbeatPort = 6667

type heartbeatPacket struct {
	Name       string
	Address    string
	AliveSince time.Time
}

func (h *heartbeatPacket) String() string {
	return fmt.Sprintf("%v~%v~%v", h.Name, h.Address, h.AliveSince.String())
}

func (h *heartbeatPacket) MarshallText() (text []byte, err error) {
	return []byte(fmt.Sprintf("%v~%v~%d", h.Name, h.Address, h.AliveSince.Unix())), nil
}

func (h *heartbeatPacket) UnmarshalText(text []byte) error {
	packet := string(text)
	firstFieldEnd := strings.Index(packet, "~")
	lastFieldStart := strings.LastIndex(packet, "~")
	if firstFieldEnd == lastFieldStart || firstFieldEnd == 0 {
		return fmt.Errorf("Unknown heartbeat packet: %v", packet)
	}
	var err error
	h.Name = packet[:firstFieldEnd]
	h.Address = packet[firstFieldEnd+1 : lastFieldStart]
	h.AliveSince, err = millisToTime(packet[lastFieldStart+1:])
	return err
}

type Discoverer interface {
	Listen()
}

type UDPDiscoverer struct {
	address string
	port    int
}

func NewUDPDiscoverer() Discoverer {
	return &UDPDiscoverer{"", DefaultHeartbeatPort}
}

func (d *UDPDiscoverer) Listen() {
	conn, err := net.ListenUDP("udp", resolveAddress(d.address, d.port))
	file, err := conn.File()
	if err != nil { panic(err)}
	syscall.SetsockoptInt(int(file.Fd()), syscall.IPPROTO_IP, syscall.SO_REUSEADDR, 1)
	if err != nil {
		panic(fmt.Errorf("Failed to listen %v", err.Error()))
	}
	defer conn.Close()
	log.Println("Listening on", conn.LocalAddr())
	d.receiveForever(conn)
}

func (d *UDPDiscoverer) receiveForever(conn *net.UDPConn) *net.UDPAddr {
	buffer := make([]byte, 1024)
	for {
		if n, sender, err := conn.ReadFromUDP(buffer); err == nil {
			packet := string(buffer[:n])
			hb := &heartbeatPacket{}
			if err := hb.UnmarshalText(buffer[:n]); err != nil {
				log.Printf("Malformed heartbeat from %v:%v %v\n", sender.IP, sender.Port, packet)
			} else {
				log.Printf("Heartbeat received from %v:%v %v\n", sender.IP, sender.Port, hb.String())
			}
		} else {
			log.Printf("Error while reading from UDP: %v", err.Error())
		}
	}
}

type HeartbeatBroadcaster interface {
	BroadcastLiveness()
}

type UDPBroadcaster struct {
	address    string
	port       int
	name       string
	aliveSince time.Time
}

func NewUDPBroadcaster(name string, aliveSince time.Time) HeartbeatBroadcaster {
	return &UDPBroadcaster{net.IPv4bcast.String(), DefaultHeartbeatPort, name, aliveSince}
}

func (b *UDPBroadcaster) BroadcastLiveness() {
	for {
		b.broadcastLiveness()
		time.Sleep(5 * time.Second)
	}
}

func (b *UDPBroadcaster) broadcastLiveness() {
	remote := resolveAddress(b.address, b.port)
	local := resolveAddress("127.0.0.1", DefaultHeartbeatPort)
	conn, err := net.DialUDP("udp", local, remote)
	if err != nil {
		panic(fmt.Errorf("failed to send heartbeat to %v:%v %v", remote.IP, remote.Port, err))
	}
	defer conn.Close()
	hb := heartbeatPacket{
		Name:       b.name,
		Address:    local.IP.String(),
		AliveSince: b.aliveSince}
	marshalledHeartbeat, _ := hb.MarshallText()
	conn.WriteToUDP(marshalledHeartbeat, remote)
}

func millisToTime(secsSinceEpochString string) (time.Time, error) {
	if secsSinceEpoch, err := strconv.ParseInt(secsSinceEpochString, 10, 64); err != nil {
		return time.Now(), err
	} else {
		return time.Unix(secsSinceEpoch, 0), nil
	}
}

func resolveAddress(address string, port int) *net.UDPAddr {
	listenInterface := address + ":" + strconv.Itoa(port)
	if address, err := net.ResolveUDPAddr("udp", listenInterface); err == nil {
		return address
	} else {
		panic(fmt.Errorf("Could resolve address %v, error: %v", listenInterface, err.Error()))
	}
}
