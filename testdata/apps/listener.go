package main

import (
	"disco/network"
	"disco/network/instance"
	"disco/network/relay"
	"disco/storage"
	"net/http"
	"time"
)

/// diskey: Distributed in-memory key-value storage

func main() {
	store := storage.NewTimedQueue(5 * 60 * time.Second)
	cfg := network.NewConfig("storage", ":6660")

	peers := network.Autodiscover("listener-one")
	relay := relay.NewRelay(peers, cfg)
	instance := instance.NewInstance(store, cfg)
	http.HandleFunc(cfg.RelayPath, relay.Handle)
	http.HandleFunc(cfg.InstancePath, instance.Handle)
	http.ListenAndServe(cfg.Host+":"+cfg.Port, nil)
}
