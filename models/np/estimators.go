package np

import (
  "errors"
  "github.com/ryadzenine/dolphin/models"
  "math"
)

type EstimatorState struct {
  Points  []models.Point
  State   []float64
  version int
}

func (r EstimatorState) Values() []float64 {
  return r.State
}
func (r EstimatorState) Version() int {
  return r.version
}

type RegressionEstimator interface {
  Predict(models.Point) (float64, error)
  ComputeDistributedStep(float64, models.SLPoint)
  ComputeStep(models.SLPoint)
  State() models.State
}

type RevezEstimator struct {
  Points    []models.Point
  state     []float64
  Step      int
  Rate      func(int) float64
  Smoothing func(int) float64
  Kernel    func(models.Point) float64
}

func (r *RevezEstimator) L2Error(testData []models.SLPoint) float64 {
  err := 0.0
  for _, v := range testData {
    p, _ := r.Predict(v.X)
    err = err + (p-v.Y)*(p-v.Y)
  }
  return math.Sqrt(err) / float64(len(testData))
}
func (r *RevezEstimator) Predict(p models.Point) (float64, error) {
  // first we seek the closest point
  for i, pt := range r.Points {
    if l1Norm(pt, p) == 0 {
      return r.state[i], nil
    }
  }
  return 0, errors.New("the point is outside the learning domain")
}

func (r *RevezEstimator) ComputeDistributedStep(convexPart []float64, l models.SLPoint) {
  r.Step++
  ht := r.Smoothing(r.Step)
  for j, point := range r.Points {
    tmp := make([]float64, len(l.X))
    for i, v := range l.X {
      tmp[i] = (point[i] - v) / ht
    }
    tmp_ker := r.Kernel(tmp) / ht
    r.state[j] = convexPart[j] - r.Rate(r.Step)*(tmp_ker*r.state[j]-l.Y*tmp_ker)
  }
}

func (r *RevezEstimator) ComputeStep(p models.SLPoint) {
  r.ComputeDistributedStep(r.state, p)
}

func (r RevezEstimator) State() models.State {
  return EstimatorState{Points: r.Points, State: r.state, version: r.Step}
}

func NewRevezEstimator(points []models.Point, smooth float64) (*RevezEstimator, error) {
  e := RevezEstimator{
    Points:    points,
    state:     make([]float64, len(points)),
    Step:      0,
    Rate:      func(i int) float64 { return 1.0 / math.Sqrt(float64(i)) },
    Smoothing: func(i int) float64 { return smooth },
    Kernel:    GaussianKernel}
  return &e, nil
}
