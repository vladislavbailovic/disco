package relay

import (
	"disco/network"
	"net/http"
	"time"
)

type Metrics struct {
	peers  *network.Peers
	apiKey *network.ApiKey
	client *http.Client
	cfg    network.Config
}

func NewMetrics(peers *network.Peers, cfg network.Config) *Metrics {
	return &Metrics{
		peers:  peers,
		apiKey: network.NewApiKey(cfg.KeyBase),
		client: &http.Client{
			Timeout: 1 * time.Second,
		},
		cfg: cfg,
	}
}

func (x *Metrics) handle(w http.ResponseWriter, r *http.Request) {
	// TODO: serve metrics
}
