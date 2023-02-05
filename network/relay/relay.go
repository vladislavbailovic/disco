package relay

import (
	"bytes"
	"disco/network"
	"disco/store"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Relay struct {
	peers           *network.Peers
	client          *http.Client
	storageEndpoint string
	storagePort     string
}

func NewRelay(peers *network.Peers, cfg network.Config) *Relay {
	return &Relay{
		peers:           peers,
		storageEndpoint: cfg.InstancePath,
		storagePort:     cfg.Port,
		client: &http.Client{
			Timeout: 1 * time.Second,
		},
	}
}

func (x *Relay) Handle(w http.ResponseWriter, r *http.Request) {
	if x.peers == nil || x.peers.Status() != network.Ready {
		w.WriteHeader(http.StatusInternalServerError)
		if x.peers != nil {
			w.Write([]byte(x.peers.Status().String()))
		}
		return
	}

	key, err := store.NewKey(r.URL.Query().Get("key"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}

	reqUrl := x.getInstanceURL(key)
	fmt.Printf("[%v]: dispatching to instance %q\n",
		network.GetOutboundIP(),
		reqUrl.String())

	switch r.Method {
	case http.MethodGet:
		x.handleGet(reqUrl, w, r)
		return
	case http.MethodPost:
		x.handlePost(reqUrl, w, r)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func (x *Relay) handleGet(reqUrl url.URL, w http.ResponseWriter, r *http.Request) {
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

func (x *Relay) handlePost(reqUrl url.URL, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := x.client.Post(
		reqUrl.String(), "text/plain", bytes.NewBuffer(value))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(resp.StatusCode)

	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func (x *Relay) getInstanceURL(key *store.Key) url.URL {
	instance := x.getInstance(key)
	return url.URL{
		Scheme:   "http",
		Host:     instance + ":" + x.storagePort,
		Path:     x.storageEndpoint,
		RawQuery: "key=" + url.QueryEscape(key.String()),
	}
}

func (x *Relay) getInstance(key *store.Key) string {
	instances := x.peers.Get()
	for _, keyspace := range store.Keyspaces {
		if keyspace.InKeyspace(key) {
			stride := keyspace.GetRange() / len(instances)
			idx := keyspace.GetPosition(key) / stride
			if idx > 0 {
				idx -= 1
			}
			return instances[idx]
		}
	}
	panic("Unable to map instance to key")
}
