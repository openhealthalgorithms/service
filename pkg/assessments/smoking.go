package assessments

func GetSmoking(current, ex, quitYear bool) (SmokingAssessment, error) {
	status := false
	smokingCalc := false
	resultCode := ""

	if current {
		status = true
		smokingCalc = true
		resultCode = "SM-1"
	} else if ex && quitYear {
		smokingCalc = true
		resultCode = "SM-2"
	} else if ex {
		resultCode = "SM-3"
	} else {
		resultCode = "SM-4"
	}

	smokingObj := NewSmokingAssessment(resultCode, status, smokingCalc)

	return smokingObj, nil
}

type SmokingAssessment struct {
	Code        string `structs:"code"`
	Status      bool   `structs:"status"`
	SmokingCalc bool   `structs:"smoking_calc"`
}

// NewSmokingAssessment returns a BP object.
func NewSmokingAssessment(code string, status, smokingCalc bool) SmokingAssessment {
	return SmokingAssessment{
		Code:        code,
		Status:      status,
		SmokingCalc: smokingCalc,
	}
}
