package models

import (
  "testing"
)

type StateMock int

func (StateMock) Version() int {
  return 1
}
func (i StateMock) Values() []float64 {
  return []float64{float64(i)}
}
func TestComputeaggregation(t *testing.T) {
  states := States([]State{
    StateMock(1),
    StateMock(2),
    StateMock(3)})
  v := states.ComputeAgregation()
  if len(v) != 1 || v[0] != 2.0 {
    t.Error("Aggregation fonction wrong")

  }
}
