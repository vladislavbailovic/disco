package storage

import (
	"strings"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	q := NewQueue()
	k, _ := NewKey("test")
	raw := `{"Payload":"test value"}`

	if err := q.Put(k, raw); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if got, err := q.Fetch(k); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if !strings.Contains(got.Value(), "test value") {
		t.Errorf("unexpected value: %q", got.Value())
	}
}

func TestTimedQueue(t *testing.T) {
	q := NewTimedQueue(10 * time.Millisecond)
	k, _ := NewKey("test")
	raw := `{"Payload":"test value"}`

	if err := q.Put(k, raw); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if got, err := q.Fetch(k); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if !strings.Contains(got.Value(), "test value") {
		t.Errorf("unexpected value: %q", got)
	}

	<-time.After(20 * time.Millisecond)

	if _, err := q.Fetch(k); err == nil {
		t.Errorf("expected cleanup, didn't get it")
	}
}
