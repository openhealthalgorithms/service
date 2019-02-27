package engine

import (
	"fmt"
	"math"
)

// MedicationConditions object
type MedicationConditions struct {
	Antihypertensive  *bool `json:"anti-hypertensive"`
	OralHypoglycaemic *bool `json:"oral-hypoglycaemic"`
	Insulin           *bool `json:"insulin"`
	LipidLowering     *bool `json:"lipid-lowering"`
	Antiplatelet      *bool `json:"anti-platelet"`
	AntiCoagulant     *bool `json:"anti-coagulant"`
	Bronchodilator    *bool `json:"bronchodilator"`
}

// BloodPressureCondition object
type BloodPressureCondition struct {
	Medications *MedicationConditions `json:"medications"`
	Diabetes    *bool                 `json:"diabetes"`
	SBP         *RangeInt             `json:"sbp"`
	DBP         *RangeInt             `json:"dbp"`
	Age         *RangeFloat           `json:"age"`
	Target      *string               `json:"target"`
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
func (b *BloodPressureGuidelines) Process(diabetes bool, sbp, dbp int, age float64, medications map[string]bool) (Response, error) {
	code := ""
	value := fmt.Sprintf("%d/%d", sbp, dbp)
	target := ""

	for _, g := range *b {
		for _, c := range *g.Conditions {
			ageFrom := 0.0
			ageTo := math.MaxFloat64

			if c.Age != nil {
				if c.Age.From != nil {
					ageFrom = *c.Age.From
				}
				if c.Age.To != nil {
					ageTo = *c.Age.To
				}
			}

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

			conditionMedication := true
			if c.Medications != nil {
				if c.Medications.Antihypertensive != nil && *c.Medications.Antihypertensive != medications["anti-hypertensive"] {
					conditionMedication = false
				}
			}

			conditionDiabetes := true
			if c.Diabetes != nil && *c.Diabetes != diabetes {
				conditionDiabetes = false
			}

			if conditionDiabetes && conditionMedication && (age >= ageFrom && age <= ageTo) && sbpFrom <= sbp && sbpTo >= sbp && dbpFrom <= dbp && dbpTo >= dbp {
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
