package main

import (
	"disco/network"
	"fmt"
	"time"
)

func main() {
	// dktest-one is a service name in k8s
	var autodiscoverySeed string = "autodiscovery-one"
	peers := network.Autodiscover(autodiscoverySeed)
	for {
		time.Sleep(5 * time.Second)
		fmt.Println("\n--- peers ---\n\t", peers.Get(), "\n\t", peers.Status(), "\n------")
	}
}
