package common

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"reflect"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ErrHandle method handles the error thrown
func ErrHandle(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

// NewUUID generates a random UUID according to RFC 4122
// --- found at go playground ---
func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// MongoConnect connects to mongo Server
func MongoConnect(mgs *MongoServer) (*mgo.Session, error) {
	session, err := mgo.Dial(mgs.Server + ":" + mgs.Port)
	if err != nil {
		panic(err)
	}

	//defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session, err
}

// VerifyAgent verifies agent against the given IP
func VerifyAgent(ip string) (string, error) {
	_, err := net.Dial("tcp", ip+":"+AgentPort)
	if err != nil {
		return "flashagent daemon is not running", err
	}
	mgs := MongoServer{
		Server:     ip,
		Port:       "27017",
		DB:         "node",
		Collection: "info",
	}
	mSession, err := MongoConnect(&mgs)
	if err != nil {
		//panic(err)
		return "failed to connect to Agent Mongo Server", err
	}
	defer mSession.Close()
	//var AgNode AgentInfo
	mCursor := mSession.DB(mgs.DB).C(mgs.Collection)

	var m bson.M
	err = mCursor.Find(nil).One(&m)
	if err != nil {
		return "Failed to find the agent info. Node not yet identified as Agent", err
	}
	//var uuid string
	uuid := reflect.ValueOf(m["uuid"]).String()
	agent := reflect.ValueOf(m)
	fmt.Println(agent)

	return uuid, nil
}

// CommandExec executes a command locally
func CommandExec(cmdName string, cmdArgs []string) error {
	fmt.Println(cmdName)
	fmt.Println(cmdArgs)
	cmd := exec.Command(cmdName, cmdArgs...)

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return err
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("Running the command | %s\n", scanner.Text())
		}
	}()
	//cmd.Run()
	//time.Sleep(30 * time.Second)

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		return err
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		return err
	}

	return nil
}
