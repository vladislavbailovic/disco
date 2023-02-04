package main

import (
	"bytes"
	"disco/network"
	"disco/storage"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

/// diskey: Distributed in-memory key-value storage

func main() {
	peers := network.Autodiscover("storage-one")
	cfg := storage.NewStoreConfig("storage", ":6660")
	dispatch := storage.NewDispatch(peers, cfg)
	store := storage.NewStore()
	http.HandleFunc(cfg.DispatchPath, dispatch.Handle)
	http.HandleFunc(cfg.StoragePath, store.Handle)
	go http.ListenAndServe(cfg.Addr, nil)

	t := time.Tick(time.Second * 5)
	for {
		select {
		case <-t:
			r, err := http.Get("http://localhost:6660/storage?key=AAA")
			if err != nil {
				fmt.Println(err)
				panic("wat")
			}
			resp, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()

			fmt.Printf("[%v] GET: %s\n",
				r.StatusCode, resp)

			if r.StatusCode != http.StatusOK {
				r, err = http.Post("http://localhost:6660/storage?key=AAA", "text/plain", bytes.NewBuffer([]byte("Yo")))
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
