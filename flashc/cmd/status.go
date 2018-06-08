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

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: AgentStatus,
}

// AgentStatus gives the status of the Agents
func AgentStatus(cmd *cobra.Command, args []string) {
	//fmt.Println("Agent status is being called")
	err := checkController(ControllerNode)
	if err != nil {
		fmt.Printf("\n\n\nCheck if Controller is running on the node with IP - " + ControllerNode + "\n")
		fmt.Println("Commands cannot be executed on that node right now")
		fmt.Printf("use -C option to provide a different controller\n\n\n")
		return
	}

	url := "http://" + ControllerNode + ":" + common.ControllerPort + "/agent/status"
	//bodyType := "application/json; charset=UTF-8"
	resp, err := http.Post(url, "application/json; charset=UTF-8", nil)
	var agents []common.AgentInfo
	err = json.NewDecoder(resp.Body).Decode(&agents)
	if err != nil {
		panic(err)
	}
	fmt.Println("Agent Status")
	fmt.Println("IP			UUID					status		Hostname")
	for _, agent := range agents {
		fmt.Printf("%v		%v	%v		%v\n", agent.AgentIP, agent.UUID, agent.Status, agent.Hostname)
	}

}

func init() {
	//RootCmd.AddCommand(statusCmd)
	agentCmd.AddCommand(statusCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
