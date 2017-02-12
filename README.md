# Network discovery

Simple UDP network discovery using network broadcasts. 

Tested on Linux and Mac

Get the sources and compile the binary:

    go get -u github.com/toefel18/p2pnet
    
Run the binary: 

    $GOPATH/bin/p2pnet -name ny-name
    
(you can run the same binary once with different names!) example:

```
greekgod : Heartbeats sent: 53
Discovery Summary - nodes in network
NAME                   ADDRESS                      ALIVE SINCE                      KNOWN SINCE                                LAST SEEN                                  TIMES SEEN    STATUS
greekgod               udp://192.168.178.18:6667    2017-02-12 11:45:08 +0100 CET    2017-02-12 11:45:08.670568904 +0100 CET    2017-02-12 11:46:50.679334059 +0100 CET    35            stable
bleetbot               udp://192.168.178.18:6667    2017-02-12 11:46:09 +0100 CET    2017-02-12 11:46:09.684798582 +0100 CET    2017-02-12 11:46:51.688555915 +0100 CET    15            stable
chef-special           udp://192.168.178.18:6667    2017-02-12 11:44:24 +0100 CET    2017-02-12 11:44:24.815930097 +0100 CET    2017-02-12 11:46:51.827372961 +0100 CET    50            stable
pietbrown              udp://192.168.178.18:6667    2017-02-12 11:45:04 +0100 CET    2017-02-12 11:45:04.223083933 +0100 CET    2017-02-12 11:46:49.232379504 +0100 CET    36            stable
keesklaas              udp://192.168.178.18:6667    2017-02-12 11:44:48 +0100 CET    2017-02-12 11:44:48.351150009 +0100 CET    2017-02-12 11:46:51.361892988 +0100 CET    42            stable
sailor                 udp://192.168.178.18:6667    2017-02-12 11:46:14 +0100 CET    2017-02-12 11:46:14.460671405 +0100 CET    2017-02-12 11:46:50.464255608 +0100 CET    13            stable
ronny                  udp://192.168.178.18:6667    2017-02-12 11:46:28 +0100 CET    2017-02-12 11:46:28.332264579 +0100 CET    2017-02-12 11:46:49.334184117 +0100 CET    8             new
spawnbot               udp://192.168.178.18:6667    2017-02-12 11:46:04 +0100 CET    2017-02-12 11:46:04.404727343 +0100 CET    2017-02-12 11:46:49.40853099 +0100 CET     16            stable

```