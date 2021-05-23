package main

import (
	"log"
	"os"

	toml "github.com/pelletier/go-toml"
)

type Config struct {
	ParentNode NodeConfig
	Checks     map[string]CheckSettings
}

type NodeConfig struct {
	Address string
}

func main() {
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

	StartScheduler(config.Checks)

	noExit := make(chan bool)
	<-noExit

	// var d net.Dialer
	// ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	// defer cancel()

	// conn, err := d.DialContext(ctx, "tcp", "localhost:12345")
	// if err != nil {
	// 	log.Fatalf("Failed to dial: %v", err)
	// }
	// defer conn.Close()

	// if _, err := conn.Write([]byte("Hello, World!")); err != nil {
	// 	log.Fatal(err)
	// }
}
