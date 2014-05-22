package regression

import (
	"math"

	"github.com/ryadzenine/dolphin/models"
	"github.com/ryadzenine/dolphin/models/optimization"
)

type SVM struct {
	alpha []float64
	data  models.SLDataset
	pen   float64
	init  optimization.Initializer
}

func (s SVM) Model() []float64   { return s.alpha }
func (s SVM) Update(v []float64) { s.alpha = v }
func (s SVM) Initialize()        { s.alpha = s.init() }

func (s SVM) f(v models.Vector) float64 {
	acc := 0.
	for key, value := range s.alpha {
		acc += value * math.Exp(-models.L22Norm(v, s.data[key].X))
	}
	return acc
}

func (s *SVM) Gradient(index int) float64 {
	return s.data[index].Y*s.alpha[index]*math.Exp(-models.L2Norm(s.data[index].X)) + s.pen*s.alpha[index]
}

func (s *SVM) Loss(vector models.SLDataset) float64 {
	acc := 0.
	for _, data := range vector {
		acc += math.Max(math.Abs(1-data.Y*s.f(data.X)), 0)
	}
	return acc
}
