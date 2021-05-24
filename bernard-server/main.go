package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	toml "github.com/pelletier/go-toml"
	"github.com/streatcodes/bernard"
)

func handleConnection(config Config, conn net.Conn) {
	err := AuthClient(config, conn)
	if err != nil {
		fmt.Printf("Failed to auth new connection: %s\n", err)
		return
	}

	fmt.Println("New client authenticated")

	for {
		checkResult := bernard.CheckResult{}
		decoder := json.NewDecoder(conn)
		err := decoder.Decode(&checkResult)
		if err != nil {
			fmt.Printf("Error decoding message: %s - closing connection\n", err)
			return
		}

		fmt.Printf("New message: %+v\n", checkResult)
	}
}

type Config struct {
	Key string
}

func main() {
	//Read config
	configPath := "./sample.toml"
	f, err := os.OpenFile(configPath, os.O_RDONLY, 0755)
	if err != nil {
		log.Fatalf("Error opening %s: %s\n", configPath, err)
	}

	var config Config
	dec := toml.NewDecoder(f)
	err = dec.Decode(&config)
	if err != nil {
		log.Fatalf("Error decoding %s: %s\n", configPath, err)
	}

	listenAddr := "localhost:8888"
	fmt.Printf("Listening on %s\n", listenAddr)
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	//Accept incoming connections
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(config, conn)
	}
}
