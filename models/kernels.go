package models

import (
	"math"
)

func EpanechnikovKernel(v Vector) float64 {
	return math.Max(1-L2Norm(v), 0)
}

func GaussianKernel(v Vector) float64 {
	return math.Exp(-L2Norm(v))
}

func NaiveKernel(v Vector) float64 {
	return math.Min(L2Norm(v), 1)
}

var Kernels = make(map[string]func(Vector) float64, 4)

func Init() {
	Kernels["naive"] = NaiveKernel
	Kernels["gaussian"] = GaussianKernel
	Kernels["Epanechnikov"] = EpanechnikovKernel
}
