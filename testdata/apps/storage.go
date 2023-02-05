package main

import (
	"bytes"
	"disco/network"
	"disco/network/discovery"
	"disco/network/instance"
	"disco/network/relay"
	"disco/storage"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

/// diskey: Distributed in-memory key-value storage

func main() {
	store := storage.NewTimedQueue(6 * time.Second)
	cfg := network.NewConfig("storage", ":6660")

	peers := discovery.Run("storage-one")

	relay := relay.NewRelay(peers, cfg)
	relay.Run()

	instance := instance.NewInstance(store, cfg)
	instance.Run()

	t := time.Tick(time.Second * 5)
	count := 0
	for {
		select {
		case <-t:
			r, err := http.Get(
				"http://localhost:6660/storage?key=ZZZ")
			if err != nil {
				fmt.Println(err)
				panic("wat")
			}
			resp, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()

			fmt.Printf("[%v] GET: %s (peers: %d)\n",
				r.StatusCode, resp, len(peers.Get()))

			if r.StatusCode != http.StatusOK {
				count += 1
				if count < 3 {
					continue
				}
				count = 0
				r, err = http.Post(
					"http://localhost:6660/storage?key=ZZZ",
					"application/json",
					bytes.NewBuffer([]byte(`{"Payload": "Yo"}`)))
				if err != nil {
					fmt.Println(err)
					panic("wat")
				}
				resp, _ := ioutil.ReadAll(r.Body)
				r.Body.Close()

				fmt.Printf("[%v] POST: %s\n",
					r.StatusCode, resp)
			}
		}
	}
}
