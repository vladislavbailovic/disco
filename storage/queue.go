package storage

import (
	"fmt"
	"sync"
	"time"
)

type JobStatus uint

const (
	JobQueued JobStatus = iota
	JobRunning
	JobDone
	JobExpired
)

type Job struct {
	Status  JobStatus
	Payload string
}

// TODO: create new job from string
func NewJob(pld string) Job {
	return Job{Payload: pld}
}

func (x Job) Value() string { return x.Payload }

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
	job, ok := x.data[key.String()]
	x.lock.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown job: %q", key)
	}

	if job.Status != JobQueued {
		return nil, fmt.Errorf(
			"job %q already dispatched: %v", key, job.Status)
	}

	job.Status = JobRunning
	if err := x.put(key, job); err != nil {
		return nil, fmt.Errorf(
			"unable to update job %q status: %w", key, err)
	}

	return job, nil
}

func (x *Queue) Put(key *Key, val string) error {
	x.lock.RLock()
	job, ok := x.data[key.String()]
	x.lock.RUnlock()

	if ok && job.Status == JobQueued {
		return fmt.Errorf("job %q already queued", key)
	}
	job = NewJob(val)
	return x.put(key, job)
}

func (x *Queue) put(key *Key, job Job) error {
	x.lock.Lock()
	defer x.lock.Unlock()
	x.data[key.String()] = job
	return nil
}

type TimedQueue struct {
	Queue
	expiry  time.Duration
	removal time.Duration
}

func NewTimedQueue(d time.Duration) *TimedQueue {
	queue := NewQueue()
	return &TimedQueue{
		Queue:   *queue,
		expiry:  d,
		removal: d * 2,
	}
}

func (x *TimedQueue) Put(key *Key, val string) error {
	go x.expire(key)
	go x.remove(key)
	return x.Queue.Put(key, val)
}

func (x *TimedQueue) expire(key *Key) {
	<-time.After(x.expiry)
	x.lock.RLock()
	job, ok := x.data[key.String()]
	x.lock.RUnlock()

	if ok {
		job.Status = JobExpired
		x.put(key, job)
	}
}

func (x *TimedQueue) remove(key *Key) {
	<-time.After(x.expiry)
	delete(x.Queue.data, key.String())
}
