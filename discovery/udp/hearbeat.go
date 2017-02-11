package udp

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ErrInvalidAddress struct {
	ErrMsg  string
	Address string
}

func (e ErrInvalidAddress) Error() string {
	return fmt.Sprintf("invalid ip:port adddress (got %v) %v", e.Address, e.ErrMsg)
}

type heartbeatPacket struct {
	Name       string
	Address    string
	Port       int
	AliveSince time.Time
}

func (h *heartbeatPacket) String() string {
	return fmt.Sprintf("%v~%v~%d~%v", h.Name, h.Address, h.Port, h.AliveSince.String())
}

func (h *heartbeatPacket) MarshallText() (text []byte, err error) {
	return []byte(fmt.Sprintf("%v~%v~%d~%d", h.Name, h.Address, h.Port, h.AliveSince.Unix())), nil
}

func (h *heartbeatPacket) UnmarshalText(text []byte) error {
	packet := string(text)
	fields := strings.Split(packet, "~")
	if len(fields) != 4 {
		return fmt.Errorf("expecting 4 fields but got %d in packet: %v", len(fields), packet)
	}
	var portErr, aliveSinceErr error
	h.Name = fields[0]
	h.Address = fields[1]
	h.Port, portErr = strconv.Atoi(fields[2])
	h.AliveSince, aliveSinceErr = hearbeatMillisToTime(fields[3])
	if portErr != nil || aliveSinceErr != nil {
		return fmt.Errorf("Error(s): %v %v", portErr, aliveSinceErr)
	}
	return nil
}

func hearbeatMillisToTime(secsSinceEpochString string) (time.Time, error) {
	if secsSinceEpoch, err := strconv.ParseInt(secsSinceEpochString, 10, 64); err != nil {
		return time.Now(), err
	} else {
		return time.Unix(secsSinceEpoch, 0), nil
	}
}
