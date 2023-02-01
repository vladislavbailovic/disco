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

type Connections struct {
	lock sync.RWMutex
	cons map[string]bool
}

func NewConnections() *Connections {
	return &Connections{
		cons: make(map[string]bool, 10),
	}
}

func (x *Connections) GetAll() []string {
	cons := x.getAll()
	sort.Strings(cons)
	return cons
}

func (x *Connections) getAll() []string {
	x.lock.RLock()
	defer x.lock.RUnlock()
	cons := make([]string, 0, len(x.cons))
	for addr, _ := range x.cons {
		cons = append(cons, addr)
	}
	return cons
}

func (x *Connections) GetConfirmed() []string {
	cons := x.getConfirmed()
	sort.Strings(cons)
	return cons
}

func (x *Connections) getConfirmed() []string {
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

func (x *Connections) TotalLenExcept(addr string) int {
	x.lock.RLock()
	defer x.lock.RUnlock()
	count := len(x.cons)
	if _, ok := x.cons[addr]; ok {
		count -= 1
	}
	return count
}

func (x *Connections) Add(cons ...string) {
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

func (x *Connections) Confirm(cons ...string) {
	x.lock.Lock()
	defer x.lock.Unlock()
	for _, c := range cons {
		// Confirm just adds address unconditionally
		x.cons[c] = true
	}
}

var connections *Connections = NewConnections()

func hello() {
	t := time.Tick(time.Second * 5)
	myself := fmt.Sprintf("%s", GetOutboundIP())

	for {
		select {
		case <-t:
			cons := []string{autodiscoverySeed}
			if connections.TotalLenExcept(myself) > 0 {
				cons = connections.GetAll()
			}
			for _, addr := range cons {
				if addr == myself {
					// fmt.Println("not going to be pinging myself", addr, myself)
					continue
				}
				// fmt.Println("pinging", addr)
				r, err := http.Get(fmt.Sprintf("http://%s:6660", addr))
				if err != nil {
					fmt.Println("well, something didn't go well", err)
					continue
				}

				var res []string
				defer r.Body.Close()
				if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
					// fmt.Println("error parsing response", err)
					continue
				}
				fmt.Println("adding cons", res)
				connections.Add(res...)
			}
			fmt.Println(myself, "connections", connections.GetConfirmed())
		}
	}
}

func main() {
	myself := fmt.Sprintf("%s", GetOutboundIP())
	connections.Confirm(myself)
	go hello()

	http.HandleFunc("/", handleHello)
	http.ListenAndServe(":6660", nil)
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	host, _ /*port*/, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println("unable to split host/port", err)
		return
	}
	fmt.Println("confirming", host)
	connections.Confirm(host)
	json.NewEncoder(w).Encode(connections.GetAll())
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
