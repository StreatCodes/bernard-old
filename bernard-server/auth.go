package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net"
)

func GenerateToken(length int) ([]byte, error) {
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	return buf, err
}

func AuthClient(config Config, conn net.Conn) error {
	authBuf := make([]byte, 256)
	count, err := io.ReadFull(conn, authBuf)
	if err != nil {
		return fmt.Errorf("error reading auth token: %s", err)
	}
	if count < len(authBuf) {
		return fmt.Errorf("error truncated auth token")
	}

	serverKey, err := hex.DecodeString(config.Key)
	if err != nil {
		return fmt.Errorf("error decoding server key")
	}

	diff := bytes.Compare(serverKey, authBuf)
	if diff != 0 {
		return fmt.Errorf("incorrect client key")
	}

	return nil
}
