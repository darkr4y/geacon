package crypt

import (
	"geacon/cmd/config"
	"math/rand"
	"time"
)

func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func RandomAESKey() {
	config.GlobalKey = make([]byte,16)
	_, err := rand.Read(config.GlobalKey[:])
	if err != nil {
		panic(err)
	}
}