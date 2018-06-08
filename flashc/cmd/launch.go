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
	"gopkg.in/lxc/go-lxc.v2"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: CreateContainers,
}

// SuperCon is the structure that holds the json input to create lxc containers

// CreateContainers - takes the input from user and issues http Post to controller
func CreateContainers(cmd *cobra.Command, args []string) {
	err := checkController(ControllerNode)
	if err != nil {
		fmt.Printf("\n\n\nCheck if Controller is running on the node with IP - " + ControllerNode + "\n")
		fmt.Println("Commands cannot be executed on that node")
		fmt.Printf("use -C option to preovide a different controller\n\n\n")
		return
	}
	// make call to controller.
	url := "http://" + ControllerNode + ":" + common.ControllerPort + "/container/create"
	bodyType := "application/json; charset=UTF-8"
	fmt.Println(url)

	// marshal the structure to json format
	bodyjson, err := json.Marshal(lCon)
	if err != nil {
		panic(err)
	}

	// convert json to io bytes
	content := bytes.NewReader(bodyjson)
	resp, err := http.Post(url, bodyType, content)
	err = json.NewDecoder(resp.Body).Decode(&respStatus)
	if err != nil {
		panic(err)
	}
	fmt.Println(respStatus)
}

func init() {
	//RootCmd.AddCommand(createCmd)
	containerCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&lCon.Con.Lxcpath, "lxcpath", "l", lxc.DefaultConfigPath(), "Enter lxc location: By default, default config path will be taken")
	createCmd.Flags().StringVarP(&lCon.Con.Template, "template", "t", "download", "Template to use")
	createCmd.Flags().StringVarP(&lCon.Con.Arch, "arch", "a", "amd64", "Template to use")
	createCmd.Flags().StringVarP(&lCon.Con.Distro, "distro", "d", "ubuntu", "ubuntu or fedora or etc....")
	createCmd.Flags().StringVarP(&lCon.Con.Release, "release", "r", "precise", "precise, trusty ...")
	createCmd.Flags().StringVarP(&lCon.Con.Name, "batchname", "b", "testcon", "name of the container batch")
	createCmd.Flags().StringVarP(&lCon.Con.Verbose, "verbose", "v", "false", "verbose parameter")
	createCmd.Flags().StringVarP(&lCon.Con.Flush, "flush", "f", "false", "requires flush or not")
	createCmd.Flags().StringVarP(&lCon.Con.Backend, "store", "s", "", "Backend store- btrfs supported right now - set it to 1")
	createCmd.Flags().StringVarP(&lCon.Con.Validation, "validation", "e", "flase", "Validation required or not ")
	createCmd.Flags().StringVarP(&lCon.Count, "count", "c", "1", "Number of containers")

}
