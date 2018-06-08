package flashlib

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/rpc"

	"github.com/bharad1988/flashc/common"
)

// ImageBuildObs builds image
func ImageBuildObs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Image build command issued....")

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	//fmt.Println(body)

	var locImage common.SuperImage
	var message string
	var agents []common.AgentInfo
	var agent common.AgentInfo
	var reply common.ConImage
	var retObj common.HTTPRetObj

	err = json.Unmarshal(body, &locImage)
	fmt.Println(locImage)
	err = r.Body.Close()
	if err != nil {
		panic(err)
	}
	agents = GetOnlineAgents() // defined in REST_AgentStatus
	fmt.Println(agents)

	for _, item := range agents {
		if item.UUID == locImage.AgentUUID {
			agent = item
			break
		}
	}

	if len(agents) == 0 {
		message = message + "\nNo Agents available"
	} else {
		client, err := rpc.DialHTTP("tcp", agent.AgentIP+":"+common.AgentPort)
		if err != nil {
			message = message + err.Error() + "\nAgent daemon might have stopped or crashed" + "\n"
			return
		}

		err = client.Call("LocalImage.ImageBuild", locImage, &reply)
	}

	if message == "" {
		message = "success"
	}
	// Build return object
	retObj.Obj = reply
	retObj.Message = message
	// ResponseWriter
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err1 := json.NewEncoder(w).Encode(retObj); err1 != nil {
		fmt.Println(err1)
	}
	return

}

// ImageListObs lists image
func ImageListObs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Image list command issued....")

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	//fmt.Println(body)

	var locImage common.SuperImage
	var message string
	var agents []common.AgentInfo
	var agent common.AgentInfo
	var reply []common.ConImage
	var retObj common.HTTPRetObj

	err = json.Unmarshal(body, &locImage)
	fmt.Println(locImage)
	err = r.Body.Close()
	if err != nil {
		panic(err)
	}
	agents = GetOnlineAgents() // defined in REST_AgentStatus
	fmt.Println(agents)

	for _, item := range agents {
		if item.UUID == locImage.AgentUUID {
			agent = item
			break
		}
	}

	if len(agents) == 0 {
		message = message + "\nNo Agents available"
	} else {
		client, err := rpc.DialHTTP("tcp", agent.AgentIP+":"+common.AgentPort)
		if err != nil {
			message = message + err.Error() + "\nAgent daemon might have stopped or crashed" + "\n"
			return
		}

		err = client.Call("LocalImage.ImageList", locImage, &reply)
	}

	if message == "" {
		message = "success"
	}
	// Build return object
	retObj.Obj = reply
	retObj.Message = message
	// ResponseWriter
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err1 := json.NewEncoder(w).Encode(retObj); err1 != nil {
		fmt.Println(err1)
	}
	return

}
