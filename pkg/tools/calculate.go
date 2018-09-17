package tools

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
