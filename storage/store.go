package storage

import (
	"disco/logging"
	"encoding/json"
)

type Storer interface {
	Fetch(*Key) (Valuer, error)
	Put(*Key, string) error
	Delete(*Key) error
	Stats() *Stats
}

type Valuer interface {
	Value() string
	MIME() ContentType
}

type Meter struct {
	Label string
	Value int
}

func NewMeter(label string, value int) Meter {
	return Meter{
		Label: label,
		Value: value,
	}
}

type Stats []Meter

func NewStats(m ...Meter) *Stats {
	s := Stats(m)
	return &s
}

func DecodeStats(from []byte) (*Stats, error) {
	var stats Stats
	if err := json.Unmarshal(from, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

func (x *Stats) MIME() ContentType {
	return ContentTypeJSON
}

func (x *Stats) Value() string {
	dst, err := json.Marshal(x)
	if err != nil {
		logging.Get().Error("Marshaling stats to JSON: %v", err)
	}
	return string(dst)
}

func (x *Stats) Len() int {
	return len(*x)
}

func (x *Stats) GetMeters() []Meter {
	return *x
}

func (x *Stats) Add(meters ...Meter) {
	for _, m := range meters {
		*x = append(*x, m)
	}
}

func (x *Stats) Merge(s *Stats) {
	if s != nil {
		x.Add((*s)...)
	}
}

func (x *Stats) Sum() *Stats {
	values := map[string]int{}
	for _, m := range *x {
		values[m.Label] += m.Value
	}
	result := NewStats()
	for label, value := range values {
		result.Add(NewMeter(label, value))
	}
	return result
}

func Default() Storer {
	return NewPlainStore()
}

type ContentType uint

const (
	ContentTypeText ContentType = iota
	ContentTypeJSON
)

func (x ContentType) String() string {
	switch x {
	case ContentTypeText:
		return "text/plain"
	case ContentTypeJSON:
		return "application/json"
	}
	logging.Get().Fatal("Unknown content type: %d", x)
	panic("Unknown content type")
}
