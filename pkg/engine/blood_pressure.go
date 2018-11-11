package engine

import (
	"fmt"
	"math"
)

// BloodPressureCondition object
type BloodPressureCondition struct {
	Diabetes *bool     `json:"diabetes"`
	SBP      *RangeInt `json:"sbp"`
	DBP      *RangeInt `json:"dbp"`
	Target   *string   `json:"target"`
}

// BloodPressureConditions slice
type BloodPressureConditions []BloodPressureCondition

// BloodPressureGuideline object
type BloodPressureGuideline struct {
	Category   *string                  `json:"category"`
	Definition *string                  `json:"definition"`
	Conditions *BloodPressureConditions `json:"conditions"`
	Code       *string                  `json:"code"`
}

// BloodPressureGuidelines slice
type BloodPressureGuidelines []BloodPressureGuideline

// Process function
func (b *BloodPressureGuidelines) Process(diabetes bool, sbp, dbp int) (Response, error) {
	code := ""
	value := fmt.Sprintf("%d/%d", sbp, dbp)
	target := ""

	for _, g := range *b {
		for _, c := range *g.Conditions {
			sbpFrom := 0
			sbpTo := math.MaxInt32

			if c.SBP != nil {
				if c.SBP.From != nil {
					sbpFrom = *c.SBP.From
				}
				if c.SBP.To != nil {
					sbpTo = *c.SBP.To
				}
			}

			dbpFrom := 0
			dbpTo := math.MaxInt32

			if c.DBP != nil {
				if c.DBP.From != nil {
					dbpFrom = *c.DBP.From
				}
				if c.DBP.To != nil {
					dbpTo = *c.DBP.To
				}
			}

			conditionDiabetes := true
			if c.Diabetes != nil && *c.Diabetes != diabetes {
				conditionDiabetes = false
			}

			if conditionDiabetes && sbpFrom <= sbp && sbpTo >= sbp && dbpFrom <= dbp && dbpTo >= dbp {
				code = *g.Code
				target = *c.Target
				break
			}
		}
		if len(code) > 0 {
			break
		}
	}

	return GetResponse("blood pressure", code, value, target)
}
