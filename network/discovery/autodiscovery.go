package discovery

import (
	"disco/network"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

func Run(seed string) *network.Peers {
	var peers *network.Peers = network.NewPeers()
	myself := fmt.Sprintf("%s", network.GetOutboundIP())
	peers.Confirm(myself)
	go hello(seed, peers)

	http.HandleFunc("/hello", handleHello(peers))
	go http.ListenAndServe(":6660", nil)

	return peers
}

func hello(seed string, peers *network.Peers) {
	t := time.Tick(time.Second * 5)
	myself := fmt.Sprintf("%s", network.GetOutboundIP())
	client := http.Client{
		Timeout: time.Second,
	}

	for {
		select {
		case <-t:
			cons := []string{seed}
			if peers.TotalLenExcept(myself) > 0 {
				cons = peers.GetAll()
			}
			for _, addr := range cons {
				if addr == myself {
					// fmt.Println("not going to be pinging myself", addr, myself)
					continue
				}
				// fmt.Println("pinging", addr)
				r, err := client.Get(fmt.Sprintf("http://%s:6660/hello", addr))
				if err != nil {
					// fmt.Println("well, something didn't go well", err)
					peers.Unconfirm(addr)
					continue
				}

				if r.StatusCode != http.StatusOK {
					// fmt.Println("NOT OK!", addr)
					peers.Unconfirm(addr)
					continue
				}

				var res []string
				defer r.Body.Close()
				if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
					// fmt.Println("error parsing response", err)
					peers.Unconfirm(addr)
					continue
				}
				// fmt.Println("adding cons from addr", addr, res)
				peers.Add(res...)

				peers.SetReady(len(res) == len(peers.GetConfirmed()))
			}
		}
	}
}

func handleHello(peers *network.Peers) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		host, _ /*port*/, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Println("unable to split host/port", err)
			return
		}
		// fmt.Println("confirming", host)
		peers.Confirm(host)
		json.NewEncoder(w).Encode(peers.Get())
	}
}
