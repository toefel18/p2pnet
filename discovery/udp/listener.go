package udp

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"strings"

	"github.com/jbenet/go-reuseport"
)

type HeartbeatListener interface {
	Listen()
}

type heartbeatListener struct {
	address     string
	port        int
	onHeartbeat chan *heartbeatPacket

	conn net.PacketConn
}

func NewHeartbeatListener(address string, port int, onHeartbeat chan *heartbeatPacket) *heartbeatListener {
	return &heartbeatListener{address, port, onHeartbeat, nil}
}

func (d *heartbeatListener) Close() {
	if d.conn != nil {
		d.conn.Close()
		close(d.onHeartbeat)
	}
}

func (d *heartbeatListener) Listen() {
	conn, err := reuseport.ListenPacket("udp", d.address+":"+strconv.Itoa(d.port))
	if err != nil {
		panic(fmt.Errorf("Failed to listen %v", err.Error()))
	}
	defer conn.Close()
	log.Println("Listening on", conn.LocalAddr())
	d.receiveForever(conn)
}

func (d *heartbeatListener) receiveForever(conn net.PacketConn) *net.UDPAddr {
	buffer := make([]byte, 1024)
	for {
		if n, sender, err := conn.ReadFrom(buffer); err == nil {
			packet := string(buffer[:n])
			hb := &heartbeatPacket{}
			if err := hb.UnmarshalText(buffer[:n]); err != nil {
				log.Printf("ERROR: Malformed heartbeat from %v: %v\n", sender.String(), packet)
			} else {
				addressOnPacket := sender.String()
				hb.Address = addressOnPacket[:strings.Index(addressOnPacket, ":")]
				select {
				case d.onHeartbeat <- hb:
				default:
					log.Printf("onHeartbeat channel full, discarding heartbeat %v", hb.String())
				}
			}
		} else {
			log.Printf("Error while reading from UDP: %v", err.Error())
		}
	}
}
