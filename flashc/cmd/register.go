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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bharad1988/flashc/common"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register an Agent node with the controller",
	Long: `flashc register is a command that allows the controller to register the
	available Agents in the network to participate in the flashc cluster
	e.g : `,
	Run: RegisterAgent,
}

var unregisterCmd = &cobra.Command{
	Use:   "unregister",
	Short: "UnRegister an Agent node with the controller",
	Long: `flashc unregister is a command that allows the controller to unregister an Agent
	e.g : `,
	Run: UnRegisterAgent,
}

var agentinfo common.AgentInfo

// RegisterAgent is a fuction that handles register
func RegisterAgent(cmd *cobra.Command, args []string) {
	err := checkController(ControllerNode)
	if err != nil {
		fmt.Printf("\n\n\nCheck if Controller is running on the node with IP - " + ControllerNode + "\n")
		fmt.Println("Commands cannot be executed on that node")
		fmt.Printf("use -C option to preovide a different controller\n\n\n")
	}

	url := "http://" + ControllerNode + ":" + common.ControllerPort + "/agent/register"
	fmt.Println(url)
	fmt.Println(agentinfo.AgentIP)
	bodyjson, err := json.Marshal(agentinfo)
	if err != nil {
		panic(err)
	}
	content := bytes.NewReader(bodyjson)
	resp, err := http.Post(url, common.JSONBodyType, content)
	err = json.NewDecoder(resp.Body).Decode(&respStatus)
	if err != nil {
		panic(err)
	}
	fmt.Println(respStatus)
}

// UnRegisterAgent unregisters the agent
func UnRegisterAgent(cmd *cobra.Command, args []string) {
	err := checkController(ControllerNode)
	if err != nil {
		fmt.Printf("\n\n\nCheck if Controller is running on the node with IP - " + ControllerNode + "\n")
		fmt.Println("Commands cannot be executed on that node")
		fmt.Printf("use -C option to preovide a different controller\n\n\n")
	}

	url := "http://" + ControllerNode + ":" + common.ControllerPort + "/agent/unregister"
	fmt.Println(url)
	fmt.Println(agentinfo.UUID)
	bodyjson, err := json.Marshal(agentinfo)
	if err != nil {
		panic(err)
	}
	content := bytes.NewReader(bodyjson)
	resp, err := http.Post(url, common.JSONBodyType, content)
	err = json.NewDecoder(resp.Body).Decode(&respStatus)
	if err != nil {
		panic(err)
	}
	fmt.Println(respStatus)
}

func init() {
	//RootCmd.AddCommand(registerCmd)
	agentCmd.AddCommand(registerCmd)
	agentCmd.AddCommand(unregisterCmd)
	registerCmd.Flags().StringVarP(&agentinfo.AgentIP, "agentip", "a", "", "Enter the IP address or DNS Name of the agent")
	unregisterCmd.Flags().StringVarP(&agentinfo.UUID, "agentuuid", "u", "", "Enter the UUID of the agent to be unregistered")
}
