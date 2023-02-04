package storage

import (
	"disco/network"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type StoreConfig struct {
	Addr         string
	Port         string
	StoragePath  string
	DispatchPath string
}

func NewStoreConfig(base, addr string) StoreConfig {
	split := strings.SplitN(addr, ":", 2)
	return StoreConfig{
		Addr:         addr,
		Port:         split[1],
		DispatchPath: "/" + base,
		StoragePath:  "/_" + base,
	}
}

type Store struct {
	lock  sync.RWMutex
	store map[string]string
}

func NewStore() *Store {
	return &Store{
		store: map[string]string{},
	}
}

func (x *Store) Handle(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		x.handleGet(key, w, r)
		return
	case http.MethodPost:
		x.handlePost(key, w, r)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func (x *Store) handleGet(key string, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%v]: gets key %q from storage\n",
		network.GetOutboundIP(), key)

	value, err := x.fetch(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", value)
}

func (x *Store) handlePost(key string, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%v]: sets key %q in storage\n",
		network.GetOutboundIP(), key)

	defer r.Body.Close()
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	x.put(key, string(value))
	w.WriteHeader(http.StatusOK)
}

func (x *Store) fetch(key string) (string, error) {
	x.lock.RLock()
	defer x.lock.RUnlock()

	if val, ok := x.store[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("unknown key: %q", key)
}

func (x *Store) put(key, val string) {
	x.lock.Lock()
	defer x.lock.Unlock()

	x.store[key] = val
}
