package udp

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jbenet/go-reuseport"
)

type Heartbeater interface {
	SendHeartbeat()
}

type heartbeater struct {
	address        string
	port           int
	name           string
	aliveSince     time.Time
	HeartbeatsSent int
}

func NewHeartbeater(name string, aliveSince time.Time, address string, port int) Heartbeater {
	return &heartbeater{address, port, name, aliveSince, 0}
}

func (b *heartbeater) SendHeartbeat() {
	remote := b.address + ":" + strconv.Itoa(b.port)
	conn, err := reuseport.Dial("udp", "", remote)
	if err != nil {
		panic(fmt.Errorf("failed to send heartbeat to %v: %v", remote, err))
	}
	defer conn.Close()
	hb := heartbeatPacket{
		Name:       b.name,
		Address:    "0.0.0.0:", // address will be on Packet, return port is important!
		Port:       b.port,
		AliveSince: b.aliveSince}
	marshalledHeartbeat, _ := hb.MarshallText()
	if _, err := conn.Write(marshalledHeartbeat); err != nil {
		log.Printf("Error sending heartbeat %v", err.Error())
	} else {
		b.HeartbeatsSent++
	}
}
