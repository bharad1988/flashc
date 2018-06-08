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

// agupdateCmd represents the agupdate command
var agupdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update agent info",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: AgentUpdate,
}
var (
	uuid string
	ip   string
)

// AgentUpdate updates the IP address of Agent
func AgentUpdate(cmd *cobra.Command, args []string) {
	url := "http://" + ControllerNode + ":" + common.ControllerPort + "/agent/update"
	var agent common.AgentInfo
	agent.UUID = uuid
	agent.AgentIP = ip
	bodyjson, err := json.Marshal(agent)
	if err != nil {
		panic(err)
	}
	content := bytes.NewReader(bodyjson)

	//bodyType := "application/json; charset=UTF-8"
	resp, err := http.Post(url, common.JSONBodyType, content)

	var message string
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		panic(err)
	}
	fmt.Println(message)
	/*
		ipaddr := net.ParseIP(ip).To4()
		_, err := net.Dial("tcp", ip+":"+common.AgentPort)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(ipaddr)
		fmt.Println(uuid)
	*/
}
func init() {
	agentCmd.AddCommand(agupdateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// agupdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// agupdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	agupdateCmd.Flags().StringVarP(&ip, "ipaddr", "i", "", "Enter the new IP address")
	agupdateCmd.Flags().StringVarP(&uuid, "uuid", "u", "", "Enter the uuid of the agent to be updated")

}
