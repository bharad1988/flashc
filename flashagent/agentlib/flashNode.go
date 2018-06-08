// Package agentlib
//Node status object , monitoring feedback to controller
//This file will have all the info. related to a node
//This will collect stats of containers, node health (LB)

package agentlib

import (
	"errors"
	"log"
	"os"

	"github.com/bharad1988/flashc/common"
	"gopkg.in/mgo.v2"
)

// MgoInsert inserts container entry to local mongo service
func MgoInsert(locon *common.ConConfig) error {

	ms := common.MongoServer{
		Server:     "127.0.0.1",
		Port:       "27017",
		DB:         "cdata",
		Collection: "info",
	}
	session, err := mgo.Dial(ms.Server + ":" + ms.Port)
	if err != nil {
		//panic(err)
		return err
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(ms.DB).C(ms.Collection)

	err = c.Insert(&locon)
	if err != nil {
		//panic(err)
		return err
	}

	// insert container monitor structure to stat collection in cdata DB
	mgs := common.MongoServer{
		Server:     "127.0.0.1",
		Port:       "27017",
		DB:         "cdata",
		Collection: "stat",
	}
	// agent session
	mSession, err := common.MongoConnect(&mgs)
	if err != nil {
		log.Print(err)
	}
	// agent cursor to fetch all containers
	mCursor := mSession.DB(mgs.DB).C(mgs.Collection)
	stats := new(common.ContStat)
	stats.ContUUID = locon.UUID
	mCursor.Insert(stats)

	return nil
}

//IdentifyNode is to Provide the Agent node with an ID
func IdentifyNode() error {
	mgs := common.MongoServer{
		Server:     "127.0.0.1",
		Port:       "27017",
		DB:         "node",
		Collection: "info",
	}
	mSession, err := common.MongoConnect(&mgs)
	if err != nil {
		//panic(err)
		return err
	}
	defer mSession.Close()
	var AgNode common.AgentInfo
	mCursor := mSession.DB(mgs.DB).C(mgs.Collection)

	n, err := mCursor.Count()
	if err != nil {
		//panic(err)
		return err
	}
	if n != 1 {
		uuid, err := common.NewUUID()
		AgNode.UUID = uuid
		AgNode.AgentIP = "127.0.0.1"
		hostname, err := os.Hostname()
		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
			//*reply = err.Error()
			return err
		}
		AgNode.Hostname = hostname
		err = mCursor.Insert(AgNode)
		if err != nil {
			//panic(err)
			return err
		}
	}
	return nil
}

// GetAgentParam function returns Agent params
func GetAgentParam() (common.AgentInfo, error) {
	agNode := common.AgentInfo{}
	mgs := common.MongoServer{
		Server:     "127.0.0.1",
		Port:       "27017",
		DB:         "node",
		Collection: "info",
	}
	mSession, err := common.MongoConnect(&mgs)
	if err != nil {
		//panic(err)
		return agNode, err
	}
	defer mSession.Close()
	mCursor := mSession.DB(mgs.DB).C(mgs.Collection)
	err = mCursor.Find(nil).One(&agNode)
	if err != nil {
		reterr := errors.New("Failed to find the agent info. Node not yet identified as Agent")
		return agNode, reterr
	}
	return agNode, err
}
