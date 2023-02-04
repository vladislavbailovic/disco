package network

import (
	"sort"
	"sync"
)

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

func (x *Peers) Get() []string {
	cons := x.getConfirmed()
	sort.Strings(cons)
	return cons
}

func (x *Peers) Status() DiscoveryStatus {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return x.status
}

/// Public because of tests
func (x *Peers) SetReady(ready bool) {
	x.lock.Lock()
	defer x.lock.Unlock()
	if ready {
		x.status = Ready
	} else {
		x.status = EstablishingQuorum
	}
}

func (x *Peers) getAll() []string {
	x.lock.RLock()
	defer x.lock.RUnlock()
	cons := make([]string, 0, len(x.cons))
	for addr, _ := range x.cons {
		cons = append(cons, addr)
	}
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

func (x *Peers) totalLenExcept(addr string) int {
	x.lock.RLock()
	defer x.lock.RUnlock()
	count := len(x.cons)
	if _, ok := x.cons[addr]; ok {
		count -= 1
	}
	return count
}

func (x *Peers) add(cons ...string) {
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

/// Confirm just adds address unconditionally
/// Public because of testing
func (x *Peers) Confirm(cons ...string) {
	x.lock.Lock()
	defer x.lock.Unlock()
	for _, c := range cons {
		x.cons[c] = true
	}
}

func (x *Peers) unconfirm(cons ...string) {
	x.lock.Lock()
	defer x.lock.Unlock()
	for _, c := range cons {
		if _, ok := x.cons[c]; ok {
			// Only unconfirm previously known addresses
			x.cons[c] = false
		}
	}
}
