package storage

import (
	"testing"
)

func TestStats_AddMergeSum_HappyPath(t *testing.T) {
	s := NewStats()
	if s.Len() != 0 {
		t.Errorf("should be empty at first")
	}

	s.Add(NewMeter("test", 1))
	if s.Len() != 1 {
		t.Errorf("should have 1 meter")
	}

	s.Add(NewMeter("test", 6), NewMeter("test", 1))
	if s.Len() != 3 {
		t.Errorf("should have 3 meters")
	}

	s1 := NewStats(
		NewMeter("test", 1),
		NewMeter("test", 3),
		NewMeter("test", 1),
		NewMeter("test", 2),
	)
	s.Merge(s1)
	if s.Len() != 7 {
		t.Errorf("expected 7 meters")
	}

	res := s.Sum()
	if res.Len() != 1 {
		t.Errorf("sum should converge to 1 meter: %v", res)
	}
}

func TestStats_MergeSum_EdgeCases(t *testing.T) {
	s := NewStats()
	s.Merge(nil) // if doesn't panic, we're good

	s1 := s.Sum()
	if s1.Len() != 0 {
		t.Errorf("expected empty stats")
	}
}
