package models

type Point []float64

// Represents a supervised learning Point
type SLPoint struct {
	X Point   // the point
	Y float64 // The label
}
