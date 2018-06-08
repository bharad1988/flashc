package agentlib

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gopkg.in/lxc/go-lxc.v2"

	"github.com/bharad1988/flashc/common"
	"gopkg.in/mgo.v2/bson"
)

// UpdateState updates the state of containers
func UpdateState() {
	var conts []common.ConConfig
	var totalMemUsage lxc.ByteSize

	mgs := common.MongoServer{
		Server:     "127.0.0.1",
		Port:       "27017",
		DB:         "cdata",
		Collection: "info",
	}
	// agent session
	mSession, err := common.MongoConnect(&mgs)
	if err != nil {
		log.Print(err)
	}
	// agent cursor to fetch all containers
	mCursor := mSession.DB(mgs.DB).C(mgs.Collection)
	statCursor := mSession.DB(mgs.DB).C("stat")
	// to ensure no empty entries
	err = mCursor.Find(bson.M{"uuid": &bson.RegEx{Pattern: "-", Options: ""}}).All(&conts)
	//c.Find(bson.M{"abc": &bson.RegEx{Pattern: "efg", Options: "i"}})
	if err != nil {
		log.Print(err)
	}
	//conts = append(conts, aconts...)
	//fmt.Println(conts)

	for _, cont := range conts {
		c, nerr := lxc.NewContainer(cont.Name, cont.Lxcpath)
		if nerr != nil {
			log.Print(err)
		}
		cont.State = c.State().String()
		//fmt.Println(cont.State)
		cstate := bson.M{"uuid": cont.UUID}
		mCursor.Update(cstate, cont)

		if c.Running() {
			// update stats
			//contstats := new(common.ContStat)
			cuuid := bson.M{"contuuid": cont.UUID}
			//ifaces := c.InterfaceStats()
			memory, nerr := c.MemoryUsage()
			fmt.Println(memory)
			if memory != -1.00 {
				totalMemUsage = totalMemUsage + memory
			}

			if nerr != nil {
				log.Print(err)
			}
			ifacestats, nerr := c.InterfaceStats()
			if nerr != nil {
				log.Print(err)
			}
			ips, nerr := c.IPAddresses()
			if nerr != nil {
				log.Print(err)
			}
			intfaces := []common.Iface{}
			for iface := range ifacestats {
				fmt.Println(iface)
				var intface common.Iface
				intface.Name = iface
				intface.Rx = ifacestats[iface]["rx"].String()
				intface.Tx = ifacestats[iface]["tx"].String()
				intfaces = append(intfaces, intface)
				//fmt.Println(intfaces)
			}
			for i, ip := range ips {
				fmt.Println(ip)
				intfaces[i].IP = ip
				fmt.Println(intfaces[i])
			}
			data := bson.M{"$set": bson.M{"memory": memory.String(), "intface": intfaces}}
			err = statCursor.Update(cuuid, data)
			if err != nil {
				log.Print(err)
			}
			cpustats, _ := c.CPUStats()
			fmt.Println(cpustats["user"])
			blkio, _ := c.BlkioUsage()

			fmt.Printf("Blokio - %v", blkio.String())
		}
		//fmt.Println(contstats)
	} // end of for loop for containers

	nmgs := common.MongoServer{
		Server:     "127.0.0.1",
		Port:       "27017",
		DB:         "node",
		Collection: "info",
	}
	// agent session
	nmSession, err := common.MongoConnect(&nmgs)
	if err != nil {
		log.Print(err)
	}
	// agent cursor to fetch all containers
	nmCursor := mSession.DB(nmgs.DB).C(nmgs.Collection)
	//var node common.AgentInfo
	totmem := bson.M{"$set": bson.M{"totalmemusage": totalMemUsage.String(), "status": "online"}}
	err = nmCursor.Update(nil, totmem)
	if err != nil {
		log.Print(err)
	}
	defer nmSession.Close()

	fmt.Println(totalMemUsage.String())
	mSession.Close()

}

// StartAgMonitor starts the monitor in a loop
func StartAgMonitor(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Monitoring container state...")
	for true {
		time.Sleep(5 * time.Second)
		UpdateState()
	}

}
