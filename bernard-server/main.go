package main

import (
	"log"
	"time"
)

type Config struct {
	Key        string
	ListenAddr string
}

func main() {
	t := make(map[string][]time.Time)
	_ = len(t)
	server := Server{
		ThrottleList: ThrottleList{connAttempts: make(map[string][]time.Time)},
	}
	err := server.Init("./sample.toml")
	if err != nil {
		log.Fatalln(err)
	}
}
