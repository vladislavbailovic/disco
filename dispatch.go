package main

import (
	"disco/network"
	"fmt"
	"io"
	"net/http"
)

type Dispatch struct {
	peers *network.Peers
}

func NewDispatch(peers *network.Peers) *Dispatch {
	return &Dispatch{
		peers: peers,
	}
}

func (x *Dispatch) handle(w http.ResponseWriter, r *http.Request) {
	if x.peers == nil || x.peers.Status() != network.Ready {
		w.WriteHeader(http.StatusInternalServerError)
		if x.peers != nil {
			w.Write([]byte(x.peers.Status().String()))
		}
		return
	}

	key := r.URL.Query()["key"][0]
	fmt.Println("Requested", key)
	instance := getInstance(key, x.peers)
	fmt.Println(network.GetOutboundIP(), "peers:", x.peers.Status(), x.peers.Get())
	requestUrl := "http://" + instance + ":6660/_storage"
	fmt.Println("Gonna ask instance", requestUrl+"?key="+key)

	resp, err := http.Get(requestUrl + "?key=" + key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println("instance responded with", resp.StatusCode)
	w.WriteHeader(resp.StatusCode)

	defer r.Body.Close()
	io.Copy(w, r.Body)
}
