package storage

import "testing"

func Test_KeyValid(t *testing.T) {
	suite := map[string]bool{
		"1312":       true,
		"onetwo":     true,
		"OneTwo":     true,
		"123test123": true,
		"12-test-12": true,
		"-test-12":   false,
		"%#$^$%&^$":  false,
	}
	for test, want := range suite {
		t.Run(test, func(t *testing.T) {
			got := true
			if _, err := NewKey(test); err != nil {
				got = false
			}
			if want != got {
				t.Errorf("%q: want %v, got %v",
					test, want, got)
			}
		})
	}
}

func Test_Keyspace_Contains(t *testing.T) {
	suite := map[byte][]bool{
		// digits
		'0': []bool{true, false, false},
		'1': []bool{true, false, false},
		'5': []bool{true, false, false},
		'9': []bool{true, false, false},
		// lc
		'a': []bool{false, true, false},
		'b': []bool{false, true, false},
		'r': []bool{false, true, false},
		'z': []bool{false, true, false},
		// uc
		'A': []bool{false, false, true},
		'B': []bool{false, false, true},
		'R': []bool{false, false, true},
		'Z': []bool{false, false, true},
		// garbage
		'+': []bool{false, false, false},
		'!': []bool{false, false, false},
		'.': []bool{false, false, false},
		'-': []bool{false, false, false},
	}
	for test, want := range suite {
		t.Run("digit", func(t *testing.T) {
			got := KeyspaceDigit.contains(test)
			if got != want[0] {
				t.Errorf("%c (%q): want %v, got %v",
					test, test, want[0], got)
			}
		})
		t.Run("lowercase", func(t *testing.T) {
			got := KeyspaceLowercase.contains(test)
			if got != want[1] {
				t.Errorf("%c (%q): want %v, got %v",
					test, test, want[1], got)
			}
		})
		t.Run("uppercase", func(t *testing.T) {
			got := KeyspaceUppercase.contains(test)
			if got != want[2] {
				t.Errorf("%c (%q): want %v, got %v",
					test, test, want[2], got)
			}
		})
	}
}
