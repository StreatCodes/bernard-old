package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/streatcodes/bernard"
)

type Server struct {
	Config       Config
	ThrottleList ThrottleList
}

func (s *Server) Init(configPath string) error {
	//Read config
	f, err := os.OpenFile(configPath, os.O_RDONLY, 0755)
	if err != nil {
		return fmt.Errorf("error opening %s: %s", configPath, err)
	}

	dec := toml.NewDecoder(f)
	err = dec.Decode(&s.Config)
	if err != nil {
		return fmt.Errorf("error decoding %s: %s", configPath, err)
	}

	//Start server
	fmt.Printf("Listening on %s\n", s.Config.ListenAddr)
	l, err := net.Listen("tcp", s.Config.ListenAddr)
	if err != nil {
		return fmt.Errorf("error starting TCP server on %s: %s", s.Config.ListenAddr, err)
	}
	defer l.Close()

	//Accept incoming connections
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting connection: %s\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()
	client, _, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Printf("Closing failed to extract host from remote address %s: %s\n", addr, err)
		return
	}

	//Close connection is remote is throttled
	fmt.Printf("New connection: %s\n", client)
	isAddrThrottled := s.ThrottleList.IsThrottled(client)
	if isAddrThrottled {
		fmt.Printf("Closing connection address is throttled: %s\n", client)
		return
	}

	//Verify client key
	encoder := gob.NewEncoder(conn)
	err = AuthClient(s.Config, conn)
	if err != nil {
		fmt.Printf("Failed to auth new connection: %s\n", err)
		s.ThrottleList.FailedAttempt(client)

		err = encoder.Encode(bernard.AuthResult{Success: false})
		if err != nil {
			fmt.Printf("Failed to write auth result to client: %s\n", err)
		}
		return
	}

	err = encoder.Encode(bernard.AuthResult{Success: true})
	if err != nil {
		fmt.Printf("Failed to write auth result to client: %s\n", err)
	}
	fmt.Printf("Client authenticated: %s\n", client)

	//Read incoming check results
	decoder := gob.NewDecoder(conn)
	for {
		checkResult := bernard.CheckResult{}
		err := decoder.Decode(&checkResult)
		if err != nil {
			fmt.Printf("Error decoding message: %s - closing connection\n", err)
			return
		}

		fmt.Printf("Check result from %s - status %d:\n%s\n", client, checkResult.Status, checkResult.Output)
	}
}
