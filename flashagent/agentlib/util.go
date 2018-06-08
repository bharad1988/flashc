package agentlib

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

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
