package models

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

func L2Norm(v Vector) float64 {
	var res float64
	for _, val := range v {
		res = res + val*val
	}
	return math.Sqrt(res)
}
func L22Norm(v, u Vector) float64 {
	var res float64
	for key, val := range v {
		res = res + math.Pow(val-u[key], 2)
	}
	return math.Sqrt(res)
}
func L1Norm(x Vector, y Vector) float64 {
	if len(x) > len(y) {
		return -1
	}
	var res float64
	for i, val := range x {
		res = res + math.Abs(val-y[i])
	}
	return res
}

func ParseLearningPoint(s string) SLPoint {
	sl := strings.Split(s, ";")
	fl := make([]float64, cap(sl))
	for i, s := range sl {
		v, _ := strconv.ParseFloat(s, 64)
		fl[i] = v
	}
	return SLPoint{Y: fl[0], X: fl[1:len(fl)]}
}

// buildMeshCoordinate will build all the values that the coordinats of the
// meshing will take
func buildMeshCoordinate(inf, sup, step float64) []float64 {
	nbVectors := int(math.Ceil(sup - inf/float64(step)))
	coordinates := make([]float64, 0, nbVectors+1)
	// Dans un premier temps on construit les coordonées de la discrétisation
	for j := inf + float64(step)/2; j < sup; j = j + float64(step) {
		coordinates = append(coordinates, j)
	}
	return coordinates
}
func buildVector(track []int, coordinates []float64) Vector {
	pt := make([]float64, 0, len(track))
	for _, v := range track {
		pt = append(pt, coordinates[v])
	}
	return Vector(pt)
}
func incTrack(track []int, d int) ([]int, bool) {
	acc := 1
	for i, v := range track {
		acc = acc + v*int(math.Pow(float64(d), float64(i)))
	}
	tr := track[0:0]
	for i := 0; i < len(track); i++ {
		tr = append(tr, acc%d)
		acc = int(acc / d)
	}
	acc = 0
	for i, v := range tr {
		acc = acc + v*int(math.Pow(float64(d), float64(i)))
	}
	if acc == 0 {
		return nil, false
	}
	return tr, true
}

// This function will build a meshing of the compact [ @inf , @sup ]^@d
// such that the euclidian distance between two points of the mesh is less
// then @step
func MeshEvalVectors(inf, sup, step float64, d int) ([]Vector, error) {
	if inf >= sup {
		return nil, errors.New("inf should be strictly smaller than sup")
	}
	if float64(sup-inf) < step {
		return nil, errors.New("sup - inf should be strictly greater than step")
	}
	//
	nbVectors := int(math.Ceil(sup - inf/step))
	points := make([]Vector, 0, (nbVectors+1)*d)
	// on va se servir de ce tableau pour construire
	coordinates := buildMeshCoordinate(inf, sup, step)
	track := make([]int, d)
	ok := true
	for ok {
		// on rajoute le point
		points = append(points, buildVector(track, coordinates))
		track, ok = incTrack(track, len(coordinates))
	}
	return points, nil
}
