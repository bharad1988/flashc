package agentlib

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bharad1988/flashc/common"

	"gopkg.in/lxc/go-lxc.v2"
	"gopkg.in/mgo.v2/bson"
)

//LocalContainers is a local struct type to define the methods that are local to this library
type LocalContainers common.SuperCon

//var lcs common.SuperCon

//Create ALL container methods that will run will be a part of this library
//Create container on the agent node
func (cont *LocalContainers) Create(ct *LocalContainers, reply *string) error {
	lxcpath := ct.Con.Lxcpath
	template := ct.Con.Template
	distro := ct.Con.Distro
	release := ct.Con.Release
	arch := ct.Con.Arch
	// Set btrfs as backend store

	// convert below variables to bool
	verbose, _ := strconv.ParseBool(ct.Con.Verbose)
	flush, _ := strconv.ParseBool(ct.Con.Flush)
	validation, _ := strconv.ParseBool(ct.Con.Validation)

	total, err := strconv.ParseInt(ct.Count, 10, 64)
	begin, err := strconv.ParseInt(ct.Start, 10, 64)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}

	//=======================================================================================
	//Run loop over the number of containers that have to be created in this batch
	// BEGIN
	//=======================================================================================
	for index := int(begin); index < int(total); index++ {
		fmt.Println(index)
		locon := common.ConConfig{}
		locon = ct.Con

		// Need to reconstruct the name for each container
		locon.Name = ct.Con.Name + "-" + strconv.FormatInt(int64(index), 10)
		c, err := lxc.NewContainer(locon.Name, lxcpath)
		if err != nil {
			//*reply = err.Error()
			log.Printf("ERROR: %s\n", err.Error())
			//fmt.Println(*reply)
			return err
		}

		log.Printf("Creating container...\n")
		if verbose {
			c.SetVerbosity(lxc.Verbose)
			fmt.Println(lxc.Verbose)
		}

		options := lxc.TemplateOptions{
			Template:             template,
			Distro:               distro,
			Release:              release,
			Arch:                 arch,
			FlushCache:           flush,
			DisableGPGValidation: validation,
		}
		if ct.Con.Backend == "1" {
			var backend lxc.BackendStore
			backend = 1
			options.Backend = backend
		}

		fmt.Println(options)
		if err = c.Create(options); err != nil {
			log.Printf("ERROR: %s\n", err.Error())
			//fmt.Println(*reply)
			return err
		}

		uuid, err := common.NewUUID()
		locon.UUID = uuid
		fmt.Println("Name : " + locon.Name + "       UUID : " + locon.UUID)
		agent, err := GetAgentParam()
		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
			//*reply = err.Error()
			return err
		}
		locon.AgentUUID = agent.UUID

		// need to maintain the each container's data in local DB.
		err = MgoInsert(&locon)
		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
			//*reply = err.Error()
			return err
		}
	}
	//=======================================================================================
	// END of loop
	//=======================================================================================

	*reply = "container(s) successfully created"
	return nil
}

// Echo is a test function
func (cont *LocalContainers) Echo(ct *LocalContainers, reply *bool) error {
	fmt.Println("This call was called ")
	fmt.Println(ct)
	*reply = true
	return nil
}

// StartContainers method starts the containers on local agents
func (cont *LocalContainers) StartContainers(ct *LocalContainers, reply *string) error {
	// connect to agent MongoServer
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
	fmt.Println(ct.Con.BatchUUID)
	var conts []common.ConConfig
	err = mCursor.Find(bson.M{"batchuuid": ct.Con.BatchUUID}).All(&conts)
	if err != nil { // if mongo server not running on controller throw this error
		return err
	}
	for _, cont := range conts {
		c, err := lxc.NewContainer(cont.Name, cont.Lxcpath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = c.Start()
		err = c.SetMemoryLimit(common.ContMemLimit)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = c.SetSoftMemoryLimit(common.ContMemSoftLimit)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	fmt.Println(conts)
	*reply = "successfully started containers"
	return nil
}

// StopContainers method stops the containers on local agents
func (cont *LocalContainers) StopContainers(ct *LocalContainers, reply *string) error {
	// connect to agent MongoServer
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
	fmt.Println(ct.Con.BatchUUID)
	var conts []common.ConConfig
	err = mCursor.Find(bson.M{"batchuuid": ct.Con.BatchUUID}).All(&conts)
	if err != nil { // if mongo server not running on controller throw this error
		return err
	}
	for _, cont := range conts {
		c, err := lxc.NewContainer(cont.Name, cont.Lxcpath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = c.Stop()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	fmt.Println(conts)
	*reply = "successfully stopped containers"
	return nil
}

// DestroyContainers method destroys the containers on local agents
func (cont *LocalContainers) DestroyContainers(ct *LocalContainers, reply *string) error {
	// connect to agent MongoServer
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
	statCursor := mSession.DB(mgs.DB).C("stat")
	fmt.Println(ct.Con.BatchUUID)
	var conts []common.ConConfig
	err = mCursor.Find(bson.M{"batchuuid": ct.Con.BatchUUID}).All(&conts)
	if err != nil { // if mongo server not running on controller throw this error
		return err
	}
	if ct.Con.BatchUUID == "" {
		err = mCursor.Find(bson.M{"uuid": ct.Con.UUID}).All(&conts)
		if err != nil { // if mongo server not running on controller throw this error
			return err
		}
	}

	for _, cont := range conts {
		c, err := lxc.NewContainer(cont.Name, cont.Lxcpath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if ct.DC.WithSnapShots == "true" {
			log.Println("Delete snapshots as well. Delete containers")
			err = c.Destroy() // attempt to destroy by default . This is to handle the cases of those which have no snapshots
			if err != nil {
				fmt.Println(err)
				err = c.DestroyWithAllSnapshots()
				if err != nil {
					fmt.Println(err)
					return err
				}
			}
		}
		if ct.DC.WithSnapShots == "false" {
			log.Println("No snapshots. Just Delete containers")
			err = c.Destroy()
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
		err = mCursor.Remove(bson.M{"uuid": cont.UUID})
		if err != nil {
			fmt.Printf("remove fail %v\n", err)
			return err
		}
		err = statCursor.Remove(bson.M{"contuuid": cont.UUID})
		if err != nil {
			fmt.Printf("remove fail %v\n", err)
			return err
		}

	}
	fmt.Println(conts)
	*reply = "successfully destroyed containers"
	return nil

}

// TestDestroy to verify destroy options
func (cont *LocalContainers) TestDestroy(ct *LocalContainers, reply *string) error {
	// connect to agent MongoServer
	mgs := common.MongoServer{
		Server:     "127.0.0.1",
		Port:       common.MongoAgentPort,
		DB:         "cdata",
		Collection: "info",
	}
	mSession, err := common.MongoConnect(&mgs)
	if err != nil { // if mongo server not running on controller throw this error
		return err
	}

	defer mSession.Close()
	mCursor := mSession.DB(mgs.DB).C(mgs.Collection)
	fmt.Println(ct.Con.BatchUUID)
	var conts []common.ConConfig
	err = mCursor.Find(bson.M{"batchuuid": ct.Con.BatchUUID}).All(&conts)
	if err != nil { // if mongo server not running on controller throw this error
		return err
	}
	for _, cont := range conts {
		c, err := lxc.NewContainer(cont.Name, cont.Lxcpath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		//c.RestoreSnapshot(snapshot, name)

		/*
			snaps,err:= c.Snapshots()
			c.RestoreSnapshot(snapshot, name)
			c.Restore(opts)
					c.DestroyWithAllSnapshots()
					c.DestroyAllSnapshots()
					c.DestroySnapshot(snapshot)

				c.SetCgroupItem(key, value)
				c.SetSoftMemoryLimit(limit)
				c.SetMemoryLimit(limit)
		*/
		err = c.Destroy()
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = mCursor.Remove(bson.M{"uuid": cont.UUID})
		if err != nil {
			fmt.Printf("remove fail %v\n", err)
			return err
		}

	}
	fmt.Println(conts)
	*reply = "successfully destroyed containers"
	return nil
}

// StatContainers method destroys the containers on local agents
func (cont *LocalContainers) StatContainers(ct *LocalContainers, stats *[]common.ContStat) error {
	// connect to agent MongoServer
	mgs := common.MongoServer{
		Server:     "127.0.0.1",
		Port:       common.MongoAgentPort,
		DB:         "cdata",
		Collection: "info",
	}
	mSession, err := common.MongoConnect(&mgs)
	if err != nil { // if mongo server not running on controller throw this error
		return err
	}

	defer mSession.Close()
	mCursor := mSession.DB(mgs.DB).C(mgs.Collection)
	statCursor := mSession.DB(mgs.DB).C("stat")
	fmt.Println(ct.Con.BatchUUID)
	var conts []common.ConConfig

	err = mCursor.Find(bson.M{"batchuuid": ct.Con.BatchUUID}).All(&conts)
	if err != nil { // if mongo server not running on controller throw this error
		return err
	}
	fmt.Println(conts)

	if ct.Con.BatchUUID == "" {
		err = mCursor.Find(bson.M{"uuid": ct.Con.UUID}).All(&conts)
		if err != nil { // if mongo server not running on controller throw this error
			return err
		}
	}
	fmt.Println("Get the stat of this cont")
	var fillstats []common.ContStat
	for _, cont := range conts {
		err = statCursor.Find(bson.M{"contuuid": cont.UUID}).All(&fillstats)
		if err != nil { // if mongo server not running on controller throw this error
			return err
		}
		*stats = append(*stats, fillstats...)
	}
	//fmt.Println(fillstats[0].Memory)
	//fmt.Println(stats)
	return nil
}

// SnapContainers method destroys the containers on local agents
func (cont *LocalContainers) SnapContainers(ct *LocalContainers, reply *string) error {
	// connect to agent MongoServer
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
	//statCursor := mSession.DB(mgs.DB).C("stat")
	fmt.Println(ct.Con.BatchUUID)
	var conts []common.ConConfig
	err = mCursor.Find(bson.M{"batchuuid": ct.Con.BatchUUID}).All(&conts)
	if err != nil { // if mongo server not running on controller throw this error
		return err
	}
	if ct.Con.BatchUUID == "" {
		err = mCursor.Find(bson.M{"uuid": ct.Con.UUID}).All(&conts)
		if err != nil { // if mongo server not running on controller throw this error
			return err
		}
	}

	for _, cont := range conts {
		c, err := lxc.NewContainer(cont.Name, cont.Lxcpath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		snap, err := c.CreateSnapshot()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(snap)
	}
	fmt.Println(conts)
	*reply = "successfully snapped containers"
	return nil

}

// SnapList method destroys the containers on local agents
func (cont *LocalContainers) SnapList(ct *LocalContainers, snaps *[]common.SnapShot) error {
	// connect to agent MongoServer
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
	//statCursor := mSession.DB(mgs.DB).C("stat")
	fmt.Println(ct.Con.UUID)
	var conts []common.ConConfig
	err = mCursor.Find(bson.M{"batchuuid": ct.Con.BatchUUID}).All(&conts)
	if err != nil { // if mongo server not running on controller throw this error
		return err
	}
	if ct.Con.BatchUUID == "" {
		err = mCursor.Find(bson.M{"uuid": ct.Con.UUID}).All(&conts)
		if err != nil { // if mongo server not running on controller throw this error
			return err
		}
	}

	for _, cont := range conts {
		c, err := lxc.NewContainer(cont.Name, cont.Lxcpath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		snapshots, err := c.Snapshots()
		if err != nil {
			fmt.Println(err)
			return err
		}

		//fmt.Println(snaps)

		for _, snapshot := range snapshots {
			fmt.Println(snapshot)
			var snap common.SnapShot
			snap.CommentPath = snapshot.CommentPath
			snap.Name = snapshot.Name
			snap.Path = snapshot.Path
			snap.TimeStamp = snapshot.Timestamp
			snap.ContainerName = cont.Name
			*snaps = append(*snaps, snap)
		}
		fmt.Println(snaps)
	}
	return nil
}

// SnapDestroy method destroys the containers on local agents
func (cont *LocalContainers) SnapDestroy(ct *LocalContainers, reply *string) error {
	// connect to agent MongoServer
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
	//statCursor := mSession.DB(mgs.DB).C("stat")
	fmt.Println(ct.Con.BatchUUID)
	var conts []common.ConConfig
	err = mCursor.Find(bson.M{"batchuuid": ct.Con.BatchUUID}).All(&conts)
	if err != nil { // if mongo server not running on controller throw this error
		return err
	}
	if ct.Con.BatchUUID == "" {
		err = mCursor.Find(bson.M{"uuid": ct.Con.UUID}).All(&conts)
		if err != nil { // if mongo server not running on controller throw this error
			return err
		}
	}

	for _, cont := range conts {
		c, err := lxc.NewContainer(cont.Name, cont.Lxcpath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(ct.Snaps.Name)
		if ct.Snaps.Name == "*" {
			fmt.Println(ct.Snaps.Name)
			err = c.DestroyAllSnapshots()
			if err != nil {
				fmt.Println(err)
				return err
			}
		} else {
			var snapname lxc.Snapshot
			snapname.Name = ct.Snaps.Name
			err = c.DestroySnapshot(snapname)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}
	fmt.Println(conts)
	*reply = "successfully destroyed the snapshots"
	return nil

}
