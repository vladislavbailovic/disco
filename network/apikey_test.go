package network

import "testing"

func TestApiKey(t *testing.T) {
	k := NewApiKey("wat")
	if k.String() != "wat" {
		t.Errorf("unexpected api key result: %q", k)
	}
}

func TestApiKey_Equals(t *testing.T) {
	k := NewApiKey("wat")
	suite := map[string]bool{
		"wat":  true,
		"nope": false,
		"111":  false,
	}
	for test, want := range suite {
		t.Run(test, func(t *testing.T) {
			o := NewApiKey(test)
			got := k.Equals(o)
			if want != got {
				t.Errorf("%q: want %v, got %v",
					o, want, got)
			}
		})
	}
}
