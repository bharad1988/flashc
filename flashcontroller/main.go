package main

// All GOroutines will be launched
// (-1-) : REST API server
// (-2-) : control unit - This does LB based on Agent node's weight

import (
	"sync"

	"github.com/bharad1988/flashc/flashcontroller/flashlib"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2) // Add the total number of GOroutines that have been started for which the main has to Wait
	// calls the rest server and waits on it. This main will not terminate unless the server returns ( This happens only in case of an error )
	go flashlib.StartRESTServer(&wg)
	go flashlib.StartController(&wg)

	wg.Wait()
}
