package main

import (
	"context"
	"encoding/json"
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
	//TODO improve this context
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", config.ParentNode.Address)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %s", config.ParentNode.Address, err)
	}
	defer conn.Close()

	err = authToServer(config, conn)
	if err != nil {
		log.Fatalf("Failed to auth to server: %s\n", err)
	}
	fmt.Println("Authenticated with server")

	//Listen on the channel and send check results upstream
	for checkResult := range parentNodeChan {
		encoder := json.NewEncoder(conn)
		// fmt.Printf("Sending %+v\n", checkResult)
		err := encoder.Encode(checkResult)
		if err != nil {
			//TODO don't fatal here
			log.Fatalf("Error writing result to parent node: %s\n", err)
		}
	}
}
