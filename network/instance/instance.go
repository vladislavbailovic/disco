package instance

import (
	"disco/network"
	"disco/store"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Instance struct {
	store.Storer
}

func NewInstance(str store.Storer) *Instance {
	if str == nil {
		str = store.Default()
	}
	return &Instance{
		Storer: str,
	}
}

func (x *Instance) Handle(w http.ResponseWriter, r *http.Request) {
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

func (x *Instance) handleGet(key *store.Key, w http.ResponseWriter, r *http.Request) {
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

func (x *Instance) handlePost(key *store.Key, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%v]: sets key %q in storage\n",
		network.GetOutboundIP(), key)

	defer r.Body.Close()
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	x.Put(key, string(value))
	w.WriteHeader(http.StatusOK)
}
