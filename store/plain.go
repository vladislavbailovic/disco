package store

import (
	"fmt"
	"sync"
)

type PlainStore struct {
	lock sync.RWMutex
	data map[string]string
}

func NewPlainStore() *PlainStore {
	return &PlainStore{
		data: map[string]string{},
	}
}

func (x *PlainStore) Fetch(key *Key) (string, error) {
	x.lock.RLock()
	defer x.lock.RUnlock()

	if val, ok := x.data[key.String()]; ok {
		return val, nil
	}
	return "", fmt.Errorf("unknown key: %q", key)
}

func (x *PlainStore) Put(key *Key, val string) error {
	x.lock.Lock()
	defer x.lock.Unlock()

	x.data[key.String()] = val
	return nil
}
