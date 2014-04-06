package utils

import (
	"github.com/ryadzenine/dolphin/models"
	"strings"
)

func ParseData(source []byte) []models.SLPoint {
	data := make([]models.SLPoint, 0)
	for _, s := range strings.Split(string(source), "\n") {
		data = append(data, models.ParseLearningPoint(s))
	}
	return data
}
