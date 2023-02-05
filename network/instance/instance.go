package instance

import (
	"disco/network"
	"disco/storage"
	"net/http"
	"path"
)

type Instance struct {
	cfg     network.Config
	data    *Data
	metrics *Metrics
}

func NewInstance(str storage.Storer, cfg network.Config) *Instance {
	return &Instance{
		cfg:     cfg,
		data:    NewData(str, cfg),
		metrics: NewMetrics(str, cfg),
	}
}

func (x *Instance) Run() {
	http.HandleFunc(x.cfg.InstancePath, x.data.handle)
	http.HandleFunc(path.Join(x.cfg.InstancePath, "metrics"), x.metrics.handle)
	go http.ListenAndServe(x.cfg.Host+":"+x.cfg.Port, nil)
}
