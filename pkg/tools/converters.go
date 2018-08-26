package tools

import (
	"math"
)

func ConvertCholesterol(cholesterol float64, unit string) int {
	if unit == "mgdl" {
		cholesterol = cholesterol * 0.02586
	}

	tmp := int(math.Floor(cholesterol)) - 4

	if tmp < 1 {
		return 0
	} else if tmp <= 4 {
		return tmp
	}

	return 4
}

func ConvertSbp(sbp int) int {
	if sbp < 140 {
		return 3
	} else if sbp >= 140 && sbp < 160 {
		return 2
	} else if sbp >= 160 && sbp < 180 {
		return 1
	}

	return 0
}

func ConvertAge(age float64) int {
	if age <= 18 {
		return 0
	} else if age < 50 {
		return 40
	} else if age < 60 {
		return 50
	} else if age < 70 {
		return 60
	}

	return 70
}
