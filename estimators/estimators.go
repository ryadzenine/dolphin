package estimators;

import (
    "github.com/ryadzenine/dolphin/kernels"
)

type RegressionEstimator interface {
    ComputeDistributedStep(float64, []float64, float64) 
    ComputeStep([]float64, float64)

}
type RevezEstimator struct {
    point   []float64
    state   float64 
    step    int
    rate    func (int) float64
    smoothing func (int) float64
    kernel  func([]float64) float64
}
func NewRevezEstimator()(*RevezEstimator){
    e := RevezEstimator{
        []float64{0,0},
        0.0,
        0.0,
        func(i int) float64 { return 1.0/float64(i) },
        func(i int) float64 {return 0.2},
        kernels.Gaussian}
    return &e
        
}
func (r *RevezEstimator) ComputeDistributedStep(convexPart float64, x []float64, y float64){
    r.step++
    tmp := make([]float64, cap(x))
    ht := r.smoothing(r.step)
    for i, v := range x {
        tmp[i] = (r.point[i] - v)/ht
    }
    tmp_ker := r.kernel(tmp)/ht
    r.state = convexPart - r.rate(r.step)*(tmp_ker*r.state - y*tmp_ker)  

}

func (r *RevezEstimator) ComputeStep(x []float64,y float64) {
    r.ComputeDistributedStep(r.state, x, y)  
}
