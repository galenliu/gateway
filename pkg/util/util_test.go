package util

import (
	"testing"
)

func TestUtil(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Log(GenerateMac())
	}
}
