package relay

import (
	"disco/network"
	"net/http"
)

type Relay struct {
	cfg     network.Config
	data    *Data
	metrics *Metrics
}

func NewRelay(peers *network.Peers, cfg network.Config) *Relay {
	return &Relay{
		cfg:     cfg,
		data:    NewData(peers, cfg),
		metrics: NewMetrics(peers, cfg),
	}
}

func (x *Relay) Run() {
	http.HandleFunc(x.cfg.RelayPath, x.data.handle)
	go http.ListenAndServe(x.cfg.Host+":"+x.cfg.Port, nil)
}
