package store

import "fmt"

type Key string

func NewKey(src string) (*Key, error) {
	k := Key(src)
	if !k.isValid() {
		return nil, fmt.Errorf("invalid key source: %q", src)
	}
	return &k, nil
}

func (x *Key) String() string {
	return string(*x)
}

func (x Key) isValid() bool {
	if len(x) < 1 {
		return false
	}
	if !isAlnum(x[0]) {
		return false
	}
	for i := 1; i < len(x); i++ {
		if !isAlnumDash(x[i]) {
			return false
		}
	}
	return true
}

func isAlnum(c byte) bool {
	if c >= KeyspaceDigit.Min && c <= KeyspaceDigit.Max {
		return true
	}
	if c >= KeyspaceLowercase.Min && c <= KeyspaceLowercase.Max {
		return true
	}
	if c >= KeyspaceUppercase.Min && c <= KeyspaceUppercase.Max {
		return true
	}
	return false
}

func isAlnumDash(c byte) bool {
	if c == uint8('-') {
		return true
	}
	return isAlnum(c)
}

type Keyspace struct {
	Min uint8
	Max uint8
}

func (x Keyspace) InKeyspace(key *Key) bool {
	if key == nil {
		return false
	}
	if (*key)[0] >= x.Min && (*key)[0] <= x.Max {
		return true
	}
	return false
}

func (x Keyspace) GetPosition(key *Key) int {
	first := (*key)[0]
	return int((first - x.Min + 1) - 1)
}

func (x Keyspace) GetRange() int {
	return int(x.Max - x.Min + 1)
}

var KeyspaceDigit Keyspace = Keyspace{
	Min: uint8('0'),
	Max: uint8('9'),
}
var KeyspaceLowercase Keyspace = Keyspace{
	Min: uint8('a'),
	Max: uint8('z'),
}
var KeyspaceUppercase Keyspace = Keyspace{
	Min: uint8('A'),
	Max: uint8('Z'),
}
var Keyspaces []Keyspace = []Keyspace{
	KeyspaceDigit,
	KeyspaceLowercase,
	KeyspaceUppercase,
}
