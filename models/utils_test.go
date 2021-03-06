package models

import (
	"testing"
)

func TestParseLearningPoint(t *testing.T) {
	l := ParseLearningPoint("5.28;1.65;2")

	if l.Y != 5.28 || l.X[0] != 1.65 || l.X[1] != 2 {
		t.Error("Parse is not working correctly")
	}
	l = ParseLearningPoint("5.28;1.65")
	if l.Y != 5.28 || l.X[0] != 1.65 {
		t.Error("Parse is not working correctly")
	}

}
func TestBuildMeshCoordinate(t *testing.T) {
	coo := buildMeshCoordinate(0, 1, 0.25)
	if len(coo) != 4 {
		t.Error("the lenght of the build coordinates for parameters ", 0, 1, 0.25, coo)
	}
}
