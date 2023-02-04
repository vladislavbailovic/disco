package store

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	name string
}

func (x Job) Value() string { return x.name }

type Queue struct {
	lock sync.RWMutex
	data map[string]Job
}

func NewQueue() *Queue {
	return &Queue{
		data: map[string]Job{},
	}
}

func (x *Queue) Fetch(key *Key) (Valuer, error) {
	x.lock.RLock()
	defer x.lock.RUnlock()

	if val, ok := x.data[key.String()]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("unknown key: %q", key)
}

func (x *Queue) Put(key *Key, val string) error {
	x.lock.Lock()
	defer x.lock.Unlock()

	// TODO: string to job
	x.data[key.String()] = Job{name: val}
	return nil
}

type TimedQueue struct {
	Queue
	expiry time.Duration
}

func NewTimedQueue(d time.Duration) *TimedQueue {
	queue := NewQueue()
	return &TimedQueue{
		Queue:  *queue,
		expiry: d,
	}
}

func (x *TimedQueue) Put(key *Key, val string) error {
	go x.cleanup(key)
	return x.Queue.Put(key, val)
}

func (x *TimedQueue) cleanup(key *Key) {
	<-time.After(x.expiry)
	delete(x.Queue.data, key.String())
}
