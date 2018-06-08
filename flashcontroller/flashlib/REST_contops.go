package flashlib

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/rpc"
	"strings"
	"sync"

	"github.com/bharad1988/flashc/common"
)

func getSuperCon(r *http.Request) common.SuperCon {
	var lcs common.SuperCon
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(body, &lcs)
	if err != nil {
		fmt.Println(err)
	}
	err = r.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	return lcs
}

func callConOp(op string, agent common.AgentInfo, lcs common.SuperCon, message *string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("here for agent : " + agent.AgentIP)
	client, err := rpc.DialHTTP("tcp", agent.AgentIP+":"+common.AgentPort)
	if err != nil {
		*message = *message + err.Error() + "\nAgent daemon might have stopped or crashed" + "\n"
		return
	}
	var reply string
	if op == "start" {
		err = client.Call("LocalContainers.StartContainers", lcs, &reply)
		if strings.Contains(reply, "success") {
			*message = *message + "containers started on " + agent.AgentIP + "\n"
			fmt.Printf(*message)
		} else {
			*message = *message + err.Error() + "\nFailed on " + agent.AgentIP + "\n"
			fmt.Println(*message)
		}
	}
	if op == "stop" {
		err = client.Call("LocalContainers.StopContainers", lcs, &reply)
		if strings.Contains(reply, "success") {
			*message = *message + "containers stopped on " + agent.AgentIP + "\n"
			fmt.Printf(*message)
		} else {
			*message = *message + err.Error() + "\nFailed on " + agent.AgentIP + "\n"
			fmt.Println(*message)
		}
	}
	if op == "destroy" {
		err = client.Call("LocalContainers.DestroyContainers", lcs, &reply)
		if strings.Contains(reply, "success") {
			*message = *message + "containers destroyed on " + agent.AgentIP + "\n"
			fmt.Printf(*message)
		} else {
			*message = *message + err.Error() + "\nFailed on " + agent.AgentIP + "\n"
			fmt.Println(*message)
		}
	}

	if op == "snap" {
		err = client.Call("LocalContainers.SnapContainers", lcs, &reply)
		if strings.Contains(reply, "success") {
			*message = *message + "containers snap'd on " + agent.AgentIP + "\n"
			fmt.Printf(*message)
		} else {
			*message = *message + err.Error() + "\nFailed on " + agent.AgentIP + "\n"
			fmt.Println(*message)
		}
	}

	if op == "snapdestroy" {
		err = client.Call("LocalContainers.SnapDestroy", lcs, &reply)
		if strings.Contains(reply, "success") {
			*message = *message + "snapshots destroyed on " + agent.AgentIP + "\n"
			fmt.Printf(*message)
		} else {
			*message = *message + err.Error() + "\nFailed on " + agent.AgentIP + "\n"
			fmt.Println(*message)
		}
	}
}

func callConStats(stats *[]common.ContStat, agent common.AgentInfo, lcs common.SuperCon, message *string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("here for agent : " + agent.AgentIP)
	client, err := rpc.DialHTTP("tcp", agent.AgentIP+":"+common.AgentPort)
	if err != nil {
		*message = *message + err.Error() + "\nAgent daemon might have stopped or crashed" + "\n"
		return
	}
	err = client.Call("LocalContainers.StatContainers", lcs, &stats)
	if err != nil {
		*message = *message + err.Error() + "\nAgent " + agent.AgentIP + "\n"
		return
	}
}

func callSnapList(snaps *[]common.SnapShot, agent common.AgentInfo, lcs common.SuperCon, message *string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("here for agent : " + agent.AgentIP)
	client, err := rpc.DialHTTP("tcp", agent.AgentIP+":"+common.AgentPort)
	if err != nil {
		*message = *message + err.Error() + "\nAgent daemon might have stopped or crashed" + "\n"
		return
	}
	err = client.Call("LocalContainers.SnapList", lcs, &snaps)
	if err != nil {
		*message = *message + err.Error() + "\nAgent " + agent.AgentIP + "\n"
		return
	}
}

// ContStart starts the containers with given id on the agents that are online
func ContStart(w http.ResponseWriter, r *http.Request) {
	var agents []common.AgentInfo
	var message string
	var lcs common.SuperCon
	var op string
	op = "start"
	fmt.Println("Container start command issued....")
	lcs = getSuperCon(r)
	agents = GetOnlineAgents() // defined in REST_AgentStatus
	fmt.Println(agents)
	//var conts []common.ConConfig
	var wg sync.WaitGroup
	wg.Add(len(agents))
	for _, agent := range agents {
		fmt.Println(agent)
		fmt.Println(message)
		go callConOp(op, agent, lcs, &message, &wg)
	}
	wg.Wait()
	if len(agents) == 0 {
		message = message + "\nNo Agents available"
	}
	// ResponseWriter
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err1 := json.NewEncoder(w).Encode(message); err1 != nil {
		fmt.Println(err1)
	}
	return
}

// ContStop issues stop containers on Agents that are online
func ContStop(w http.ResponseWriter, r *http.Request) {
	var agents []common.AgentInfo
	var message string
	var lcs common.SuperCon
	var op string
	op = "stop"
	fmt.Println("Container start command issued....")
	lcs = getSuperCon(r)
	agents = GetOnlineAgents() // defined in REST_AgentStatus
	fmt.Println(agents)
	var wg sync.WaitGroup
	wg.Add(len(agents))
	for _, agent := range agents {
		fmt.Println(agent)
		fmt.Println(message)
		go callConOp(op, agent, lcs, &message, &wg)
	}
	wg.Wait()
	if len(agents) == 0 {
		message = message + "\nNo Agents available"
	}
	// ResponseWriter
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err1 := json.NewEncoder(w).Encode(message); err1 != nil {
		fmt.Println(err1)
	}
	return
}

// ContDestroy issues stop containers on Agents that are online
func ContDestroy(w http.ResponseWriter, r *http.Request) {
	var agents []common.AgentInfo
	var message string
	var lcs common.SuperCon
	var op string
	op = "destroy"
	fmt.Println("Container destroy command issued....")
	lcs = getSuperCon(r)
	agents = GetOnlineAgents() // defined in REST_AgentStatus
	fmt.Println(agents)
	var wg sync.WaitGroup

	// set wait based on batch or individual container
	if lcs.Con.BatchUUID != "" {
		wg.Add(len(agents))
	} else {
		wg.Add(1)
	}

	for _, agent := range agents {
		fmt.Println(agent)
		fmt.Println(message)
		if lcs.Con.BatchUUID != "" { // run on all Agents that are online
			go callConOp(op, agent, lcs, &message, &wg)
		} else if lcs.Con.AgentUUID == agent.UUID { // run only on the agent with same AgentUUID
			fmt.Println(agent.Hostname, agent.AgentIP)
			go callConOp(op, agent, lcs, &message, &wg)
		} else { // A case where the container is not present in any online agents ( individual container destroy op )
			message = message + "\nNo Agents available"
			wg.Done()
		}
	}
	wg.Wait()
	if len(agents) == 0 {
		message = message + "\nNo Agents available"
	}
	// ResponseWriter
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err1 := json.NewEncoder(w).Encode(message); err1 != nil {
		fmt.Println(err1)
	}
	return
}

// ContSnap issues stop containers on Agents that are online
func ContSnap(w http.ResponseWriter, r *http.Request) {
	var agents []common.AgentInfo
	var message string
	var lcs common.SuperCon
	var op string
	op = "snap"
	fmt.Println("Container snapshot command issued....")
	lcs = getSuperCon(r)
	agents = GetOnlineAgents() // defined in REST_AgentStatus
	fmt.Println(agents)
	var wg sync.WaitGroup

	// set wait based on batch or individual container
	if lcs.Con.BatchUUID != "" {
		wg.Add(len(agents))
	} else {
		wg.Add(1)
	}

	for _, agent := range agents {
		fmt.Println(agent)
		fmt.Println(message)
		if lcs.Con.BatchUUID != "" { // run on all Agents that are online
			go callConOp(op, agent, lcs, &message, &wg)
		} else if lcs.Con.AgentUUID == agent.UUID { // run only on the agent with same AgentUUID
			fmt.Println(agent.Hostname, agent.AgentIP)
			go callConOp(op, agent, lcs, &message, &wg)
		} else { // A case where the container is not present in any online agents ( individual container destroy op )
			message = message + "\nNo Agents available"
			wg.Done()
		}
	}
	wg.Wait()
	if len(agents) == 0 {
		message = message + "\nNo Agents available"
	}
	// ResponseWriter
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err1 := json.NewEncoder(w).Encode(message); err1 != nil {
		fmt.Println(err1)
	}
	return
}

// ContStat issues stop containers on Agents that are online
func ContStat(w http.ResponseWriter, r *http.Request) {
	var agents []common.AgentInfo
	var message string
	var lcs common.SuperCon
	var stats []common.ContStat

	fmt.Println("Container stat command issued....")
	lcs = getSuperCon(r)
	agents = GetOnlineAgents() // defined in REST_AgentStatus
	fmt.Println(agents)
	var wg sync.WaitGroup

	// set wait based on batch or individual container
	if lcs.Con.BatchUUID != "" {
		wg.Add(len(agents))
	} else {
		wg.Add(1)
	}

	for _, agent := range agents {
		fmt.Println(agent)
		fmt.Println(message)
		if lcs.Con.BatchUUID != "" { // run on all Agents that are online
			go callConStats(&stats, agent, lcs, &message, &wg)
		} else if lcs.Con.AgentUUID == agent.UUID { // run only on the agent with same AgentUUID
			fmt.Println(agent.Hostname, agent.AgentIP)
			go callConStats(&stats, agent, lcs, &message, &wg)
		} else { // A case where the container is not present in any online agents ( individual container destroy op )
			message = message + "\nNo Agents available"
			wg.Done()
		}
	}
	wg.Wait()
	if len(agents) == 0 {
		message = message + "\nNo Agents available"
	}
	//fmt.Println(stats)
	// ResponseWriter
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err1 := json.NewEncoder(w).Encode(stats); err1 != nil {
		fmt.Println(err1)
	}
	return
}

// SnapList issues containers snap list
func SnapList(w http.ResponseWriter, r *http.Request) {
	var agents []common.AgentInfo
	var message string
	var lcs common.SuperCon
	var snaps []common.SnapShot

	fmt.Println("Container snapshot list command issued....")
	lcs = getSuperCon(r)
	agents = GetOnlineAgents() // defined in REST_AgentStatus
	fmt.Println(agents)
	var wg sync.WaitGroup

	// set wait based on batch or individual container
	if lcs.Con.BatchUUID != "" {
		wg.Add(len(agents))
	} else {
		wg.Add(1)
	}

	for _, agent := range agents {
		fmt.Println(agent.Hostname)
		fmt.Println(lcs.Con.BatchUUID)
		if lcs.Con.BatchUUID != "" { // run on all Agents that are online
			go callSnapList(&snaps, agent, lcs, &message, &wg)
		} else if lcs.Con.AgentUUID == agent.UUID { // run only on the agent with same AgentUUID
			fmt.Println(agent.Hostname, agent.AgentIP)
			go callSnapList(&snaps, agent, lcs, &message, &wg)
		} else { // A case where the container is not present in any online agents ( individual container destroy op )
			message = message + "\nNo Agents available"
			wg.Done()
		}
	}
	wg.Wait()
	fmt.Println(snaps)
	if len(agents) == 0 {
		message = message + "\nNo Agents available"
	}
	// ResponseWriter
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err1 := json.NewEncoder(w).Encode(snaps); err1 != nil {
		fmt.Println(err1)
	}
	return
}

// SnapDestroy issues stop containers on Agents that are online
func SnapDestroy(w http.ResponseWriter, r *http.Request) {
	var agents []common.AgentInfo
	var message string
	var lcs common.SuperCon
	var op string
	op = "snapdestroy"
	fmt.Println("Container snapshot command issued....")
	lcs = getSuperCon(r)
	agents = GetOnlineAgents() // defined in REST_AgentStatus
	fmt.Println(agents)
	var wg sync.WaitGroup

	// set wait based on batch or individual container
	if lcs.Con.BatchUUID != "" {
		wg.Add(len(agents))
	} else {
		wg.Add(1)
	}

	for _, agent := range agents {
		fmt.Println(agent)
		fmt.Println(message)
		if lcs.Con.BatchUUID != "" { // run on all Agents that are online
			go callConOp(op, agent, lcs, &message, &wg)
		} else if lcs.Con.AgentUUID == agent.UUID { // run only on the agent with same AgentUUID
			fmt.Println(agent.Hostname, agent.AgentIP)
			go callConOp(op, agent, lcs, &message, &wg)
		} else { // A case where the container is not present in any online agents ( individual container destroy op )
			message = message + "\nNo Agents available"
			wg.Done()
		}
	}
	wg.Wait()
	if len(agents) == 0 {
		message = message + "\nNo Agents available"
	}
	// ResponseWriter
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err1 := json.NewEncoder(w).Encode(message); err1 != nil {
		fmt.Println(err1)
	}
	return
}
