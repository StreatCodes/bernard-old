package main

import (
	"bytes"
	"context"
	"fmt"
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
	stdoutBuf := bytes.NewBuffer(make([]byte, 4096))
	command.Stdout = stdoutBuf

	stderrBuf := bytes.NewBuffer(make([]byte, 4096))
	command.Stderr = stderrBuf

	//Exec check command
	err := command.Start()
	if err != nil {
		//TODO write error to result channel
		fmt.Println("Error starting cmd")
		fmt.Println(err)
	}

	err = command.Wait()
	if err != nil {
		//TODO write error to result channel
		fmt.Println("Error running cmd")
		// if exitError, ok := err.(*exec.ExitError); ok {
		// 	exictCode := exitError.ExitCode()
		// }
		fmt.Println(err)
	}

	// io.Copy(os.Stdout, stdoutBuf)

	//TODO expand this
	return bernard.CheckResult{
		Status: 0,
	}
}
