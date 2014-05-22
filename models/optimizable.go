package models

type RateFunc func(int) float64
type Optimizable interface {
	PartialDerivative(int) []float64 // Should return the i-th partial derivative of the Loss function
	Model() []float64                // Returns the actual parameters
	Update([]float64)                // Updates the parameters
	Loss(SLPoint) float64
	Dataset() Dataset // Returns the optimizable dataset
}

type DistributedSGD struct {
	instance  Optimizable
	rate      RateFunc
	iteration int
}

func (s DistributedSGD) State() State {
	return DistributedSGDState{State: s.instance.Model(), version: s.iteration}
}
func (s DistributedSGD) Error(errSet []SLPoint) float64 {
	acc := 0.
	for _, v := range errSet {
		acc += s.instance.Loss(v)
	}
	return acc / float64(len(s.instance.Dataset()))
}

func (s *DistributedSGD) Compute(iteration int) {
	index := iteration % len(s.instance.Dataset())
	grad := s.instance.PartialDerivative(index)
	rt := s.rate(iteration)
	model := s.instance.Model()
	for k, v := range s.instance.Model() {
		model[k] = v - rt*grad[k]
	}
	s.iteration += 1
}
func (s *DistributedSGD) Average(states States, iteration int) {
	index := iteration % len(s.instance.Dataset())
	states = append(states, s.State())
	convexPart := states.Average()
	grad := s.instance.PartialDerivative(index)
	rt := s.rate(iteration)
	model := s.instance.Model()
	for k := range s.instance.Model() {
		model[k] = convexPart[k] - rt*grad[k]
	}
	s.iteration += 1

}

type DistributedSGDState struct {
	State   []float64
	version int
}

func (r DistributedSGDState) Values() []float64 {
	return r.State
}

func (r DistributedSGDState) Version() int {
	return r.version
}
