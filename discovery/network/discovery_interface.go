package network

type Discovery interface {
	Start()
	Stop()
	Peers() map[string]Peer
}

type Peer interface {
	Name() string
	Address() string
}
