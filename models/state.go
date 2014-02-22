package models

import (
  "github.com/ryadzenine/dolphin/mpi"
)

type State interface {
  mpi.Versionable
  Values() []float64
}
type States []State

func (states States) ComputeAgregation() []float64 {
  N := len(states)
  agg := make([]float64, len(states[0].Values()))
  for _, s := range states {
    for j, v := range s.Values() {
      agg[j] = agg[j] + v
    }
  }
  for i, v := range agg {
    agg[i] = v / float64(N)
  }
  return agg
}
