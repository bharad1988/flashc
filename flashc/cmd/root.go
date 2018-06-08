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
	"fmt"
	"net"
	"os"

	"gopkg.in/mgo.v2/bson"

	"github.com/bharad1988/flashc/common"
	"github.com/spf13/cobra"
)

//var cfgFile string

// ControllerNode is the DNS name or IP address of Controller Node
var ControllerNode string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "flashc",
	Short: "flashc is the command line client for controlling flashc cluster",
	Long: `flashc
	This is a command line REST client to communicate with controller
	This command is used to manage the flashc cluster with cli tools
	Use flashc -h for more detailed help`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func checkController(cn string) error {
	_, err := net.Dial("tcp", cn+":"+common.ControllerPort)
	if err != nil {
		return err
	}
	mgs := common.MongoServer{
		Server:     cn,
		Port:       "27017",
		DB:         "ctlr",
		Collection: "info",
	}
	mSession, err := common.MongoConnect(&mgs)
	if err != nil {
		//panic(err)
		return err
	}
	defer mSession.Close()
	//var AgNode AgentInfo
	mCursor := mSession.DB(mgs.DB).C(mgs.Collection)

	var m bson.M
	err = mCursor.Find(nil).One(&m)
	if err != nil {
		return err
	}
	//var uuid string
	//uuid := reflect.ValueOf(m["uuid"]).String()

	return nil
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)
	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.flashc.yaml)")
	RootCmd.PersistentFlags().StringVarP(&ControllerNode, "controllernode", "C", "127.0.0.1", "Provide the IP/DNS name of Controller Node")
	//	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
/*
func initConfig() {
	//This is just an empty set
}
*/
