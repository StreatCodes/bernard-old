package main

import (
	"log"
)

type Config struct {
	Key                 string
	ListenAddr          string
	AuthAttemptsAllowed int
	AuthTimeout         int
}

func main() {
	server, err := NewServer("./sample.toml")
	if err != nil {
		log.Fatalln(err)
	}

	err = server.Listen()
	if err != nil {
		log.Fatalln(err)
	}
}
