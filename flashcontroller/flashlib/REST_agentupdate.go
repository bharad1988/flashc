package flashlib

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/bharad1988/flashc/common"
)

// AgentUpdate updates agent ip on the Controller
func AgentUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Agent update command issued....")

	// connect to controller MongoServer
	mgs := common.MongoServer{
		Server:     "127.0.0.1",
		Port:       "27017",
		DB:         "ctlr",
		Collection: "info",
	}
	//controller db session
	mSession, err := common.MongoConnect(&mgs)
	if err != nil { // if mongo server not running on controller throw this error
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err1 := json.NewEncoder(w).Encode(err.Error()); err1 != nil {
			panic(err1)
		}
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	//fmt.Println(body)
	var agent common.AgentInfo
	err = json.Unmarshal(body, &agent)
	if err != nil {
		panic(err)
	}
	err = r.Body.Close()
	if err != nil {
		panic(err)
	}
	_, err = net.Dial("tcp", agent.AgentIP+":"+common.AgentPort)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err1 := json.NewEncoder(w).Encode(err.Error()); err1 != nil {
			panic(err1)
		}
		return

	}

	// controller db cursor
	mCursor := mSession.DB(mgs.DB).C(mgs.Collection)
	var nagent common.AgentInfo

	agname := bson.M{"uuid": agent.UUID}
	fmt.Println(agname)
	fmt.Println(nagent)
	err = mCursor.Update(agname, bson.M{"$set": bson.M{"agentip": agent.AgentIP}})

	mCursor.Find(bson.M{"uuid": agent.UUID}).One(&nagent)

	if err != nil {
		log.Print(err)
	}
	message := "New IP " + nagent.AgentIP + " of Agent with UUID " + agent.UUID
	// ResponseWriter
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err1 := json.NewEncoder(w).Encode(message); err1 != nil {
		panic(err1)
	}
	return
}
