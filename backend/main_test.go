package main

import (
	"math"
	"testing"
)

func TestDummy(t *testing.T) {
	got := math.Abs(-1)
	if got != 1 {
		t.Error("Abs not working")
	}
}
