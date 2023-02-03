package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

func autodiscover() *Peers {
	var peers *Peers = NewPeers()
	myself := fmt.Sprintf("%s", GetOutboundIP())
	peers.confirm(myself)
	go hello(peers)

	http.HandleFunc("/", handleHello(peers))
	go http.ListenAndServe(":6660", nil)

	return peers
}

func hello(peers *Peers) {
	t := time.Tick(time.Second * 5)
	myself := fmt.Sprintf("%s", GetOutboundIP())
	client := http.Client{
		Timeout: time.Second,
	}

	for {
		select {
		case <-t:
			cons := []string{autodiscoverySeed}
			if peers.totalLenExcept(myself) > 0 {
				cons = peers.getAll()
			}
			for _, addr := range cons {
				if addr == myself {
					// fmt.Println("not going to be pinging myself", addr, myself)
					continue
				}
				// fmt.Println("pinging", addr)
				r, err := client.Get(fmt.Sprintf("http://%s:6660", addr))
				if err != nil {
					// fmt.Println("well, something didn't go well", err)
					peers.unconfirm(addr)
					continue
				}

				if r.StatusCode != http.StatusOK {
					// fmt.Println("NOT OK!", addr)
					peers.unconfirm(addr)
					continue
				}

				var res []string
				defer r.Body.Close()
				if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
					// fmt.Println("error parsing response", err)
					peers.unconfirm(addr)
					continue
				}
				// fmt.Println("adding cons from addr", addr, res)
				peers.add(res...)

				peers.setReady(len(res) == len(peers.getConfirmed()))
			}
		}
	}
}

func handleHello(peers *Peers) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		host, _ /*port*/, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Println("unable to split host/port", err)
			return
		}
		// fmt.Println("confirming", host)
		peers.confirm(host)
		json.NewEncoder(w).Encode(peers.Get())
	}
}
