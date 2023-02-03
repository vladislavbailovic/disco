package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sort"
	"sync"
	"time"
)

// dktest-one is a service name in k8s
const autodiscoverySeed string = "dktest-one"

func main() {
	peers := autodiscover()
	for {
		time.Sleep(5 * time.Second)
		fmt.Println("\n--- peers ---\n\t", peers.GetConfirmed(), "\n\t", peers.Status(), "\n------")
	}
}

func autodiscover() *Peers {
	var peers *Peers = NewPeers()
	myself := fmt.Sprintf("%s", GetOutboundIP())
	peers.Confirm(myself)
	go hello(peers)

	http.HandleFunc("/", handleHello(peers))
	go http.ListenAndServe(":6660", nil)

	return peers
}

type DiscoveryStatus uint8

const (
	Init DiscoveryStatus = iota
	EstablishingQuorum
	Ready
)

func (x DiscoveryStatus) String() string {
	switch x {
	case Init:
		return "Initialized Discovery"
	case EstablishingQuorum:
		return "Establishing Quorum"
	case Ready:
		return "Ready"
	}
	panic("unknown discovery status")
}

type Peers struct {
	status DiscoveryStatus
	lock   sync.RWMutex
	cons   map[string]bool
}

func NewPeers() *Peers {
	return &Peers{
		cons: make(map[string]bool, 10),
	}
}

func (x *Peers) SetReady(ready bool) {
	x.lock.Lock()
	defer x.lock.Unlock()
	if ready {
		x.status = Ready
	} else {
		x.status = EstablishingQuorum
	}
}

func (x *Peers) Status() DiscoveryStatus {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return x.status
}

func (x *Peers) GetAll() []string {
	cons := x.getAll()
	sort.Strings(cons)
	return cons
}

func (x *Peers) getAll() []string {
	x.lock.RLock()
	defer x.lock.RUnlock()
	cons := make([]string, 0, len(x.cons))
	for addr, _ := range x.cons {
		cons = append(cons, addr)
	}
	return cons
}

func (x *Peers) GetConfirmed() []string {
	cons := x.getConfirmed()
	sort.Strings(cons)
	return cons
}

func (x *Peers) getConfirmed() []string {
	x.lock.RLock()
	defer x.lock.RUnlock()
	cons := make([]string, 0, len(x.cons))
	for addr, confirmed := range x.cons {
		if confirmed {
			cons = append(cons, addr)
		}
	}
	return cons
}

func (x *Peers) TotalLenExcept(addr string) int {
	x.lock.RLock()
	defer x.lock.RUnlock()
	count := len(x.cons)
	if _, ok := x.cons[addr]; ok {
		count -= 1
	}
	return count
}

func (x *Peers) Add(cons ...string) {
	x.lock.Lock()
	defer x.lock.Unlock()
	for _, c := range cons {
		if _, ok := x.cons[c]; !ok {
			// Only add if we don't know about its status previously
			// This is so that we don't trump its status if it's already confirmed
			x.cons[c] = false
		}
	}
}

func (x *Peers) Confirm(cons ...string) {
	x.lock.Lock()
	defer x.lock.Unlock()
	for _, c := range cons {
		// Confirm just adds address unconditionally
		x.cons[c] = true
	}
}

func (x *Peers) Unconfirm(cons ...string) {
	x.lock.Lock()
	defer x.lock.Unlock()
	for _, c := range cons {
		if _, ok := x.cons[c]; ok {
			// Only unconfirm previously known addresses
			x.cons[c] = false
		}
	}
}

func hello(peers *Peers) {
	t := time.Tick(time.Second * 5)
	myself := fmt.Sprintf("%s", GetOutboundIP())
	client := http.Client{
		Timeout: time.Second,
	}

	for {
		select {
		case <-t:
			cons := []string{autodiscoverySeed}
			if peers.TotalLenExcept(myself) > 0 {
				cons = peers.GetAll()
			}
			for _, addr := range cons {
				if addr == myself {
					// fmt.Println("not going to be pinging myself", addr, myself)
					continue
				}
				// fmt.Println("pinging", addr)
				r, err := client.Get(fmt.Sprintf("http://%s:6660", addr))
				if err != nil {
					// fmt.Println("well, something didn't go well", err)
					peers.Unconfirm(addr)
					continue
				}

				if r.StatusCode != http.StatusOK {
					// fmt.Println("NOT OK!", addr)
					peers.Unconfirm(addr)
					continue
				}

				var res []string
				defer r.Body.Close()
				if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
					// fmt.Println("error parsing response", err)
					peers.Unconfirm(addr)
					continue
				}
				// fmt.Println("adding cons from addr", addr, res)
				peers.Add(res...)

				peers.SetReady(len(res) == len(peers.GetConfirmed()))
			}
		}
	}
}

func handleHello(peers *Peers) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		host, _ /*port*/, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Println("unable to split host/port", err)
			return
		}
		// fmt.Println("confirming", host)
		peers.Confirm(host)
		json.NewEncoder(w).Encode(peers.GetConfirmed())
	}
}

// Get preferred outbound ip of this machine
// https://stackoverflow.com/a/37382208/12221657
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
