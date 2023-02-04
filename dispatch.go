package main

import (
	"disco/network"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Dispatch struct {
	peers           *network.Peers
	client          *http.Client
	storageEndpoint string
	storagePort     string
}

func NewDispatch(peers *network.Peers) *Dispatch {
	return &Dispatch{
		peers:           peers,
		storageEndpoint: "_storage",
		storagePort:     "6660",
		client: &http.Client{
			Timeout: 1 * time.Second,
		},
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

	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqUrl := x.getInstanceURL(key)
	fmt.Printf("[%v]: dispatching to instance %q\n",
		network.GetOutboundIP(),
		reqUrl.String())

	resp, err := x.client.Get(reqUrl.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(resp.StatusCode)

	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func (x *Dispatch) getInstanceURL(key string) url.URL {
	instance := x.getInstance(key)
	return url.URL{
		Scheme:   "http",
		Host:     instance + ":" + x.storagePort,
		Path:     x.storageEndpoint,
		RawQuery: "key=" + url.QueryEscape(key),
	}
}

func (x *Dispatch) getInstance(key string) string {
	keyspaces := []struct {
		min, max int
	}{
		{min: int('0'), max: int('9')},
		{min: int('A'), max: int('Z')},
		{min: int('a'), max: int('z')},
	}
	instances := x.peers.Get()
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
