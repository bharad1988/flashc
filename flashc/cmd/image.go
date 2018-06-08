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

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("image called")
	},
}

// contstartCmd represents the contstart command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build an flashc image from given container ",
	Long:  `BUild an image for flashc from a container or snapshot`,
	Run:   ImgBuild,
}

// ImgBuild is the function that builds image
func ImgBuild(cmd *cobra.Command, args []string) {
	fmt.Println("Build Image called")
	fmt.Println("Container start command issued")
	err := checkController(ControllerNode)
	if err != nil {
		fmt.Printf("\n\n\nCheck if Controller is running on the node with IP - " + ControllerNode + "\n")
		fmt.Println("Commands cannot be executed on that node")
		fmt.Printf("use -C option to preovide a different controller\n\n\n")
		return
	}
	if containerUUID == "" || agentUUID == "" || imagename == "" {
		fmt.Println("Please provide the required parameters- container UUID, agent UUID and imagename")
	}
	// make call to controller.
	url := "http://" + ControllerNode + ":" + common.ControllerPort + "/image/build"
	//bodyType := "application/json; charset=UTF-8"
	fmt.Println(url)
	superimag.ConUUID = containerUUID
	superimag.AgentUUID = agentUUID
	superimag.Snap = snapname
	superimag.Image.Name = imagename

	// marshal the structure to json format
	bodyjson, err := json.Marshal(superimag)
	if err != nil {
		panic(err)
	}

	// convert json to io bytes
	content := bytes.NewReader(bodyjson)
	resp, err := http.Post(url, common.JSONBodyType, content)
	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(resp.Body).Decode(&httpResp)
	if err != nil {
		panic(err)
	}
	fmt.Println(httpResp)
}

func init() {
	RootCmd.AddCommand(imageCmd)
	imageCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVarP(&agentUUID, "agentUUID", "a", "", "[Required] provide the AgentUUID of the container")
	buildCmd.Flags().StringVarP(&containerUUID, "containerUUID", "c", "", "[Required] provide the containerUUID of the container")
	buildCmd.Flags().StringVarP(&snapname, "snapname", "s", "", "[Optional] provide the name of the snapshot to be deleted- e.g: snap0 or \"*\" for destroy all snapshots. QUOTES are mandatory for *")
	buildCmd.Flags().StringVarP(&imagename, "imagename", "i", "", "[Required] provide the name of the image. The image will be saved in this name")
}
