package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
)

func authToServer(config Config, conn net.Conn) error {
	key, err := hex.DecodeString(config.ParentNode.Key)
	if err != nil {
		return fmt.Errorf("failed to decode auth key: %s", err)
	}

	buf := bytes.NewBuffer(key)
	_, err = buf.WriteTo(conn)
	if err != nil {
		return fmt.Errorf("failed to write key to server conn: %s", err)
	}

	return nil
}
