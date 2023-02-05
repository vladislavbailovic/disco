package storage

import (
	"fmt"
	"testing"
)

func TestDefault(t *testing.T) {
	s := Default()
	fmt.Println(s)
}
