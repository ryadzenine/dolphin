package estimators;

import (
    "github.com/ryadzenine/dolphin/kernels"
    "math"
)

type LearningPoint struct{
    X []float64
    Y float64
}
type RegressionEstimator interface {
    ComputeDistributedStep(float64, LearningPoint) 
    ComputeStep( LearningPoint)

}

type RevezEstimator struct {
    point   []float64
    State   float64 
    step    int
    rate    func (int) float64
    smoothing func (int) float64
    kernel  func([]float64) float64
}

func NewRevezEstimator()(*RevezEstimator){
    e := RevezEstimator{
        point: []float64{0,0},
        State: 0.0,
        step: 0.0,
        rate: func(i int) float64 { return 1.0/math.Sqrt(float64(i)) },
        smoothing: func(i int) float64 {return 0.2},
        kernel: kernels.Gaussian}
    return &e
        
}

func (r *RevezEstimator) ComputeDistributedStep(convexPart float64, p LearningPoint){
    r.step++
    tmp := make([]float64, cap(p.X))
    ht := r.smoothing(r.step)
    for i, v := range p.X {
        tmp[i] = (r.point[i] - v)/ht
    }
    tmp_ker := r.kernel(tmp)/ht
    r.State = convexPart - r.rate(r.step)*(tmp_ker*r.State - p.Y*tmp_ker)  
}

func (r *RevezEstimator) ComputeStep(p LearningPoint) {
    r.ComputeDistributedStep(r.State, p)  
}
