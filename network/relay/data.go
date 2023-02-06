package relay

import (
	"bytes"
	"disco/logging"
	"disco/network"
	"disco/storage"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Data struct {
	peers  *network.Peers
	apiKey *network.ApiKey
	client *http.Client
	cfg    network.Config
}

func NewData(peers *network.Peers, cfg network.Config) *Data {
	return &Data{
		peers:  peers,
		apiKey: network.NewApiKey(cfg.KeyBase),
		client: &http.Client{
			Timeout: 1 * time.Second,
		},
		cfg: cfg,
	}
}

func (x *Data) handle(w http.ResponseWriter, r *http.Request) {

	log := logging.Get()
	if x.peers == nil || x.peers.Status() != network.Ready {
		w.WriteHeader(http.StatusInternalServerError)
		if x.peers != nil {
			w.Write([]byte(x.peers.Status().String()))
		} else {
			log.Debug("no peers")
		}
		return
	}

	key, err := storage.NewKey(r.URL.Query().Get("key"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		log.Error("Unable to parse key: %v", err)
		return
	}

	reqUrl := x.getInstanceURL(key)
	// fmt.Printf("[%v]: dispatching to instance %q\n",
	// 	network.GetOutboundIP(),
	// 	reqUrl.String())
	log.Info("[%v]: dispatching to instance %q\n",
		network.GetOutboundIP(),
		reqUrl.String())

	switch r.Method {
	case http.MethodGet:
		x.handleGet(reqUrl, w, r)
		return
	case http.MethodPost:
		x.handlePost(reqUrl, w, r)
		return
	case http.MethodDelete:
		x.handleDelete(reqUrl, w, r)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func (x *Data) handleDelete(reqUrl url.URL, w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(
		http.MethodDelete, reqUrl.String(), nil)
	req.Header.Add("x-relay-key", x.apiKey.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	resp, err := x.client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type",
		resp.Header.Get("content-type"))
	w.WriteHeader(resp.StatusCode)

	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func (x *Data) handleGet(reqUrl url.URL, w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(
		http.MethodGet, reqUrl.String(), nil)
	req.Header.Add("x-relay-key", x.apiKey.String())
	resp, err := x.client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type",
		resp.Header.Get("content-type"))
	w.WriteHeader(resp.StatusCode)

	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func (x *Data) handlePost(reqUrl url.URL, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		reqUrl.String(),
		bytes.NewBuffer(value))
	req.Header.Add("x-relay-key", x.apiKey.String())
	resp, err := x.client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type",
		resp.Header.Get("content-type"))
	w.WriteHeader(resp.StatusCode)

	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func (x *Data) getInstanceURL(key *storage.Key) url.URL {
	instance := x.getInstance(key)
	return url.URL{
		Scheme:   "http",
		Host:     instance + ":" + x.cfg.Port,
		Path:     x.cfg.InstancePath,
		RawQuery: "key=" + url.QueryEscape(key.String()),
	}
}

func (x *Data) getInstance(key *storage.Key) string {
	instances := x.peers.Get()
	for _, keyspace := range storage.Keyspaces {
		if keyspace.InKeyspace(key) {
			stride := keyspace.GetRange() / len(instances)
			idx := keyspace.GetPosition(key) / stride
			if idx > 0 {
				idx -= 1
			}
			return instances[idx]
		}
	}
	log := logging.Get()
	log.Fatal("Unable to map instance to key %v", key)
	log.Debug("Instances: %v", instances)
	panic("Unable to map instance to key")
}
