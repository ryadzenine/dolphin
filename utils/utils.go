package utils

import (
	"strings"

	"github.com/ryadzenine/dolphin/models"
)

func ParseData(source []byte) []models.SLPoint {
	var data []models.SLPoint
	for _, s := range strings.Split(string(source), "\n") {
		data = append(data, models.ParseLearningPoint(s))
	}
	return data
}
