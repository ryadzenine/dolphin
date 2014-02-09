package kernels;

import (
    "math"
)
func l2Norm(v []float64) float64{
    var res float64
    for _,val := range(v) {
        res = res + val*val
    }
    return math.Sqrt(res)
}

func Epanechnikov(v []float64) float64 {
        return math.Max(1-l2Norm(v), 0)
}

func Gaussian(v []float64) float64 {
    return math.Exp(-l2Norm(v))
}

func Naive(v []float64) float64 {
    return math.Min(l2Norm(v), 1)  
}
