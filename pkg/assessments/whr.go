package assessments

import (
	"fmt"
)

func GetWHR(waist, hip float64, gender string) (WHRAssessment, error) {
	resultCode := "WHR-0"

	whr := waist / hip
	currentWhr := fmt.Sprintf("%.2f", whr)
	target := "0.85"
	if gender == "m" {
		target = "0.9"
	}

	if whr >= 0.85 && gender == "f" {
		resultCode = "WHR-1"
	} else if whr >= 0.9 && gender == "m" {
		resultCode = "WHR-2"
	}

	whrObj := NewWHRAssessment(currentWhr, resultCode, target)

	return whrObj, nil
}

type WHRAssessment struct {
	WHR    string `structs:"value"`
	Code   string `structs:"code"`
	Target string `structs:"target"`
}

// NewWHRAssessment returns a BP object.
func NewWHRAssessment(whr, code, target string) WHRAssessment {
	return WHRAssessment{
		WHR:    whr,
		Code:   code,
		Target: target,
	}
}
