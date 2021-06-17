package main

import "crypto/rand"

func generateNewKey() ([]byte, error) {
	b := make([]byte, 256)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}
