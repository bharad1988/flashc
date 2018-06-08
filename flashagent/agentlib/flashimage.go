package agentlib

import (
	"fmt"
	"log"

	"github.com/bharad1988/flashc/common"
	"gopkg.in/mgo.v2/bson"
)

// LocalImage is the local struct for serving RPC
type LocalImage common.SuperImage

// ImageBuild build the image locally
func (image *LocalImage) ImageBuild(ig *LocalImage, reply *common.ConImage) error {
	// Get the container UUID
	// Get the container path
	// Get the check if snapshot is given
	// check for BackendStore
	// rsync the rootfs to given path for storing image
	// Thats it !!!
	mgs := common.MongoServer{
		Server:     "127.0.0.1",
		Port:       "27017",
		DB:         "cdata",
		Collection: "info",
	}
	mSession, err := common.MongoConnect(&mgs)
	if err != nil { // if mongo server not running on controller throw this error
		return err
	}

	defer mSession.Close()
	mCursor := mSession.DB(mgs.DB).C(mgs.Collection)
	var conts []common.ConConfig
	var cont common.ConConfig
	query := bson.M{"uuid": ig.ConUUID}
	err = mCursor.Find(query).All(&conts)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	cont = conts[0]

	ig.ConPath = cont.Lxcpath + "/" + cont.Name + "/"
	if ig.Snap != "" {
		ig.ConPath = ig.ConPath + "/snaps/" + ig.Snap + "/"
	}

	destPath := ig.Image.Path + "/" + ig.Image.Name
	// copy the rootfs to the image path , as well config file
	cmdName := "mkdir"
	cmdArgs := []string{destPath}
	err = CommandExec(cmdName, cmdArgs)

	cmdName = "rsync"
	rootfs := ig.ConPath + "rootfs"
	destPathRootfs := destPath + "/"
	cmdArgs = []string{"-zvra", rootfs, destPathRootfs}
	err = CommandExec(cmdName, cmdArgs)
	if err != nil {
		log.Print(err)
		return err
	}
	configFile := ig.ConPath + "/config"
	destPathConfig := destPath + "/config"
	cmdArgs = []string{"-zvra", configFile, destPathConfig}
	err = CommandExec(cmdName, cmdArgs)
	if err != nil {
		log.Print(err)
		return err
	}
	reply.Name = ig.Image.Name
	reply.Path = destPath
	reply.UUID, err = common.NewUUID()
	if err != nil {
		log.Print(err)
		return err
	}
	imgCursor := mSession.DB("node").C("images")
	err = imgCursor.Insert(reply)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

// ImageList lists out the images in local agent
func (image *LocalImage) ImageList(ig *LocalImage, reply *[]common.ConImage) error {
	var images []common.ConImage

	mgs := common.MongoServer{
		Server:     "127.0.0.1",
		Port:       common.MongoAgentPort,
		DB:         "node",
		Collection: "images",
	}
	mSession, err := common.MongoConnect(&mgs)
	if err != nil { // if mongo server not running on controller throw this error
		return err
	}

	defer mSession.Close()
	mCursor := mSession.DB(mgs.DB).C(mgs.Collection)
	err = mCursor.Find(nil).All(&images)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println(images)
	*reply = images
	return nil

}
