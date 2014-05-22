package models

type Vector []float64

// Represents a supervised learning Point
type SLPoint struct {
	X Vector  // the point
	Y float64 // The label
}

type Dataset []Vector
type SLDataset []SLPoint
