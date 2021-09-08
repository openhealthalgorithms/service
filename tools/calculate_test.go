package tools

import (
	"testing"
)

func TestCalculateAge(t *testing.T) {
	type ageTest struct {
		value  float64
		unit   string
		result float64
	}

	var ageTests = []ageTest{
		{25, "year", 25},
		{12, "month", 1},
		{17, "week", 0.33},
		{715, "day", 1.96},
		{17, "hours", 17},
	}

	for _, at := range ageTests {
		actual := CalculateAge(at.value, at.unit)
		if actual != at.result {
			t.Errorf("CalculateAge(%f, %s): expected %f, actual %f", at.value, at.unit, at.result, actual)
		}
	}
}
