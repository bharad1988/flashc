package flashlib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/bharad1988/flashc/common"
)

// ContList handles erst call for agent status command
func ContList(w http.ResponseWriter, r *http.Request) {
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

	// get all agents first
	err = mCursor.Find(bson.M{"status": "online"}).All(&agents)
	fmt.Println(agents)
	mSession.Close()

	fmt.Println("this command got executed")
	if err != nil {
		log.Print(err)
	}

	var conts []common.ConConfig

	for _, agent := range agents {
		var aconts []common.ConConfig
		mgs = common.MongoServer{
			Server:     agent.AgentIP,
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
		mCursor = mSession.DB(mgs.DB).C(mgs.Collection)
		// to ensure no empty entries
		err = mCursor.Find(bson.M{"uuid": &bson.RegEx{Pattern: "-", Options: ""}}).All(&aconts)
		//c.Find(bson.M{"abc": &bson.RegEx{Pattern: "efg", Options: "i"}})
		if err != nil {
			log.Print(err)
		}
		conts = append(conts, aconts...)
		fmt.Println(conts)
		mSession.Close()

	}
	// ResponseWriter
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err1 := json.NewEncoder(w).Encode(conts); err1 != nil {
		panic(err1)
	}
	return
}
