package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/toefel18/p2pnet/discovery/udp"
)

// TODO WARNING, vendored code has been changed to also set SO_BROADCAST (required for sending and listening to broadcast address)
// TODO files changed in vendor: const_bsd.go const_linux.go impl_unix.go
const discoverySummaryFormat = `Discovery Summary - nodes in network
NAME	ADDRESS	ALIVE SINCE	KNOWN SINCE	LAST SEEN	TIMES SEEN	STATUS
{{ range $key, $p := . }}{{$p.PeerName}}	{{$p.NetProtocol}}://{{$p.NetAddress}}:{{$p.NetPort}}	{{$p.AliveSince}}	{{$p.KnownSince}}	{{$p.LastSeen}}	{{$p.TimesSeen}}	{{$p.Status}}
{{ end }}`

func main() {
	log.Println("staring p2pnet...")
	nodeName, err := os.Hostname()
	if err != nil {
		nodeName = os.Getenv("USER")
	}
	if len(nodeName) == 0 {
		nodeName = "unknown-" + strconv.Itoa(rand.Int())
	}
	discovery := udp.NewDefaultDiscovey(nodeName)
	discovery.Start()
	fmt.Println("Waiting for hearbeats...")
	reportDiscoveryContinuously(discovery)
}

func reportDiscoveryContinuously(discovery *udp.Discovery) {
	summaryTemplate, err := template.New("DiscoverySummary").Parse(discoverySummaryFormat)
	if err != nil {
		panic(err)
	}
	w := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
	for {
		time.Sleep(2 * time.Second)
		if len(discovery.Peers()) == 0 {
			continue
		}
		clearStdOut()
		fmt.Println("Heartbeats sent:", discovery.HeartbeatsSent)
		summaryTemplate.Execute(w, discovery.Peers())
		w.Flush()
	}
}

func clearStdOut() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}
