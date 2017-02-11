# Network discovery

Simple UDP network discovery using network broadcasts. 

Tested on Linux and Mac

Get the sources and compile the binary:

    go get -u github.com/toefel18/p2pnet
    
Run the binary: 

    $GOPATH/bin/p2pnet
    
Expected output:

```text
Heartbeats sent: 20
Discovery Summary - nodes in network
NAME                   ADDRESS                      ALIVE SINCE                      KNOWN SINCE                                LAST SEEN                                  TIMES SEEN    STATUS
memyself-ubuntu        udp://192.168.178.18:6667    2017-02-11 19:05:43 +0100 CET    2017-02-11 19:05:46.957899889 +0100 CET    2017-02-11 19:06:40.961708284 +0100 CET    19            stable

```