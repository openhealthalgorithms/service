package assessments

import (
	"fmt"
)

func GetBMI(weight, height float64) (BMIAssessment, error) {
	resultCode := ""

	bmi := weight / (height * height)
	currentBmi := fmt.Sprintf("%.2f", bmi)
	target := "18.5 - 24.9"

	if bmi < 18.5 {
		resultCode = "BMI-1"
	} else if bmi < 25 {
		resultCode = "BMI-0"
	} else if bmi < 30 {
		resultCode = "BMI-2"
	} else {
		resultCode = "BMI-3"
	}

	bmiObj := NewBMIAssessment(currentBmi, resultCode, target)

	return bmiObj, nil
}

type BMIAssessment struct {
	BMI    string `structs:"value"`
	Code   string `structs:"code"`
	Target string `structs:"target"`
}

// NewBMIAssessment returns a BP object.
func NewBMIAssessment(bmi, code, target string) BMIAssessment {
	return BMIAssessment{
		BMI:    bmi,
		Code:   code,
		Target: target,
	}
}
