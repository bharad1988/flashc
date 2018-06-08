package flashlib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/bharad1988/flashc/common"
)

// AgentStatus handles erst call for agent status command
func AgentStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Agent status command issued....")
	/*
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			panic(err)
		}
		//fmt.Println(body)
		var agents []common.AgentInfo
		err = json.Unmarshal(body, &agents)
		if err != nil {
			panic(err)
		}
		err = r.Body.Close()
		if err != nil {
			panic(err)
		}
	*/
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
	//================================================
	// ResponseWriter
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err1 := json.NewEncoder(w).Encode(agents); err1 != nil {
		panic(err1)
	}
	return
}

// GetOnlineAgents returns the agents that are online
func GetOnlineAgents() []common.AgentInfo {
	var agents []common.AgentInfo
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
	defer mSession.Close()
	// controller db cursor
	mCursor := mSession.DB(mgs.DB).C(mgs.Collection)

	// get all agents first
	err = mCursor.Find(bson.M{"status": "online"}).All(&agents)
	if err != nil {
		log.Print(err)
	}
	return agents
}
