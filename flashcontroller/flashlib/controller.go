package flashlib

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/bharad1988/flashc/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GetAgentStatus gets the status of agents
func GetAgentStatus(mCursor *mgo.Collection, agents *[]common.AgentInfo) error {

	for _, agent := range *agents {
		_, err := net.Dial("tcp", agent.AgentIP+":"+common.AgentPort)
		if err != nil { // agent not running
			//update controller db
			agent.Status = "offline"
			agname := bson.M{"uuid": agent.UUID}
			mCursor.Update(agname, agent)
			continue
		}
		amgs := common.MongoServer{
			Server:     agent.AgentIP,
			Port:       "27017",
			DB:         "node",
			Collection: "info",
		}
		//check if agent mongo is running by connection to session
		amSession, err := common.MongoConnect(&amgs)
		if err != nil { // agent mongo not running
			//update controller db
			agent.Status = "offline"
			agname := bson.M{"uuid": agent.UUID}
			mCursor.Update(agname, agent)
			continue
		}

		amCursor := amSession.DB(amgs.DB).C(amgs.Collection)
		var AgInfo common.AgentInfo
		err = amCursor.Find(nil).One(&AgInfo)
		if err != nil {
			log.Print(err)
		}

		//update controller db
		agent.Status = AgInfo.Status // Gets the status as online
		agent.Hostname = AgInfo.Hostname
		agent.TotalMemUsage = AgInfo.TotalMemUsage
		agname := bson.M{"uuid": agent.UUID}
		mCursor.Update(agname, agent)
		amSession.Close()
	}
	return nil
}

// MonitorAgentHealth checks Agents' health
func MonitorAgentHealth() {
	fmt.Println("Checking status....")
	//============================================================================
	// For loop to pool for status
	//============================================================================
	for true {
		time.Sleep(5 * time.Second)
		mgs := common.MongoServer{
			Server:     "127.0.0.1",
			Port:       "27017",
			DB:         "ctlr",
			Collection: "info",
		}
		//controller db session
		mSession, err := common.MongoConnect(&mgs)
		if err != nil {
			log.Print(err)
		}
		// controller db cursor
		mCursor := mSession.DB(mgs.DB).C(mgs.Collection)
		var agents []common.AgentInfo
		err = mCursor.Find(nil).All(&agents)
		if err != nil {
			log.Print(err)
		}
		fmt.Println(agents)

		err = GetAgentStatus(mCursor, &agents)
		if err != nil {
			log.Print(err)
		}
		mSession.Close()
	}
	//============================================================================
	// loop ENDs
	//============================================================================
}

// StartController does maintenece work
// polls for status
// load balacing
func StartController(wg *sync.WaitGroup) {
	defer wg.Done()
	MonitorAgentHealth()
}
