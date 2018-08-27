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

func ConvertLength(length float64, unit string) float64 {
	result := length

	if unit == "cm" {
		result = length / 100
	} else if unit == "ft" {
		result = length * 3.28084
	} else if unit == "inch" {
		result = length * 39.3701
	}

	return result
}

func ConvertWeight(weight float64, unit string) float64 {
	result := weight

	if unit == "lb" {
		result = weight / 0.45359237
	}

	return result
}
