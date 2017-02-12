package report

import (
	"fmt"
	"os"
	"os/exec"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/toefel18/p2pnet/discovery/udp"
)

const discoverySummaryFormat = `Discovery Summary - nodes in network
NAME	ADDRESS	ALIVE SINCE	KNOWN SINCE	LAST SEEN	TIMES SEEN	STATUS
{{ range $key, $p := . }}{{$p.PeerName}}	{{$p.NetProtocol}}://{{$p.NetAddress}}:{{$p.NetPort}}	{{$p.AliveSince}}	{{$p.KnownSince}}	{{$p.LastSeen}}	{{$p.TimesSeen}}	{{$p.Status}}
{{ end }}`

var (
	summaryTemplate *template.Template = nil
	writer          *tabwriter.Writer
)

func init() {
	var err error
	summaryTemplate, err = template.New("DiscoverySummary").Parse(discoverySummaryFormat)
	if err != nil {
		panic(err)
	}
	writer = tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
}

func PrintDiscoverySummaryContinuously(nodeName string, discovery *udp.Discovery, stopSignal chan struct{}) {
	for {
		select {
		case <-stopSignal:
			return
		case <-time.Tick(2 * time.Second):
			PrintDiscoverySummary(nodeName, discovery)
		}
	}
}

func clearStdOut() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func PrintDiscoverySummary(nodeName string, discovery *udp.Discovery) {
	if len(discovery.Peers()) > 0 {
		clearStdOut()
		fmt.Println(nodeName, ": Heartbeats sent:", discovery.HeartbeatsSent)
		summaryTemplate.Execute(writer, discovery.Peers())
		writer.Flush()
	}
}
