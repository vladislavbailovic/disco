package storage

import (
	"fmt"
	"sync"
)

type PlainValue string

func (x PlainValue) Value() string {
	return string(x)
}

func (x PlainValue) MIME() ContentType {
	return ContentTypeText
}

type PlainStore struct {
	lock sync.RWMutex
	data map[string]PlainValue
}

func NewPlainStore() *PlainStore {
	return &PlainStore{
		data: map[string]PlainValue{},
	}
}

func (x *PlainStore) Fetch(key *Key) (Valuer, error) {
	x.lock.RLock()
	defer x.lock.RUnlock()

	if val, ok := x.data[key.String()]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("unknown key: %q", key)
}

func (x *PlainStore) Put(key *Key, val string) error {
	x.lock.Lock()
	defer x.lock.Unlock()

	x.data[key.String()] = PlainValue(val)
	return nil
}

func (x *PlainStore) Delete(key *Key) error {
	x.lock.Lock()
	defer x.lock.Unlock()

	delete(x.data, key.String())
	return nil
}

func (x *PlainStore) Stats() *Stats {
	x.lock.RLock()
	defer x.lock.RUnlock()

	return NewStats(NewMeter("Total", len(x.data)))
}
