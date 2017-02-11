package udp

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/toefel18/p2pnet/discovery/network"
)

const (
	DiscoveryPort               = 6667
	staleAfterDuration          = 15 * time.Second
	removeAfterInactiveDuration = 62 * time.Second
	pruneInterval               = 5 * time.Second
)

type Peers map[string]*Peer

type Peer struct {
	PeerName    string
	NetProtocol string
	NetAddress  string
	NetPort     int
	AliveSince  time.Time
	KnownSince  time.Time

	LastSeen  time.Time
	TimesSeen int64
	Status    string
}

func (p *Peer) Name() string {
	return p.PeerName
}

func (p *Peer) Address() string {
	return fmt.Sprintf("%v://%v:%v", p.NetProtocol, p.NetAddress, p.NetPort)
}

type Discovery struct {
	heartbeatInterval    time.Duration
	peersByName          Peers
	peerLock             sync.Mutex
	closeHeartbeatSender chan struct{}
	closeOutdatedPrune   chan struct{}

	*heartbeater
	*heartbeatListener
}

func NewDefaultDiscovey(localName string) *Discovery {
	return &Discovery{
		heartbeatInterval:    3 * time.Second,
		peersByName:          make(Peers),
		closeHeartbeatSender: make(chan struct{}),
		heartbeater: &heartbeater{
			address:    "255.255.255.255",
			port:       DiscoveryPort,
			name:       localName,
			aliveSince: time.Now(),
		},
		heartbeatListener: &heartbeatListener{
			address:     "0.0.0.0",
			port:        DiscoveryPort,
			onHeartbeat: make(chan *heartbeatPacket, 10),
		},
	}
}

func (d *Discovery) Start() {
	// start with initial heartbeat to check if everything works before entering forever loop.
	d.SendHeartbeat()
	go d.Listen()
	go d.sendHeartbeats()
	go d.receiveHeartbeats()
	go d.pruneOutdatedPeers()
}

func (d *Discovery) Stop() {
	d.closeHeartbeatSender <- struct{}{}
	d.closeOutdatedPrune <- struct{}{}
	d.heartbeatListener.Close()
}

func (d *Discovery) removePeer(name string) {
	d.peerLock.Lock()
	defer d.peerLock.Unlock()
	delete(d.peersByName, name)
}

func (d *Discovery) addOrUpdatePeer(peer *Peer) {
	d.peerLock.Lock()
	defer d.peerLock.Unlock()
	if knownPeer, exists := d.peersByName[peer.PeerName]; exists {
		peer.KnownSince = knownPeer.KnownSince
		peer.TimesSeen += knownPeer.TimesSeen
		if peer.TimesSeen > 10 {
			peer.Status = "stable"
		}
	}
	d.peersByName[peer.PeerName] = peer
}

func (d *Discovery) Peers() map[string]network.Peer {
	d.peerLock.Lock()
	defer d.peerLock.Unlock()
	copy := make(map[string]network.Peer)
	for name, peer := range d.peersByName {
		copyOfPeer := *peer
		copy[name] = &copyOfPeer
	}
	return copy
}

func (d *Discovery) sendHeartbeats() {
	for {
		select {
		case <-d.closeHeartbeatSender:
			return
		case <-time.Tick(d.heartbeatInterval):
			d.SendHeartbeat()
		}
	}
}

func (d *Discovery) receiveHeartbeats() {
	for {
		if hb, ok := <-d.heartbeatListener.onHeartbeat; !ok {
			return
		} else {
			if peer, err := toPeer(hb); err != nil {
				log.Printf("Error receiving heartbeat from peer %v", err.Error())
			} else {
				d.addOrUpdatePeer(peer)
			}
		}
	}
}

func (d *Discovery) pruneOutdatedPeers() {
	for {
		select {
		case <-d.closeHeartbeatSender:
			return
		case <-time.Tick(pruneInterval):
			d.doPruneOutdatedPeers()
		}
	}
}

func (d *Discovery) doPruneOutdatedPeers() {
	d.peerLock.Lock()
	defer d.peerLock.Unlock()
	for name, peer := range d.peersByName {
		secondsSinceLastSeen := time.Now().Sub(peer.LastSeen)
		if secondsSinceLastSeen > removeAfterInactiveDuration {
			delete(d.peersByName, name)
		} else if secondsSinceLastSeen > staleAfterDuration {
			peer.Status = "stale, removing in " + strconv.Itoa(int((removeAfterInactiveDuration - secondsSinceLastSeen).Seconds())) + "s"
		}
	}
}

func toPeer(hb *heartbeatPacket) (*Peer, error) {
	return &Peer{
		PeerName:    hb.Name,
		NetProtocol: "udp",
		NetAddress:  hb.Address,
		NetPort:     hb.Port,
		AliveSince:  hb.AliveSince,
		KnownSince:  time.Now(),
		LastSeen:    time.Now(),
		TimesSeen:   1,
		Status:      "new",
	}, nil
}
