package store

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
