package instance

import (
	"disco/network"
	"disco/storage"
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
	// TODO: serve metrics
}
