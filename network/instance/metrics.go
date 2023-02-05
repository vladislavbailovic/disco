package instance

import (
	"disco/network"
	"disco/storage"
	"fmt"
	"net/http"
)

type Metrics struct {
	storage.Storer
	apiKey *network.ApiKey
	cfg    network.Config
}

func NewMetrics(str storage.Storer, cfg network.Config) *Metrics {
	return &Metrics{
		Storer: str,
		apiKey: network.NewApiKey(cfg.KeyBase),
		cfg:    cfg,
	}
}

func (x *Metrics) handle(w http.ResponseWriter, r *http.Request) {
	// Validate x-relay-key
	relayKey := network.NewApiKey(r.Header.Get("x-relay-key"))
	if !x.apiKey.Equals(relayKey) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "invalid relayKey: %q", relayKey)
		return
	}

	value := x.Stats()
	w.Header().Add("content-type", value.MIME().String())
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", value.Value())
}
