package main

import (
	"disco/network"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type StoreConfig struct {
	addr         string
	port         string
	storagePath  string
	dispatchPath string
}

func NewStoreConfig(base, addr string) StoreConfig {
	split := strings.SplitN(addr, ":", 2)
	return StoreConfig{
		addr:         addr,
		port:         split[1],
		dispatchPath: "/" + base,
		storagePath:  "/_" + base,
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

func (x *Store) handle(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
