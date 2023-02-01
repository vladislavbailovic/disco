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
	cons []string
}

func NewConnections() *Connections {
	return &Connections{}
}

func (x *Connections) GetAll() []string {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return x.cons
}

func (x *Connections) Len() int {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return len(x.cons)
}

func (x *Connections) LenExcept(addr string) int {
	x.lock.RLock()
	defer x.lock.RUnlock()
	count := 0
	for _, c := range x.cons {
		if c == addr {
			continue
		}
		count += 1
	}
	return count
}

func (x *Connections) Add(cons ...string) {
	for _, c := range cons {
		x.add(c)
	}
	x.lock.Lock()
	defer x.lock.Unlock()
	sort.Strings(x.cons)
}

func (x *Connections) add(c string) {
	x.lock.Lock()
	defer x.lock.Unlock()
	alreadyPresent := false
	for _, con := range x.cons {
		if con == c {
			alreadyPresent = true
			break
		}
	}
	if !alreadyPresent {
		x.cons = append(x.cons, c)
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
			if connections.LenExcept(myself) > 0 {
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
				// fmt.Println("request: adding cons", res)
				connections.Add(res...)
			}
			fmt.Println(myself, "connections", connections.GetAll())
		}
	}
}

func main() {
	myself := fmt.Sprintf("%s", GetOutboundIP())
	connections.Add(myself)
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
	fmt.Println("response: adding host con", host)
	connections.Add(host)
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
