package models

type RateFunc func(int) float64
type StopFunc func([]float64, []float64) bool

type Optimizable interface {
	PartialDerivative(int) []float64 // Should return the i-th partial derivative of the Loss function
	Model() []float64                // Returns the actual parameters
	Update([]float64)                // Updates the parameters
	Loss(SLPoint) float64
	Dataset() SLDataset // Returns the optimizable dataset
}

type DistributedSGD struct {
	Instance  Optimizable
	Name      string
	Id        int
	Rate      RateFunc
	Stop      StopFunc
	iteration int
}

func (s DistributedSGD) State() State {
	return DistributedSGDState{State: s.Instance.Model(), version: s.iteration}
}
func (s DistributedSGD) Error(errSet []SLPoint) float64 {
	acc := 0.
	for _, v := range errSet {
		acc += s.Instance.Loss(v)
	}
	return acc / float64(len(s.Instance.Dataset()))
}

func (s *DistributedSGD) Compute(iteration int) {
	index := iteration % len(s.Instance.Dataset())
	grad := s.Instance.PartialDerivative(index)
	rt := s.Rate(iteration)
	model := s.Instance.Model()
	for k, v := range s.Instance.Model() {
		model[k] = v - rt*grad[k]
	}
	s.iteration += 1
}
func (s *DistributedSGD) Average(states States, iteration int) {
	index := iteration % len(s.Instance.Dataset())
	states = append(states, s.State())
	convexPart := states.Average()
	grad := s.Instance.PartialDerivative(index)
	rt := s.Rate(iteration)
	model := s.Instance.Model()
	for k := range s.Instance.Model() {
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
