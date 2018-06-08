package flashlib

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/bharad1988/flashc/common"

	"gopkg.in/mgo.v2/bson"
)

//Register function will register the Agent on the controller
func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("register command issued....")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	fmt.Println(body)
	err = r.Body.Close()
	if err != nil {
		panic(err)
	}

	// var AgNode Agent
	var AgNode common.AgentInfo
	err = json.Unmarshal(body, &AgNode)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err1 := json.NewEncoder(w).Encode(err); err1 != nil {
			panic(err1)
		}
	}
	//verify Agent ( check if node is identified as agent or not and get the UUID)
	uuid, err := common.VerifyAgent(AgNode.AgentIP)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		//w.WriteHeader() // unprocessable entity
		if err1 := json.NewEncoder(w).Encode(err.Error()); err1 != nil {
			panic(err1)
		}
		return
	}

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
	AgNode.UUID = uuid

	//insert if not inserted
	count, err := mCursor.Find(bson.M{"uuid": uuid}).Count()
	if count == 0 {
		err = mCursor.Insert(AgNode)
		if err != nil {
			panic(err)
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		//w.WriteHeader() // unprocessable entity
		err1 := json.NewEncoder(w).Encode("flashagent already registered")
		if err1 != nil {
			panic(err1)
		}
		return
	}
	if err != nil {
		panic(err)
	}

	// write status back to REST API client
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		//w.WriteHeader() // unprocessable entity
		if err1 := json.NewEncoder(w).Encode("Failed to update local mongo. Agent not registered on controller"); err1 != nil {
			panic(err1)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(uuid)
	//fmt.Println(AgNode.AgentIP)
}

//UnRegister function will unregister the Agent on the controller
func UnRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Println("unregister command issued....")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	fmt.Println(body)
	err = r.Body.Close()
	if err != nil {
		panic(err)
	}

	// var AgNode Agent
	var AgNode common.AgentInfo
	err = json.Unmarshal(body, &AgNode)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err1 := json.NewEncoder(w).Encode(err); err1 != nil {
			panic(err1)
		}
	}

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

	err = mCursor.Remove(bson.M{"uuid": AgNode.UUID})
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		//w.WriteHeader() // unprocessable entity
		err1 := json.NewEncoder(w).Encode("Failed to remove from Mongo")
		if err1 != nil {
			panic(err1)
		}
		return
	}

	// write status back to REST API client
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		//w.WriteHeader() // unprocessable entity
		if err1 := json.NewEncoder(w).Encode("Failed to update local mongo. Agent not registered on controller"); err1 != nil {
			panic(err1)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode("successfully removed the Agent with UUID" + AgNode.UUID)
	//fmt.Println(AgNode.AgentIP)
}
