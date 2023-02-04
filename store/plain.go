package store

import "fmt"

type PlainStore map[string]string

func NewPlainStore() *PlainStore {
	return &PlainStore{}
}

func (x *PlainStore) Fetch(key *Key) (string, error) {
	if val, ok := (*x)[key.String()]; ok {
		return val, nil
	}
	return "", fmt.Errorf("unknown key: %q", key)
}

func (x *PlainStore) Put(key *Key, val string) error {
	(*x)[key.String()] = val
	return nil
}
