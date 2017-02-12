package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"

	"github.com/toefel18/p2pnet/discovery/network"
	"github.com/toefel18/p2pnet/discovery/report"
	"github.com/toefel18/p2pnet/discovery/udp"
)

// TODO WARNING, vendored code has been changed to also set SO_BROADCAST (required for sending to broadcast address)
// TODO files changed in vendor: const_bsd.go const_linux.go impl_unix.go

var nodeName string

func init() {
	name, err := os.Hostname()
	if err != nil {
		name = os.Getenv("USER")
	}
	if len(name) == 0 {
		name = "unknown-" + strconv.Itoa(rand.Int())
	}
	flag.StringVar(&nodeName, "name", name, "the name to advertise to others (defaults to hostname)")
}

func main() {
	flag.Parse()
	log.Println("staring p2pnet using name: ", nodeName)
	discovery := udp.NewDefaultDiscovey(nodeName)
	discovery.Start()
	fmt.Println("Waiting for hearbeats...")
	stopSignal := make(chan struct{})
	go report.PrintDiscoverySummaryContinuously(nodeName, discovery, stopSignal)
	waitForSigInterrupt(discovery, stopSignal)
}

func waitForSigInterrupt(discovery network.Discovery, stopSignal chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	fmt.Println("Ctrl+c to stop")
	signal := <-c
	fmt.Println("got signal", signal.String(), "shutting down ...")
	discovery.Stop()
	fmt.Println("stopping reporting")
	stopSignal <- struct{}{}
	fmt.Println("bye")
}
