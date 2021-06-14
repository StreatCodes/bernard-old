package main

import (
	"log"
)

type Config struct {
	Key                 string
	ListenAddr          string
	DBPath              string
	HTTPListenAddr      string
	AuthAttemptsAllowed int
	AuthTimeout         int
}

func startHttpServer(server *Server) {
	err := server.ListenHTTP()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	server, err := NewServer("./sample.toml")
	if err != nil {
		log.Fatalln(err)
	}

	go startHttpServer(server)

	err = server.Listen()
	if err != nil {
		log.Fatalln(err)
	}
}
