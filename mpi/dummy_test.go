package mpi

import "testing"

type Mock int

func (i Mock) Version() int {
	return int(i)
}
func TestDummys(t *testing.T) {
	mpi := Dummy(make(map[string]Versioner))
	mpi.Write("q1", Mock(5))
	if name, ok := mpi["q1"]; ok {
		if name.(Mock) != 5 {
			t.Error("Expected an element in the DummyQueue")
		}
	}
	if mpi.ReadFirst("q1").(Mock) != 5 {
		t.Error("Expected Value 5")
	}
	mpi.Write("q2", Mock(6))
	data := mpi.ReadFirstAll()
	if v, ok := data["q1"]; !ok || v.(Mock) != 5 {
		t.Error("Expected queue Q1 to be filled with 5")
	}
	if v, ok := data["q2"]; !ok || v.(Mock) != 6 {
		t.Error("Expected queue Q2 to be filled with 6")
	}
	dd := mpi.ReadStates(map[string]int{"q1": 5, "q2": 4})
	if v, ok := dd["q2"]; !ok || v.(Mock) != 6 {
		t.Error("Expected Q2 to be returned in ReadStates")
	}
	if _, ok := dd["q1"]; ok {
		t.Error("not expecting QA to be returned in ReadStates")
	}
}
