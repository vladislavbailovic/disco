package main

import (
	"disco/network"
	"disco/network/discovery"
	"disco/network/instance"
	"disco/network/relay"
	"disco/storage"
	"time"
)

/// diskey: Distributed in-memory key-value storage

func main() {
	store := storage.NewTimedQueue(5 * 60 * time.Second)
	cfg := network.NewConfig("storage", ":6660")

	peers := discovery.Run("listener-one")

	relay := relay.NewRelay(peers, cfg)
	relay.Run()

	instance := instance.NewInstance(store, cfg)
	instance.Run()

	for {
		<-time.After(time.Second)
	}
}
