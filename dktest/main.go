package main

import (
	"fmt"
	"time"
)

func main() {
	// dktest-one is a service name in k8s
	var autodiscoverySeed string = "dktest-one"
	peers := autodiscover(autodiscoverySeed)
	for {
		time.Sleep(5 * time.Second)
		fmt.Println("\n--- peers ---\n\t", peers.Get(), "\n\t", peers.Status(), "\n------")
	}
}
