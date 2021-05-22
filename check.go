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
	Command string
	Args    []string
	Stdin   []byte
	Env     []string
	Dir     string
	Timeout time.Duration
}

type CheckScheduler struct {
	Checks []CheckSettings
}

func checkScheduler() {
	checkSettings := CheckSettings{
		Command: "ping",
		Args:    []string{"-c 10", "www.google.com"},
		Env:     []string{""},
		Dir:     "",
		Timeout: 5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), checkSettings.Timeout)
	defer cancel()

	//Setup check command
	command := exec.CommandContext(ctx, checkSettings.Command, checkSettings.Args...)

	command.Stdin = bytes.NewReader(checkSettings.Stdin)
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
