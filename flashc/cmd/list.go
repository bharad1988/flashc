// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bharad1988/flashc/common"
	"github.com/spf13/cobra"
)

// sub command container list
// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: containerList,
}

func containerList(cmd *cobra.Command, args []string) {
	err := checkController(ControllerNode)
	if err != nil {
		fmt.Printf("\n\n\nCheck if Controller is running on the node with IP - " + ControllerNode + "\n")
		fmt.Println("Commands cannot be executed on that node")
		fmt.Printf("use -C option to preovide a different controller\n\n\n")
		return
	}

	// connect to all agents that are registered and running
	// get list from cdata
	url := "http://" + ControllerNode + ":" + common.ControllerPort + "/container/list"
	//bodyType := "application/json; charset=UTF-8"
	resp, err := http.Post(url, "application/json; charset=UTF-8", nil)
	var conts []common.ConConfig

	err = json.NewDecoder(resp.Body).Decode(&conts)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Name\t\t\tUUID\t\t\t\t\tState\t\tBatchUUID\t\t\t\tAgentUUID\n")
	for _, cont := range conts {

		fmt.Printf("%v\t\t%v\t%v\t\t%v\t%v\n", cont.Name, cont.UUID, cont.State, cont.BatchUUID, cont.AgentUUID)

	}
	//fmt.Println(conts)

}

func init() {
	containerCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
