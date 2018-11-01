package repolib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/bharad1988/btrfslib"
	"github.com/bharad1988/flashc/common"
)

// LocalImage is the local image type for repo service
type LocalImage common.ConImage

// GetGroup gives the group for the username
func GetGroup(username string) (string, error) {
	cmd := "id"
	cmdArgs := []string{"-g", "-n", username}
	cmdOut, err := exec.Command(cmd, cmdArgs...).Output()
	if err != nil {
		log.Println(err)
		return "", err
	}
	s := string(cmdOut[:])

	return strings.TrimSpace(s), nil
}

// ChownFile changes the ownership of the file with given username and group
func ChownFile(path, user, group string) error {
	cmdName := "chown"
	usrgrp := user + ":" + group
	cmdArgs := []string{usrgrp, path}
	err := common.CommandExec(cmdName, cmdArgs)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GetImage gets the image details from repository
func GetImage(repoRoot string, repoImage *common.ConImage) error {
	var snaplist []int

	// Declare the btrfs subvolume and add its path to the object
	btrfsVol := new(btrfslib.SubVolume)
	btrfsVol.Path = repoImage.User + "/" + repoImage.Name + "/latest/"

	//var versions []os.FileInfo
	// imagePath is the path of the image inside the user directory
	imagePath := repoRoot + repoImage.User + "/" + repoImage.Name
	versions, err := ioutil.ReadDir(imagePath + "/versions/")
	if err != nil {
		log.Println(err)
	}

	grp, err := GetGroup(repoImage.User)
	if err != nil {
		log.Println(err)
		return err
	}

	if len(versions) > 0 { // Image exists with versions
		// set the version that has to be created
		//repoImage.Version = strconv.FormatInt(int64(snaplist[len(snaplist)-1]), 10)

		for _, version := range versions {
			//fmt.Println(version.Name())
			i64, err := strconv.ParseInt(version.Name(), 10, 32)
			if err != nil {
				log.Println(err)
			}
			i := int(i64)
			snaplist = append(snaplist, i)
		}
		fmt.Print(len(snaplist))
		nextVersion := snaplist[len(snaplist)-1] + 1
		fmt.Print(nextVersion)
		repoImage.Version = strconv.FormatInt(int64(nextVersion), 10)

		// Get the latest version for this image
		sort.Ints(snaplist)
		//fmt.Print(snaplist[len(snaplist)-1])
		repoImage.Path = imagePath
	} else { // If this is the first time image is getting created, then set the version to 0 and prepare the directories
		repoImage.Version = "0"
		err := os.MkdirAll(imagePath, 0755)
		if err != nil {
			log.Println(err)
			return err
		}
		// chown the imagepath dir with user and group
		err = ChownFile(imagePath, repoImage.User, grp)

		// create a btrfs subvolume to hold the latest image
		// latestImageDir is created here
		btrfsVol.Create(repoRoot)
		// change the ownership to the user who is hosting it
		latestImageDir := imagePath + "/latest"
		err = ChownFile(latestImageDir, repoImage.User, grp)

		// create the directory to hold the versions
		err = os.MkdirAll(imagePath+"/versions/", 0755)
		if err != nil {
			log.Println(err)
			return err
		}
		//change the ownership of the versions directory
		versionDir := imagePath + "/versions/"
		err = ChownFile(versionDir, repoImage.User, grp)
		if err != nil {
			log.Println(err)
			return err
		}

		repoImage.Path = imagePath
	}
	return nil
}

// GetRepoPath gives the path for repository root
func GetRepoPath() (string, error) {
	var repos common.RepoService
	jsonData, err := ioutil.ReadFile("/home/user/repos.json")
	if err != nil {
		log.Println(err)
		return "", err
	}

	err = json.Unmarshal(jsonData, &repos)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return repos.Path, nil
}

// ImageCreate creates an image in the repo server
func (locimage *LocalImage) ImageCreate(li *LocalImage, reply *common.ConImage) error {
	// check image exists or not
	// if exists check latest Version returns latest Version to client
	// Creates directories for latest image and versions
	repoRoot, err := GetRepoPath()
	if err != nil {
		log.Print(err)
		return err
	}

	reply.Name = li.Name
	reply.User = li.User

	log.Print(reply)

	// Get the image data and the latest version that has to be created
	err = GetImage(repoRoot, reply)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// PostSync is a post rsync call to take snapshot for the version
func (locimage *LocalImage) PostSync(li *LocalImage, reply *common.ConImage) error {
	repoRoot, err := GetRepoPath()
	if err != nil {
		log.Print(err)
		return err
	}
	btrfsVol := new(btrfslib.SubVolume)
	btrfsVol.Path = li.User + "/" + li.Name + "/latest/"
	snapShotPath := li.User + "/" + li.Name + "/versions/" + li.Version

	// create snapshot object to create a snapshot of the previous version
	btrfsSnapObj := new(btrfslib.Snap)
	btrfsSnapObj.Readonly = true
	btrfsSnapObj.DestVol.Path = snapShotPath
	btrfsVol.Snapshot(repoRoot, btrfsSnapObj)
	return nil
}
