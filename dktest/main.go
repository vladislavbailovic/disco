package main

import (
	"fmt"
	"time"
)

// dktest-one is a service name in k8s
const autodiscoverySeed string = "dktest-one"

func main() {
	peers := autodiscover()
	for {
		time.Sleep(5 * time.Second)
		fmt.Println("\n--- peers ---\n\t", peers.Get(), "\n\t", peers.Status(), "\n------")
	}
}
