package tools

import (
	"strings"
)

// CalculateAge to year
func CalculateAge(age float64, unit string) float64 {
	calAge := age
	if unit != "year" {
		switch unit {
		case "month":
			calAge = age / 12
		case "week":
			calAge = age / 52
		case "day":
			calAge = age / 365
		}
	}

	return calAge
}

// CalculateExercise to minutes and moderate type weekle
func CalculateExercise(value int, unit, frequency, intensity string) int {
	converted := value

	if unit == "hours" {
		converted = converted * 60
	}

	if intensity == "high" {
		converted = converted * 2
	} else if intensity == "low" {
		converted = (int)(converted / 2)
	}

	if frequency == "daily" {
		converted = converted * 7
	} else if frequency == "monthly" {
		converted = (int)(converted / 4)
	}

	return converted
}

// CalculateDietConsumption in servings weekly
func CalculateDietConsumption(value int, frequency string) int {
	converted := value

	if frequency == "daily" {
		converted = converted * 7
	} else if frequency == "monthly" {
		converted = (int)(converted / 4)
	}

	return converted
}

// CalculateMMOLValue to convert value into mmol/L
func CalculateMMOLValue(value float64, unit string) float64 {
	if strings.ToLower(unit) == "mg/dl" {
		return value / 18
	}

	return value
}

// CalculateLength to convert units
func CalculateLength(value float64, unit, toUnit string) float64 {
	result := value
	denominator := 1.0

	if unit == "cm" && toUnit == "m" {
		result = result / 100
	} else if unit == "m" && toUnit == "cm" {
		result = result * 100
	} else if unit == "cm" && toUnit == "inch" {
		result = result / 2.54
	} else if unit == "inch" && toUnit == "cm" {
		result = result * 2.54
	} else if unit == "cm" && toUnit == "ft" {
		result = result / 30.48
	} else if unit == "ft" && toUnit == "cm" {
		result = result * 30.48
	} else if unit == "m" && toUnit == "ft" {
		result = result * 3.281
	} else if unit == "ft" && toUnit == "m" {
		result = result / 3.281
	} else if unit == "m" && toUnit == "inch" {
		result = result * 39.37
	} else if unit == "inch" && toUnit == "m" {
		result = result / 39.37
	} else if unit == "ft" && toUnit == "inch" {
		result = result / 12
	} else if unit == "inch" && toUnit == "ft" {
		result = result * 12
	}

	return result / denominator
}

// CalculateAlcoholConsumption in servings weekly
func CalculateAlcoholConsumption(value float64, frequency string) float64 {
	converted := value

	if frequency == "daily" {
		converted = converted * 7
	} else if frequency == "monthly" {
		converted = converted / 4
	}

	return converted
}
