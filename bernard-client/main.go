package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	toml "github.com/pelletier/go-toml"
	"github.com/streatcodes/bernard"
)

type Config struct {
	ParentNode NodeConfig
	Checks     map[string]CheckSettings
}

type NodeConfig struct {
	Address string
	Key     string
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

	//Initialise the check scheduler
	parentNodeChan := make(chan bernard.CheckResult)
	StartScheduler(parentNodeChan, config.Checks)

	//Connect to parent node
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", config.ParentNode.Address)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %s", config.ParentNode.Address, err)
	}
	defer conn.Close()

	//Authenticate with server
	err = authToServer(config, conn)
	if err != nil {
		log.Fatalf("Failed to init auth with server: %s\n", err)
	}

	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)

	var authResult bernard.AuthResult
	err = decoder.Decode(&authResult)
	if err != nil {
		log.Fatalf("Failed to decode auth result from server: %s\n", err)
	}

	if !authResult.Success {
		log.Fatalln("Failed to authenticate with server")
	}

	fmt.Println("Authenticated with server")

	//Listen on the channel and send check results upstream
	for checkResult := range parentNodeChan {
		// fmt.Printf("Sending %+v\n", checkResult)
		err := encoder.Encode(checkResult)
		if err != nil {
			log.Fatalf("Error writing result to parent node: %s\n", err)
		}
	}
}
