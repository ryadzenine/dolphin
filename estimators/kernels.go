package estimators;

import (
    "math"
)

func l2Norm(v Point) float64{
    var res float64
    for _,val := range(v) {
        res = res + val*val
    }
    return math.Sqrt(res)
}
func l1Norm(x Point,y Point) float64{
    if len(x) > len(y) {
        return -1
    }
    var res float64
    for i,val := range(x) {
        res = res + math.Abs(val - y[i])
    }
    return res
}
func EpanechnikovKernel(v Point) float64 {
        return math.Max(1-l2Norm(v), 0)
}

func GaussianKernel(v Point) float64 {
    return math.Exp(-l2Norm(v))
}

func NaiveKernel(v Point) float64 {
    return math.Min(l2Norm(v), 1)  
}

var Kernels map[string] func(Point) float64 = make(map[string] func(Point) float64 , 4)
func Init(){
    Kernels["naive"] = NaiveKernel
    Kernels["gaussian"] = GaussianKernel
    Kernels["Epanechnikov"] = EpanechnikovKernel
}
