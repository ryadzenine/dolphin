package estimators;

import (
    "math"
    "errors"
    "github.com/ryadzenine/dolphin/mpi"
)
type Point []float64

type LearningPoint struct{
    X Point
    Y float64
}

type RegressionEstimator interface {
    Predict(Point) (float64,error) 
    ComputeDistributedStep(float64, LearningPoint) 
    ComputeStep(LearningPoint)
    State() State
}

type RevezEstimator struct {
    Points  []Point
    state   []float64 
    Step    int
    Rate    func (int) float64
    Smoothing func (int) float64
    Kernel  func(Point) float64
}
func (r *RevezEstimator) L2Error(testData []LearningPoint) float64 {
    err := 0.0
    for _,v := range testData{
        p,_ := r.Predict(v.X)
        err = err + (p - v.Y)*(p - v.Y)
    }
    return math.Sqrt(err) 
}
func (r *RevezEstimator) Predict(p Point) (float64, error) {
    // first we seek the closest point 
    for i,pt := range r.Points  {
        if l1Norm(pt, p) < r.Smoothing(0)*float64(len(pt))/2 {
            return r.state[i], nil
        }
    }
    return 0,errors.New("the point is outside the learning domain") 

}
func (r *RevezEstimator) ComputeDistributedStep(convexPart []float64, l LearningPoint){
    r.Step++
    ht := r.Smoothing(r.Step)
    for j, point := range r.Points {
        tmp := make([]float64, len(l.X))
        for i,v := range l.X {
            tmp[i] = (point[i] - v)/ht
        }
        tmp_ker := r.Kernel(tmp)/ht
        r.state[j] = convexPart[j] - r.Rate(r.Step)*(tmp_ker*r.state[j] - l.Y*tmp_ker)
    }
}

func (r *RevezEstimator) ComputeStep(p LearningPoint) {
    r.ComputeDistributedStep(r.state, p)  
}

func (r RevezEstimator) State() State {
    return EstimatorState{ Points: r.Points, State: r.state, version: r.Step} 
}

func NewRevezEstimator(points []Point, step float64)(*RevezEstimator, error){
    e := RevezEstimator{
        Points: points,
        state: make([]float64, len(points)),
        Step: 0,
        Rate: func(i int) float64 { return 1.0/math.Sqrt(float64(i)) },
        Smoothing: func(i int) float64 {return step},
        Kernel: GaussianKernel}
    return &e, nil
}

type State interface {
    mpi.Versionable
    Values() []float64
}
type States []State

func (states States) ComputeAgregation() []float64{
    N := len(states)
    agg := make([]float64, len(states[0].Values())) 
    for _,s := range states {
        for j,v := range s.Values() {
            agg[j]= agg[j] + v
        }
    }
    for i,v := range agg {
        agg[i] = v/float64(N)
    }
    return agg 
}

type EstimatorState struct{
    Points []Point 
    State []float64
    version int
}
func (r EstimatorState) Values() []float64 {
    return r.State
}
func (r EstimatorState) Version() int {
    return r.version 
}




