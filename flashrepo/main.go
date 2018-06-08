package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"

	"github.com/bharad1988/flashc/common"
	"github.com/bharad1988/flashc/flashrepo/repolib"
	//"github.com/bharad1988/flashc/flashagent/agentlib"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	//lc := new(agentlib.LocalContainers)
	li := new(repolib.LocalImage)
	rpc.Register(li)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+common.RepoPort)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	defer l.Close()
	http.Serve(l, nil)

}
