package np

import (
	"errors"
	"math"

	"github.com/ryadzenine/dolphin/models"
)

type EstimatorState struct {
	State   []float64
	version int
}

func (r EstimatorState) Values() []float64 {
	return r.State
}
func (r EstimatorState) Version() int {
	return r.version
}

type RevezEstimator struct {
	Vectors   []models.Vector
	state     []float64
	Step      int
	Rate      func(int) float64
	Smoothing func(int) float64
	Kernel    func(models.Vector) float64
}

func (r *RevezEstimator) Error(testData []models.SLPoint) float64 {
	return r.FastL2Error(testData)
}

// Fast L2 Error: is only to use when we are sure that the points are aligned
func (r *RevezEstimator) FastL2Error(testData []models.SLPoint) float64 {
	err := 0.0
	for key, v := range testData {
		p := r.state[key]
		err = err + (p-v.Y)*(p-v.Y)
	}
	return math.Sqrt(err) / float64(len(testData))
}

func (r *RevezEstimator) L2Error(testData []models.SLPoint) float64 {
	err := 0.0
	for _, v := range testData {
		p, _ := r.Predict(v.X)
		err = err + (p-v.Y)*(p-v.Y)
	}
	return math.Sqrt(err) / float64(len(testData))
}
func (r *RevezEstimator) Predict(p models.Vector) (float64, error) {
	// first we seek the closest point
	for i, pt := range r.Vectors {
		if models.L1Norm(pt, p) == 0 {
			return r.state[i], nil
		}
	}
	return 0, errors.New("the point is outside the learning domain")
}

func (r *RevezEstimator) Average(convexPart []float64, l models.SLPoint) {
	r.Step++
	ht := r.Smoothing(r.Step)
	for j, point := range r.Vectors {
		tmp := make([]float64, len(l.X))
		for i, v := range l.X {
			tmp[i] = (point[i] - v) / ht
		}
		tmpKer := r.Kernel(tmp) / math.Pow(ht, len(l.X))
		r.state[j] = convexPart[j] - r.Rate(r.Step)*(tmpKer*r.state[j]-l.Y*tmpKer)
	}
}

func (r *RevezEstimator) Compute(p models.SLPoint) {
	r.Average(r.state, p)
}

func (r RevezEstimator) State() models.State {
	return EstimatorState{State: r.state, version: r.Step}
}

func NewRevezEstimator(points []models.Vector) (*RevezEstimator, error) {
	dim := len(points[0])
	e := RevezEstimator{
		Vectors:   points,
		state:     make([]float64, len(points)),
		Step:      0,
		Rate:      func(i int) float64 { return 1.0 / float64(i) },
		Smoothing: func(t int) float64 { return math.Pow(t, -1/(dim+2)) },
		Kernel:    models.GaussianKernel}
	return &e, nil
}
