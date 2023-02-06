package discovery

import (
	"disco/logging"
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
	log := logging.Get()
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
					log.Trace("[%v] not going to be pinging myself (%v)", addr, myself)
					continue
				}

				log.Trace("pinging %v", addr)
				r, err := client.Get(fmt.Sprintf("http://%s:6660/hello", addr))
				if err != nil {
					log.Trace("error pinging %v: %v", addr, err)
					peers.Unconfirm(addr)
					continue
				}

				if r.StatusCode != http.StatusOK {
					log.Trace("%v responded with non-200 code: %d", addr, r.StatusCode)
					peers.Unconfirm(addr)
					continue
				}

				var res []string
				defer r.Body.Close()
				if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
					log.Trace("error parsing response from %v: %v", addr, err)
					peers.Unconfirm(addr)
					continue
				}
				log.Trace("adding connections %v from addr %v", res, addr)
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
			logging.Get().Error("discovery: unable to split host/port: %v", err)
			return
		}
		logging.Get().Trace("confirming %v", host)
		peers.Confirm(host)
		json.NewEncoder(w).Encode(peers.Get())
	}
}
