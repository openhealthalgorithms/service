package assessments

func GetHighRisks(sbp, dbp int, age float64, current map[string]bool) (HighRisksAssessment, error) {
	hrc := []string{"CVD", "CAD", "AMI", "HEART ATTACK", "CVA", "TIA", "STROKE", "CKD", "PVD"}
	return GetHighRisksWithHrc(sbp, dbp, age, current, hrc)
}

func GetHighRisksWithHrc(sbp, dbp int, age float64, current map[string]bool, hrc []string) (HighRisksAssessment, error) {
	hasHighRisk := false
	resultCode := ""
	reason := ""

	for _, h := range hrc {
		if _, ok := current[h]; ok {
			hasHighRisk = true
			resultCode = "HR-0"
			reason = h
		}
	}

	if !hasHighRisk {
		if sbp > 200 || dbp > 120 {
			hasHighRisk = true
			resultCode = "HR-1"
		} else if age < 40 && (sbp >= 140 || dbp >= 90) {
			resultCode = "HR-2"
		}
	}

	highRisksObj := NewHighRisksAssessment(hasHighRisk, reason, resultCode)

	return highRisksObj, nil
}

type HighRisksAssessment struct {
	Status bool   `structs:"status"`
	Reason string `structs:"reason"`
	Code   string `structs:"code"`
}

// NewHighRisksAssessment returns a BP object.
func NewHighRisksAssessment(status bool, reason, code string) HighRisksAssessment {
	return HighRisksAssessment{
		Status: status,
		Reason: reason,
		Code:   code,
	}
}
