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
	if KeyspaceDigit.contains(c) {
		return true
	}
	if KeyspaceLowercase.contains(c) {
		return true
	}
	if KeyspaceUppercase.contains(c) {
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
	min uint8
	max uint8
}

func (x Keyspace) InKeyspace(key *Key) bool {
	if key == nil {
		return false
	}
	return x.contains((*key)[0])
}

func (x Keyspace) contains(c byte) bool {
	if c >= x.min && c <= x.max {
		return true
	}
	return false
}

func (x Keyspace) GetPosition(key *Key) int {
	first := (*key)[0]
	return int((first - x.min + 1) - 1)
}

func (x Keyspace) GetRange() int {
	return int(x.max - x.min + 1)
}

var KeyspaceDigit Keyspace = Keyspace{
	min: uint8('0'),
	max: uint8('9'),
}
var KeyspaceLowercase Keyspace = Keyspace{
	min: uint8('a'),
	max: uint8('z'),
}
var KeyspaceUppercase Keyspace = Keyspace{
	min: uint8('A'),
	max: uint8('Z'),
}
var Keyspaces []Keyspace = []Keyspace{
	KeyspaceDigit,
	KeyspaceLowercase,
	KeyspaceUppercase,
}
