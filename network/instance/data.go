package instance

import (
	"disco/logging"
	"disco/network"
	"disco/storage"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Data struct {
	storage.Storer
	apiKey *network.ApiKey
	cfg    network.Config
}

func NewData(str storage.Storer, cfg network.Config) *Data {
	if str == nil {
		str = storage.Default()
	}
	return &Data{
		Storer: str,
		apiKey: network.NewApiKey(cfg.KeyBase),
		cfg:    cfg,
	}
}

func (x *Data) handle(w http.ResponseWriter, r *http.Request) {
	// Validate x-relay-key
	relayKey := network.NewApiKey(r.Header.Get("x-relay-key"))
	if !x.apiKey.Equals(relayKey) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "invalid relayKey: %q", relayKey)
		return
	}

	key, err := storage.NewKey(r.URL.Query().Get("key"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		x.handleGet(key, w, r)
		return
	case http.MethodPost:
		x.handlePost(key, w, r)
		return
	case http.MethodDelete:
		x.handleDelete(key, w, r)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func (x *Data) handleGet(key *storage.Key, w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("[%v]: gets key %q from storage\n",
	// 	network.GetOutboundIP(), key)
	logging.Get().Info("[%v]: gets key %q from storage\n",
		network.GetOutboundIP(), key)

	value, err := x.Fetch(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", err)
		return
	}

	w.Header().Add("content-type", value.MIME().String())
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", value.Value())
}

func (x *Data) handleDelete(key *storage.Key, w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("[%v]: deletes key %q from storage\n",
	// 	network.GetOutboundIP(), key)
	logging.Get().Info("[%v]: deletes key %q from storage\n",
		network.GetOutboundIP(), key)

	if err := x.Delete(key); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (x *Data) handlePost(key *storage.Key, w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("[%v]: sets key %q in storage\n",
	// 	network.GetOutboundIP(), key)
	logging.Get().Info("[%v]: sets key %q in storage\n",
		network.GetOutboundIP(), key)

	defer r.Body.Close()
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := x.Put(key, string(value)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
