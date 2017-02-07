package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"

	"time"

	"github.com/toefel18/p2pnet/node"
)

func main() {
	log.Println("staring p2pnet...")
	user := os.Getenv("USER")
	if len(user) == 0 {
		user = "uknwn" + strconv.Itoa(rand.Int())
	}
	log.Println("starting broadcaster")
	broadcaster := node.NewUDPBroadcaster(user, time.Now())
	go broadcaster.BroadcastLiveness()
	discoverer := node.NewUDPDiscoverer()
	discoverer.Listen()
}
