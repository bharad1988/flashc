package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bharad1988/flashc/common"
	"github.com/spf13/cobra"
)

// contstartCmd represents the contstart command
var contstartCmd = &cobra.Command{
	Use:   "start",
	Short: "starts the containers with the given UUID",
	Long: `This command starts containers with given batch UUID.
	It looks for all the registered agents that are up and
	starts the containers from the given batch`,
	Run: ContStart,
}
var contstopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stops the containers with the given UUID",
	Long: `This command stops containers with given batch UUID.
	It looks for all the registered agents that are up and
	stops the containers from the given batch`,
	Run: ContStop,
}

var contdestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroys the containers with the given UUID",
	Long: `This command destroys containers with given batch UUID.
	It looks for all the registered agents that are up and
	destroys the containers from the given batch`,
	Run: ContDestroy,
}

var contsnapCmd = &cobra.Command{
	Use:   "snap",
	Short: "snapshots the containers with the given UUID",
	Long: `This command snapshots containers with given batch UUID.
	It looks for all the registered agents that are up and
	snapshots the containers from the given batch`,
	Run: ContSnap,
}

var contstatCmd = &cobra.Command{
	Use:   "stat",
	Short: "stats the containers with the given UUID",
	Long: `This command stats containers with given batch UUID.
	It looks for all the registered agents that are up and
	destroys the containers from the given batch`,
	Run: ContStat,
}

// sub commands for snapshots
var snaplistCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the snapshots for the given container",
	Long: `This command lists the snapshots of the given container
	displays them`,
	Run: SnapList,
}

// snap shot destroy command
var snapdestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the snapshots for the given container",
	Long: `This command destroys the snapshots of the given container
	displays them`,
	Run: SnapDestroy,
}

// ContStart start containers
func ContStart(cmd *cobra.Command, args []string) {
	fmt.Println("Container start command issued")
	err := checkController(ControllerNode)
	if err != nil {
		fmt.Printf("\n\n\nCheck if Controller is running on the node with IP - " + ControllerNode + "\n")
		fmt.Println("Commands cannot be executed on that node")
		fmt.Printf("use -C option to preovide a different controller\n\n\n")
		return
	}
	// make call to controller.
	url := "http://" + ControllerNode + ":" + common.ControllerPort + "/container/start"
	//bodyType := "application/json; charset=UTF-8"
	fmt.Println(url)
	var lCon common.SuperCon
	lCon.Con.BatchUUID = batchUUID
	// marshal the structure to json format
	bodyjson, err := json.Marshal(lCon)
	if err != nil {
		panic(err)
	}

	// convert json to io bytes
	content := bytes.NewReader(bodyjson)
	resp, err := http.Post(url, common.JSONBodyType, content)
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(resp.Body).Decode(&respStatus)
	if err != nil {
		panic(err)
	}
	fmt.Println(respStatus)

}

// ContStop stops containers
func ContStop(cmd *cobra.Command, args []string) {
	fmt.Println("Container stop command issued")
	err := checkController(ControllerNode)
	if err != nil {
		fmt.Printf("\n\n\nCheck if Controller is running on the node with IP - " + ControllerNode + "\n")
		fmt.Println("Commands cannot be executed on that node")
		fmt.Printf("use -C option to preovide a different controller\n\n\n")
		return
	}
	// make call to controller.
	url := "http://" + ControllerNode + ":" + common.ControllerPort + "/container/stop"
	//bodyType := "application/json; charset=UTF-8"
	fmt.Println(url)
	var lCon common.SuperCon
	lCon.Con.BatchUUID = batchUUID
	// marshal the structure to json format
	bodyjson, err := json.Marshal(lCon)
	if err != nil {
		panic(err)
	}

	// convert json to io bytes
	content := bytes.NewReader(bodyjson)
	resp, err := http.Post(url, common.JSONBodyType, content)
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(resp.Body).Decode(&respStatus)
	if err != nil {
		panic(err)
	}
	fmt.Println(respStatus)

}

// ContDestroy stops containers
func ContDestroy(cmd *cobra.Command, args []string) {
	var lCon common.SuperCon
	setArg := false
	err := checkController(ControllerNode)
	if err != nil {
		fmt.Printf("\n\n\nCheck if Controller is running on the node with IP - " + ControllerNode + "\n")
		fmt.Println("Commands cannot be executed on that node")
		fmt.Printf("use -C option to preovide a different controller\n\n\n")
		return
	}
	if batchUUID != "" {
		setArg = true
	}
	lCon.Con.BatchUUID = batchUUID
	if containerUUID != "" && !setArg {
		if agentUUID == "" {
			fmt.Println("If a specific container has to be deleted, then its agent UUID has to be entered")
		} else {
			lCon.Con.AgentUUID = agentUUID
			lCon.Con.UUID = containerUUID
			fmt.Println(lCon.Con)
			setArg = true
		}
	}
	if !setArg {
		fmt.Println("Provide the batchUUID or containerUUID along with its agentUUID")
		fmt.Println(setArg)
		return
	}
	lCon.DC.WithSnapShots = destroySnaps
	/*
		if setArg {
			fmt.Println(lCon.Con)
			fmt.Println(setArg)
			return
		}
	*/
	fmt.Println("Container destroy command issued")
	// make call to controller.
	url := "http://" + ControllerNode + ":" + common.ControllerPort + "/container/destroy"
	//bodyType := "application/json; charset=UTF-8"
	fmt.Println(url)

	// marshal the structure to json format
	bodyjson, err := json.Marshal(lCon)
	if err != nil {
		panic(err)
	}

	// convert json to io bytes
	content := bytes.NewReader(bodyjson)
	resp, err := http.Post(url, common.JSONBodyType, content)
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(resp.Body).Decode(&respStatus)
	if err != nil {
		panic(err)
	}
	fmt.Println(respStatus)
}

// ContStat stops containers
func ContStat(cmd *cobra.Command, args []string) {
	var lCon common.SuperCon
	setArg := false
	err := checkController(ControllerNode)
	if err != nil {
		fmt.Printf("\n\n\nCheck if Controller is running on the node with IP - " + ControllerNode + "\n")
		fmt.Println("Commands cannot be executed on that node")
		fmt.Printf("use -C option to preovide a different controller\n\n\n")
		return
	}
	if batchUUID != "" {
		setArg = true
	}
	lCon.Con.BatchUUID = batchUUID
	if containerUUID != "" && !setArg {
		if agentUUID == "" {
			fmt.Println("If a specific container has to be deleted, then its agent UUID has to be entered")
		} else {
			lCon.Con.AgentUUID = agentUUID
			lCon.Con.UUID = containerUUID
			setArg = true
		}
	}
	if !setArg {
		fmt.Println("Provide the batchUUID or containerUUID along with its agentUUID")
		return
	}
	/*
		if setArg {
			fmt.Println(lCon.Con)
			fmt.Println(setArg)
			return
		}
	*/
	fmt.Println("Container stat command issued")
	// make call to controller.
	url := "http://" + ControllerNode + ":" + common.ControllerPort + "/container/stat"
	bodyjson, err := json.Marshal(lCon)
	if err != nil {
		panic(err)
	}

	// convert json to io bytes
	content := bytes.NewReader(bodyjson)
	resp, err := http.Post(url, common.JSONBodyType, content)
	if err != nil {
		panic(err)
	}
	var containerStats []common.ContStat
	err = json.NewDecoder(resp.Body).Decode(&containerStats)
	if err != nil {
		panic(err)
	}
	for _, conStats := range containerStats {
		fmt.Println("---------------------------------------------------------------------------------------------------")
		fmt.Printf("Container UUID\t\t\t: %v\n\tMemory\t\t\t: %v\n", conStats.ContUUID, conStats.Memory)
		for _, iface := range conStats.Intface {
			fmt.Printf("\tInterface Name\t\t: %v  \n", iface.Name)
			fmt.Printf("\t\tInterface IP\t: %v  \n", iface.IP)
			fmt.Printf("\t\tData Recvd.\t: %v  \n", iface.Rx)
			fmt.Printf("\t\tData Sent\t: %v  \n", iface.Tx)
		}
	}
	fmt.Println("---------------------------------------------------------------------------------------------------")
}

// ContSnap stops containers
func ContSnap(cmd *cobra.Command, args []string) {
	var lCon common.SuperCon
	setArg := false
	err := checkController(ControllerNode)
	if err != nil {
		fmt.Printf("\n\n\nCheck if Controller is running on the node with IP - " + ControllerNode + "\n")
		fmt.Println("Commands cannot be executed on that node")
		fmt.Printf("use -C option to preovide a different controller\n\n\n")
		return
	}
	if batchUUID != "" {
		setArg = true
	}
	lCon.Con.BatchUUID = batchUUID
	if containerUUID != "" && !setArg {
		if agentUUID == "" {
			fmt.Println("If a specific container has to be snapped, then its agent UUID has to be entered")
		} else {
			lCon.Con.AgentUUID = agentUUID
			lCon.Con.UUID = containerUUID
			fmt.Println(lCon.Con)
			setArg = true
		}
	}
	if !setArg {
		fmt.Println("Provide the batchUUID or containerUUID along with its agentUUID")
		fmt.Println(setArg)
		return
	}

	fmt.Println("Container snapshot command issued")
	// make call to controller.
	url := "http://" + ControllerNode + ":" + common.ControllerPort + "/container/snapshot"
	//bodyType := "application/json; charset=UTF-8"
	fmt.Println(url)

	// marshal the structure to json format
	bodyjson, err := json.Marshal(lCon)
	if err != nil {
		panic(err)
	}

	// convert json to io bytes
	content := bytes.NewReader(bodyjson)
	resp, err := http.Post(url, common.JSONBodyType, content)
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(resp.Body).Decode(&respStatus)
	if err != nil {
		panic(err)
	}
	fmt.Println(respStatus)
}

// SnapList lists the snapshots of the container
func SnapList(cmd *cobra.Command, args []string) {
	var lCon common.SuperCon
	setArg := false
	err := checkController(ControllerNode)
	if err != nil {
		fmt.Printf("\n\n\nCheck if Controller is running on the node with IP - " + ControllerNode + "\n")
		fmt.Println("Commands cannot be executed on that node")
		fmt.Printf("use -C option to preovide a different controller\n\n\n")
		return
	}
	if batchUUID != "" {
		setArg = true
	}
	lCon.Con.BatchUUID = batchUUID

	if containerUUID != "" && !setArg {
		if agentUUID == "" {
			fmt.Println("If a specific container has to be deleted, then its agent UUID has to be entered")
		} else {
			lCon.Con.AgentUUID = agentUUID
			lCon.Con.UUID = containerUUID
			fmt.Println(lCon.Con)
			setArg = true
		}
	}
	if !setArg {
		fmt.Println("Provide the batchUUID or containerUUID along with its agentUUID")
		fmt.Println(setArg)
		return
	}

	fmt.Println("Container snap list command issued")
	// make call to controller.
	url := "http://" + ControllerNode + ":" + common.ControllerPort + "/container/snapshot/list"
	//bodyType := "application/json; charset=UTF-8"
	fmt.Println(url)

	// marshal the structure to json format
	bodyjson, err := json.Marshal(lCon)
	if err != nil {
		panic(err)
	}

	// convert json to io bytes
	content := bytes.NewReader(bodyjson)
	resp, err := http.Post(url, common.JSONBodyType, content)
	if err != nil {
		panic(err)
	}
	var containerSnaps []common.SnapShot
	err = json.NewDecoder(resp.Body).Decode(&containerSnaps)
	if err != nil {
		panic(err)
	}
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Printf("ContainerName\t\t\t| Timestamp\t\t\t\t| SnapName\t\t\t\t| Path\n")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------")
	for _, snap := range containerSnaps {
		fmt.Printf("%s\t\t\t| %s\t\t\t| %s\t\t\t\t| %s\n", snap.ContainerName, snap.TimeStamp, snap.Name, snap.Path)
	}
}

// SnapDestroy lists the snapshots of the container
func SnapDestroy(cmd *cobra.Command, args []string) {
	var lCon common.SuperCon
	setArg := false
	err := checkController(ControllerNode)
	if err != nil {
		fmt.Printf("\n\n\nCheck if Controller is running on the node with IP - " + ControllerNode + "\n")
		fmt.Println("Commands cannot be executed on that node")
		fmt.Printf("use -C option to preovide a different controller\n\n\n")
		return
	}

	if batchUUID != "" {
		setArg = true
	}
	lCon.Con.BatchUUID = batchUUID
	if containerUUID != "" && !setArg {
		if agentUUID == "" {
			fmt.Println("If a specific container has to be snapped, then its agent UUID has to be entered")
		} else {
			lCon.Con.AgentUUID = agentUUID
			lCon.Con.UUID = containerUUID
			fmt.Println(lCon.Con)
			setArg = true
		}
	}
	if setArg { // set the name of snapshot to be deleted [name like snap0 or * for destroy all ]
		lCon.Snaps.Name = snapname
	}
	if !setArg {
		fmt.Println("Provide containerUUID along with its agentUUID")
		fmt.Println(setArg)
		return
	}

	fmt.Println("Container snap destroy command issued")
	// make call to controller.
	url := "http://" + ControllerNode + ":" + common.ControllerPort + "/container/snapshot/destroy"
	//bodyType := "application/json; charset=UTF-8"
	fmt.Println(url)

	// marshal the structure to json format
	bodyjson, err := json.Marshal(lCon)
	if err != nil {
		panic(err)
	}

	// convert json to io bytes
	content := bytes.NewReader(bodyjson)
	resp, err := http.Post(url, common.JSONBodyType, content)
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(resp.Body).Decode(&respStatus)
	if err != nil {
		panic(err)
	}
	fmt.Println(respStatus)
}

func init() {
	containerCmd.AddCommand(contstartCmd)
	containerCmd.AddCommand(contstopCmd)
	containerCmd.AddCommand(contstatCmd)
	containerCmd.AddCommand(contdestroyCmd)
	containerCmd.AddCommand(contsnapCmd)

	contstartCmd.Flags().StringVarP(&batchUUID, "batchUUID", "b", "", "provide the batchUUID of the container batch to be started")
	contstopCmd.Flags().StringVarP(&batchUUID, "batchUUID", "b", "", "provide the batchUUID of the container batch to be stopped")
	contdestroyCmd.Flags().StringVarP(&batchUUID, "batchUUID", "b", "", "provide the BatchUUID of the container batch to be destroyed, batchUUID overrides other params")
	contdestroyCmd.Flags().StringVarP(&agentUUID, "agentUUID", "a", "", "provide the AgentUUID of the container batch to be destroyed")
	contdestroyCmd.Flags().StringVarP(&containerUUID, "containerUUID", "c", "", "provide the containerUUID of the container batch to be destroyed")
	contdestroyCmd.Flags().StringVarP(&destroySnaps, "destroySnaps", "s", "false", "set the flag to true if you want the snapshots to be deleted along with destroy. If not snapped containers will not be deleted")
	contstatCmd.Flags().StringVarP(&batchUUID, "batchUUID", "b", "", "provide the batchUUID of the container batch to be started")
	contstatCmd.Flags().StringVarP(&agentUUID, "agentUUID", "a", "", "provide the AgentUUID of the container batch to be statd")
	contstatCmd.Flags().StringVarP(&containerUUID, "containerUUID", "c", "", "provide the containerUUID of the container batch to be stat'd")
	contsnapCmd.Flags().StringVarP(&batchUUID, "batchUUID", "b", "", "provide the BatchUUID of the container batch to be snapped, batchUUID overrides other params")
	contsnapCmd.Flags().StringVarP(&agentUUID, "agentUUID", "a", "", "provide the AgentUUID of the container batch to be snapped")
	contsnapCmd.Flags().StringVarP(&containerUUID, "containerUUID", "c", "", "provide the containerUUID of the container batch to be snapped")

	// snapshot sub Commands
	contsnapCmd.AddCommand(snaplistCmd)
	contsnapCmd.AddCommand(snapdestroyCmd)

	snaplistCmd.Flags().StringVarP(&batchUUID, "batchUUID", "b", "", "provide the batchUUID of the containers")
	snaplistCmd.Flags().StringVarP(&agentUUID, "agentUUID", "a", "", "provide the AgentUUID of the container")
	snaplistCmd.Flags().StringVarP(&containerUUID, "containerUUID", "c", "", "provide the containerUUID of the container")
	snapdestroyCmd.Flags().StringVarP(&batchUUID, "batchUUID", "b", "", "provide the batchUUID of the containers")
	snapdestroyCmd.Flags().StringVarP(&agentUUID, "agentUUID", "a", "", "provide the AgentUUID of the container")
	snapdestroyCmd.Flags().StringVarP(&containerUUID, "containerUUID", "c", "", "provide the containerUUID of the container")
	snapdestroyCmd.Flags().StringVarP(&snapname, "snapname", "s", "", "provide the name of the snapshot to be deleted- e.g: snap0 or \"*\" for destroy all snapshots. QUOTES are mandatory for *")

}
