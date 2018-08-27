package assessments

import (
	"fmt"
)

func GetBP(sbp, dbp int, diabetes bool) (BPAssessment, error) {
	resultCode := ""

	currentBp := fmt.Sprintf("%d/%d", sbp, dbp)
	target := currentBp

	if diabetes {
		if sbp > 130 {
			resultCode = "BP-3B"
			target = "130/80"
		} else {
			resultCode = "BP-3A"
			target = "130/80"
		}
	} else {
		if sbp > 160 {
			resultCode = "BP-2"
			target = "140/90"
		} else if sbp > 140 {
			resultCode = "BP-1B"
			target = "140/90"
		} else if sbp <= 140 && sbp >= 120 {
			resultCode = "BP-1A"
			target = "140/90"
		} else {
			resultCode = "BP-0"
			target = "140/90"
		}
	}

	bpObj := NewBPAssessment(currentBp, sbp, dbp, resultCode, target)

	return bpObj, nil
}

type BPAssessment struct {
	BP     string `structs:"bp"`
	SBP    int    `structs:"sbp"`
	DBP    int    `structs:"dbp"`
	Code   string `structs:"code"`
	Target string `structs:"target"`
}

// NewBPAssessment returns a BP object.
func NewBPAssessment(bp string, sbp, dbp int, code, target string) BPAssessment {
	return BPAssessment{
		BP:     bp,
		SBP:    sbp,
		DBP:    dbp,
		Code:   code,
		Target: target,
	}
}
