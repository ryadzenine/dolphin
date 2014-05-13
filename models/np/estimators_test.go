package np

import (
	"testing"

	"github.com/ryadzenine/dolphin/models"
)

func TestRevezEstimatorPredict(t *testing.T) {
	// first we setup a face revez estimator with two
	// points
	points := []models.Point{models.Point{0}, models.Point{0.5}}
	r, _ := NewRevezEstimator(points)
	r.state = []float64{0, 1}
	res, _ := r.Predict(models.Point{0.2})
	if res != 0 {
		t.Error("Prediction failled")
	}
}

func TestRevezEstimatorL2Error(t *testing.T) {
	points := []models.Point{models.Point{0}, models.Point{0.5}}
	r, _ := NewRevezEstimator(points)
	r.state = []float64{0, 1}
	res := r.L2Error([]models.SLPoint{
		models.SLPoint{X: models.Point{0.2}, Y: 1}})
	if res != 1 {
		t.Error("L2Error Failled")
	}
}

func TestRevezEstimatorComputeDistributedStep(t *testing.T) {
	points := []models.Point{models.Point{0}, models.Point{0.5}}
	r, _ := NewRevezEstimator(points)
	r.state = []float64{0, 1}
	r.ComputeDistributedStep(
		[]float64{0.25, 0.75},
		models.SLPoint{X: models.Point{0.25}, Y: 0.3})
}
func TestComputeAgregation(t *testing.T) {
	st := EstimatorState{
		Points: []models.Point{
			models.Point{1, 0},
			models.Point{0, 1},
		},
		State:   []float64{0, 1},
		version: 1,
	}
	st1 := EstimatorState{
		Points: []models.Point{
			models.Point{1, 0},
			models.Point{0, 1},
		},
		State:   []float64{2.5, 2},
		version: 1,
	}
	st3 := EstimatorState{
		Points: []models.Point{
			models.Point{1, 0},
			models.Point{0, 1},
		},
		State:   []float64{0.5, 3},
		version: 1,
	}
	states := models.States{st, st1, st3}
	r := states.ComputeAgregation()
	if r[0] != 1 || r[1] != 2.0 {
		t.Error("Compute agregation failled", r)
	}
}
