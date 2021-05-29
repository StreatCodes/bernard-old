package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/streatcodes/bernard"
)

type CheckSettings struct {
	Name     string
	Command  string
	Args     []string
	Stdin    string
	Env      []string
	Dir      string
	Interval int64
	Timeout  int64
}

func StartScheduler(parentNodeChan chan bernard.CheckResult, checkSettings map[string]CheckSettings) {
	for _, check := range checkSettings {
		go checkRunner(parentNodeChan, check)
	}
}

func checkRunner(parentNodeChan chan bernard.CheckResult, checkSettings CheckSettings) {
	fmt.Printf("Scheduling check %s - interval %d seconds\n", checkSettings.Name, checkSettings.Interval)

	ticker := time.NewTicker(time.Duration(checkSettings.Interval) * time.Second)

	for {
		t := <-ticker.C
		fmt.Printf("%s - Running check %s\n", t, checkSettings.Name)
		parentNodeChan <- runCheck(checkSettings)
	}
}

func runCheck(checkSettings CheckSettings) bernard.CheckResult {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(checkSettings.Timeout)*time.Second)
	defer cancel()

	//Setup check command
	command := exec.CommandContext(ctx, checkSettings.Command, checkSettings.Args...)

	command.Stdin = bytes.NewBufferString(checkSettings.Stdin)
	command.Env = checkSettings.Env
	command.Dir = checkSettings.Dir
	var stdoutBuf bytes.Buffer
	command.Stdout = &stdoutBuf

	//Exec check command
	err := command.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting command (%s) - %s", checkSettings.Name, err)

		return bernard.CheckResult{
			Status: -1,
			Output: stdoutBuf.Next(4096),
		}
	}

	//Get result
	err = command.Wait()
	if err != nil {
		exitCode := -1
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
		fmt.Fprintf(os.Stderr, "Command exited with %d (%s) - %s", exitCode, checkSettings.Name, err)
		return bernard.CheckResult{
			Status: exitCode,
			Output: stdoutBuf.Next(4096),
		}
	}

	return bernard.CheckResult{
		Status: 0,
		Output: stdoutBuf.Next(4096),
	}
}
