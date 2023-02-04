package main

import (
	"bytes"
	"disco/network"
	"disco/network/instance"
	"disco/network/relay"
	"disco/store"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

/// diskey: Distributed in-memory key-value storage

func main() {
	store := store.NewPlainStore()
	cfg := network.NewConfig("storage", ":6660")

	peers := network.Autodiscover("storage-one")
	dispatchHandler := relay.NewDispatch(peers, cfg)
	storageHandler := instance.NewStorage(store)
	http.HandleFunc(cfg.RelayPath, dispatchHandler.Handle)
	http.HandleFunc(cfg.InstancePath, storageHandler.Handle)
	go http.ListenAndServe(cfg.Addr, nil)

	t := time.Tick(time.Second * 5)
	count := 0
	for {
		select {
		case <-t:
			r, err := http.Get("http://localhost:6660/storage?key=ZZZ")
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
				r, err = http.Post("http://localhost:6660/storage?key=ZZZ", "text/plain", bytes.NewBuffer([]byte("Yo")))
				if err != nil {
					fmt.Println(err)
					panic("wat")
				}

				fmt.Printf("[%v] POST: %s\n",
					r.StatusCode, resp)
			}
		}
	}
}
