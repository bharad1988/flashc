//This will have the main() function for agent
// Basically the RPC server

// create a New Flash container object which on command will execute everything related to one container

//ajay@ubuntu:~$ chmod 755 ~/.local/share/
//ajay@ubuntu:~$ chmod 777 ~/.local/share/lxc/ ?

package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"

	"github.com/bharad1988/flashc/common"

	"github.com/bharad1988/flashc/flashagent/agentlib"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	err := agentlib.IdentifyNode()
	if err != nil {
		log.Fatal(err)
	}
	go agentlib.StartAgMonitor(&wg)

	lc := new(agentlib.LocalContainers)
	li := new(agentlib.LocalImage)
	rpc.Register(lc)
	rpc.Register(li)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+common.AgentPort)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)

}
