package models

import "github.com/ryadzenine/dolphin/mpi"

type stateErrorer interface {
	State() State
	Error([]SLPoint) float64
}
type computeAverager interface {
	Compute(SLPoint)
	Average([]float64, SLPoint)
}
type Estimate interface {
	computeAverager
	stateErrorer
}
type RegressionEstimate interface {
	Estimate
	Predict(Point) (float64, error)
}
type ClassificationEstimate interface {
	Estimate
	Predict(Point) (int, error)
}

type State interface {
	mpi.Versionable
	Values() []float64
}

type States []State

func (states States) Average() []float64 {
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
