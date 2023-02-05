package storage

import (
	"encoding/json"
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

func (x JobStatus) String() string {
	switch x {
	case JobQueued:
		return "Queued"
	case JobRunning:
		return "Running"
	case JobDone:
		return "Done"
	case JobExpired:
		return "Expired"
	}
	panic("Unknown job status")
}

type Job struct {
	Status  JobStatus
	Payload string
}

func NewJob(pld string) (Job, error) {
	var job Job
	if err := json.Unmarshal([]byte(pld), &job); err != nil {
		return Job{}, err
	}
	return job, nil
}

func (x Job) MIME() ContentType { return ContentTypeJSON }

func (x Job) Value() string {
	dst, err := json.Marshal(x)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
	}
	return string(dst)
}

func (x Job) isFinished() bool {
	return x.Status == JobDone || x.Status == JobExpired
}

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

	if ok && !job.isFinished() {
		return fmt.Errorf("job %q already queued", key)
	}
	job, err := NewJob(val)
	if err != nil {
		return fmt.Errorf(
			"error serializing job %q: %v", key, err)
	}
	return x.put(key, job)
}

func (x *Queue) Delete(key *Key) error {
	x.lock.Lock()
	defer x.lock.Unlock()
	delete(x.data, key.String())
	return nil
}

func (x *Queue) put(key *Key, job Job) error {
	x.lock.Lock()
	defer x.lock.Unlock()
	x.data[key.String()] = job
	return nil
}

func (x *Queue) Stats() *Stats {
	x.lock.RLock()
	defer x.lock.RUnlock()

	var q, r, d, e int
	for _, job := range x.data {
		switch job.Status {
		case JobQueued:
			q += 1
		case JobRunning:
			r += 1
		case JobDone:
			d += 1
		case JobExpired:
			e += 1
		}
	}
	return NewStats(
		NewMeter("Queued", q),
		NewMeter("Running", r),
		NewMeter("Done", d),
		NewMeter("Expired", e),
	)
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
	<-time.After(x.removal)
	x.Delete(key)
}
