package models

import (
	"math"
	"testing"
)

func TestL2Norm(t *testing.T) {
	v1 := []float64{1, 0, 0}
	n := L2Norm(v1)
	if n != 1 {
		t.Error("The norm should be 1")
	}
	v2 := []float64{1, 1}
	n2 := L2Norm(v2)
	if n2 != math.Sqrt(2) {
		t.Error("The norm should be sqrt(2)")
	}
}
func TestL1Norm(t *testing.T) {
	v1 := []float64{1, 0, 0}
	u1 := []float64{2, 0, 0}
	n := L1Norm(v1, u1)
	if n != 1 {
		t.Error("The norm should be 1")
	}
	v2 := []float64{1, 1}
	u2 := []float64{1, 3}
	n2 := L1Norm(v2, u2)
	if n2 != 2 {
		t.Error("The norm should be sqrt(2)")
	}
}
