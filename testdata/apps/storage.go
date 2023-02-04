package main

import (
	"disco/network"
	"fmt"
	"net/http"
	"time"
)

/// diskey: Distributed in-memory key-value storage

func main() {
	peers := network.Autodiscover("storage-one")
	dispatch := NewDispatch(peers)
	http.HandleFunc("/storage", dispatch.handle)
	http.HandleFunc("/_storage", handleStorage)
	go http.ListenAndServe(":6660", nil)

	t := time.Tick(time.Second * 5)
	for {
		select {
		case <-t:
			r, err := http.Get("http://localhost:6660/storage?key=AAA")
			if err != nil {
				fmt.Println(err)
				panic("wat")
			}
			fmt.Println("client status:", r.StatusCode)
		}
	}
}

func handleStorage(http.ResponseWriter, *http.Request) {
	fmt.Println(network.GetOutboundIP(), "gets key from storage")
}

func getInstance(key string, peers *network.Peers) string {
	keyspaces := []struct {
		min, max int
	}{
		{min: int('0'), max: int('9')},
		{min: int('A'), max: int('Z')},
		{min: int('a'), max: int('z')},
	}
	instances := peers.Get()
	for _, keyspace := range keyspaces {
		test := int(key[0])
		if test >= keyspace.min && test <= keyspace.max {
			total := keyspace.max - keyspace.min + 1
			test -= keyspace.min + 1
			stride := total / len(instances)
			idx := (test - 1) / stride
			return instances[idx]
		}
	}
	panic("GTFO")
}

var _storage map[string]string = make(map[string]string)

func getStored(key string) (string, error) {
	return _storage[key], nil
}

func setStored(key, value string) error {
	_storage[key] = "STORED " + value
	return nil
}
