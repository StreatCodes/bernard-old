package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

type CheckSettings struct {
	Name     string
	Command  string
	Args     []string
	Stdin    string
	Env      []string
	Dir      string
	Interval time.Duration
	Timeout  time.Duration
}

type CheckScheduler struct {
	Checks []CheckSettings
}

func (s *CheckScheduler) Start() {
	for _, check := range s.Checks {
		go checkRunner(check)
	}
}

func checkRunner(checkSettings CheckSettings) {
	fmt.Printf("Scheduling check %s - interval %s\n", checkSettings.Name, checkSettings.Interval)
	ticker := time.NewTicker(checkSettings.Interval)

	for {
		t := <-ticker.C
		fmt.Printf("%s - Running check %s\n", t, checkSettings.Name)
		runCheck(checkSettings)
	}
}

func runCheck(checkSettings CheckSettings) {
	ctx, cancel := context.WithTimeout(context.Background(), checkSettings.Timeout)
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
		fmt.Println(err)
	}

	io.Copy(os.Stdout, stdoutBuf)

	//TODO write to result channel
}
