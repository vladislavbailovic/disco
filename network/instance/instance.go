package instance

import (
	"disco/network"
	"disco/storage"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Instance struct {
	storage.Storer
}

func NewInstance(str storage.Storer) *Instance {
	if str == nil {
		str = storage.Default()
	}
	return &Instance{
		Storer: str,
	}
}

func (x *Instance) Handle(w http.ResponseWriter, r *http.Request) {
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

func (x *Instance) handleGet(key *storage.Key, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%v]: gets key %q from storage\n",
		network.GetOutboundIP(), key)

	value, err := x.Fetch(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", value.Value())
}

func (x *Instance) handleDelete(key *storage.Key, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%v]: deletes key %q from storage\n",
		network.GetOutboundIP(), key)

	if err := x.Delete(key); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (x *Instance) handlePost(key *storage.Key, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%v]: sets key %q in storage\n",
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
