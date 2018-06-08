package flashlib

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/mgo.v2/bson"

	"github.com/bharad1988/flashc/common"
)

/*
func loadbalance() (agents common.AgentInfo) {
}
*/

func callSuperCon(w http.ResponseWriter, agent common.AgentInfo, lcs common.SuperCon, message *string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("here for agent : " + agent.AgentIP)
	client, err := rpc.DialHTTP("tcp", agent.AgentIP+":"+common.AgentPort)
	if err != nil {
		*message = *message + err.Error() + "\nAgent daemon might have stopped or crashed" + "\n"
		return
	}
	var reply string
	err = client.Call("LocalContainers.Create", lcs, &reply)

	//waitForResults(statusObj)
	if strings.Contains(reply, "success") {
		*message = *message + "Create containers issued on " + agent.AgentIP + "\n"
		fmt.Printf(*message)
	} else {
		*message = *message + err.Error() + "\nFailed on " + agent.AgentIP + "\n"
		fmt.Println(*message)
	}

}

// Create lcx
func Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("launch command issued....")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	//fmt.Println(body)

	var lCon common.SuperCon
	var message string

	err = json.Unmarshal(body, &lCon)
	fmt.Println(lCon)
	err = r.Body.Close()
	if err != nil {
		panic(err)
	}

	// check agents in MongoServer ctrl db
	// stage 1: Pre LB
	// Find the online agents
	// ditribute the conainers on all agents
	// Call the agent rpc to create launcher - send in the cound in a single rpc call
	// Agent adds suffix to container name (unique to local agents ?)
	// Connect to mongo and insert agent's info ( if already not inserted )
	mgs := common.MongoServer{
		Server:     "127.0.0.1",
		Port:       "27017",
		DB:         "ctlr",
		Collection: "info",
	}
	mSession, err := common.MongoConnect(&mgs)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		//w.WriteHeader() // unprocessable entity
		err1 := json.NewEncoder(w).Encode("Most likely : Controller Mongo Not running")
		if err1 != nil {
			panic(err1)
		}
		return
	}
	defer mSession.Close()
	mCursor := mSession.DB(mgs.DB).C(mgs.Collection)
	num, err := mCursor.Find(bson.M{"status": "online"}).Count()

	var agents []common.AgentInfo
	err = mCursor.Find(bson.M{"status": "online"}).All(&agents)
	if err != nil {
		panic(err)
	}
	fmt.Println(num)
	lCon.Con.BatchUUID, err = common.NewUUID()
	if err != nil {
		panic(err)
	}

	//	var statusObjs []*rpc.Call
	var wg sync.WaitGroup
	wg.Add(len(agents))
	for i, agent := range agents {
		// LB will decide - number of containers on which agent
		// later stages -
		// need to do some research on lxc's usage and available resources

		lcs := common.SuperCon{}
		fmt.Println(agent.AgentIP)
		fmt.Println(i)
		// tmp value

		total, err := strconv.ParseInt(lCon.Count, 10, 64)
		if err != nil {
			panic(err)
		}
		scount := (int(total) / num) * (i)
		acount := (int(total) / num) * (i + 1)
		fmt.Println(scount)
		fmt.Println(acount)
		lcs = lCon
		lcs.Start = strconv.FormatInt(int64(scount), 16)
		lcs.Count = strconv.FormatInt(int64(acount), 16)

		//Then it can make a remote call:

		go callSuperCon(w, agent, lcs, &message, &wg)
		//==================================================================================

	}
	wg.Wait()
	if num == 0 {
		message = message + "\nNo Agents available"
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err1 := json.NewEncoder(w).Encode(message); err1 != nil {
		panic(err1)
	}
	return

}
