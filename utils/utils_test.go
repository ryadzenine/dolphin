package utils

import (
	"fmt"
	"testing"
)

func TestParseData(t *testing.T) {
	bytes := []byte{'0', ';', '1', '\n', '2', ';', '3'}
	data := ParseData(bytes)
	fmt.Println(data)
	if len(data) != 2 {
		t.Error("the data should be of length 2")
	}
	if data[0].X[0] != 1 || data[0].Y != 0 {
		t.Error("error parsing the first point")
	}
	if data[1].X[0] != 3 || data[1].Y != 2 {
		t.Error("error parsing the second point")
	}
}
