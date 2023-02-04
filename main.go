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

var _storage map[string]string = make(map[string]string)

func getStored(key string) (string, error) {
	return _storage[key], nil
}

func setStored(key, value string) error {
	_storage[key] = "STORED " + value
	return nil
}
