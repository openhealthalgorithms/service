package assessments

func GetDiabetes(bsl float64, bslUnit, bslType string, diabetes bool) (DiabetesAssessment, error) {
	status := false
	resultCode := ""

	bslValue := bsl
	if bslUnit == "mg/dl" {
		bslValue = bsl / 18
	}

	if diabetes {
		status = true
		resultCode = "DM-4"
	} else {
		if bslValue > 7 {
			status = true
			resultCode = "DM-3"
		} else if bslValue > 6.1 {
			status = true
			resultCode = "DM-2"
		} else if bslType == "hba1c" {
			status = true
		}
	}

	diabetesObj := NewDiabetesAssessment(bslValue, resultCode, status)

	return diabetesObj, nil
}

type DiabetesAssessment struct {
	BSL    float64 `structs:"value"`
	Code   string  `structs:"code"`
	Status bool    `structs:"status"`
}

// NewDiabetesAssessment returns a BP object.
func NewDiabetesAssessment(bsl float64, code string, status bool) DiabetesAssessment {
	return DiabetesAssessment{
		BSL:    bsl,
		Code:   code,
		Status: status,
	}
}
