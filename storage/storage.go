package storage

import (
	"disco/network"
	"disco/store"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type StorageConfig struct {
	Addr         string
	Port         string
	StoragePath  string
	DispatchPath string
}

func NewStorageConfig(base, addr string) StorageConfig {
	split := strings.SplitN(addr, ":", 2)
	return StorageConfig{
		Addr:         addr,
		Port:         split[1],
		DispatchPath: "/" + base,
		StoragePath:  "/_" + base,
	}
}

type Storage struct {
	lock  sync.RWMutex
	store store.Storer
}

func NewStorage(str store.Storer) *Storage {
	if str == nil {
		str = store.Default()
	}
	return &Storage{
		store: str,
	}
}

func (x *Storage) Handle(w http.ResponseWriter, r *http.Request) {
	key, err := store.NewKey(r.URL.Query().Get("key"))
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
	}

	w.WriteHeader(http.StatusBadRequest)
}

func (x *Storage) handleGet(key *store.Key, w http.ResponseWriter, r *http.Request) {
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

func (x *Storage) handlePost(key *store.Key, w http.ResponseWriter, r *http.Request) {
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

func (x *Storage) fetch(key *store.Key) (string, error) {
	x.lock.RLock()
	defer x.lock.RUnlock()

	if val, err := x.store.Fetch(key); err != nil {
		return "", err
	} else {
		return val, nil
	}
	return "", fmt.Errorf("unknown key: %q", key)
}

func (x *Storage) put(key *store.Key, val string) {
	x.lock.Lock()
	defer x.lock.Unlock()

	x.store.Put(key, val)
}
