package id

import (
	"crypto/rand"
	"encoding/hex"
)

type Generator interface {
	NewID() string
}

type SecureGenerator struct{}

func (SecureGenerator) NewID() string {
	b := make([]byte, 16) // 128bit
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
