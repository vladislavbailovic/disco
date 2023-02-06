package main

import (
	"disco/logging"
	"disco/network/discovery"
	"time"
)

func main() {
	// dktest-one is a service name in k8s
	var autodiscoverySeed string = "autodiscovery-one"
	peers := discovery.Run(autodiscoverySeed)
	logging.Initialize(logging.Config{
		Level: logging.LevelAll,
	})
	log := logging.Get()
	for {
		time.Sleep(5 * time.Second)

		log.Debug("----")
		log.Debug("peers: %v", peers.Get())
		log.Debug("status: %v", peers.Status())
		log.Debug("----")
	}
}
