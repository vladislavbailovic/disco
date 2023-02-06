package main

import (
	"disco/logging"
	"disco/network/discovery"
	"time"
)

func main() {
	logging.Initialize(logging.Config{
		Level: logging.LevelTrace,
	})
	log := logging.Get()
	// dktest-one is a service name in k8s
	var autodiscoverySeed string = "autodiscovery-one"
	peers := discovery.Run(autodiscoverySeed)
	for {
		time.Sleep(5 * time.Second)

		log.Debug("peers: %v", peers.Get())
		log.Debug("status: %v", peers.Status())
	}
}
