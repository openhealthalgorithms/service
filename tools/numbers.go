package tools

// Float64InBetween function
func Float64InBetween(num, from, to float64) bool {
	return (num >= from && num <= to)
}

// IntInBetween function
func IntInBetween(num, from, to int) bool {
	return (num >= from && num <= to)
}
