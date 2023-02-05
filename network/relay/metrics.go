package relay

import (
	"disco/network"
	"disco/storage"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sync"
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
	if x.peers == nil || x.peers.Status() != network.Ready {
		w.WriteHeader(http.StatusInternalServerError)
		if x.peers != nil {
			w.Write([]byte(x.peers.Status().String()))
		}
		return
	}

	// TODO: switch output type
	outputType := "prometheus"

	stats := x.getAllStats().Sum()
	if outputType == "json" {
		w.Header().Add("content-type", stats.MIME().String())
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, stats.Value())
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "HELP disco Disco storage metrics\n")
		fmt.Fprintf(w, "Type disco gauge\n")
		for _, m := range stats.GetMeters() {
			fmt.Fprintf(w, "disco(name=%q) %d\n", m.Label, m.Value)
		}
	}
}

func (x *Metrics) getAllStats() *storage.Stats {
	stats := storage.NewStats()
	var wg sync.WaitGroup

	for _, peer := range x.peers.Get() {
		wg.Add(1)
		go func(stats *storage.Stats, p string) {
			reqUrl := url.URL{
				Scheme: "http",
				Host:   p + ":" + x.cfg.Port,
				Path:   path.Join(x.cfg.InstancePath, "metrics"),
			}
			req, err := http.NewRequest(
				http.MethodGet, reqUrl.String(), nil)
			if err != nil {
				fmt.Printf("Error building request to %q: %v\n", p, err)
				return
			}

			req.Header.Add("x-relay-key", x.apiKey.String())
			resp, err := x.client.Do(req)
			if err != nil {
				fmt.Printf("Error getting %q: %v\n", p, err)
				return
			}

			defer resp.Body.Close()
			value, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response from %q: %v\n", p, err)
				return
			}

			is, err := storage.DecodeStats(value)
			if err != nil {
				fmt.Printf("Error decoding stats from %q: %v", p, err)
				return
			}

			stats.Merge(is)
			wg.Done()
		}(stats, peer)
	}

	wg.Wait()
	return stats
}
