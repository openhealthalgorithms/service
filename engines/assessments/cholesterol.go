package assessments

import (
	"fmt"
	"math"

	"github.com/openhealthalgorithms/service/tools"
)

// CholesterolGuidelinesFull object
type CholesterolGuidelinesFull struct {
	TotalCholesterol *CholesterolGuidelines `json:"total_cholesterol"`
	HDL              *CholesterolGuidelines `json:"hdl"`
	LDL              *CholesterolGuidelines `json:"ldl"`
	TG               *CholesterolGuidelines `json:"tg"`
}

// CholesterolCondition object
type CholesterolCondition struct {
	Medications *MedicationConditions `json:"medications"`
	Age         *RangeFloat           `json:"age"`
	CVD         *RangeFloat           `json:"cvd"`
	Range       *RangeFloat           `json:"range"`
	Target      *string               `json:"target"`
}

// CholesterolConditions slice
type CholesterolConditions []CholesterolCondition

// CholesterolGuideline object
type CholesterolGuideline struct {
	Category   *string                `json:"category"`
	Definition *string                `json:"definition"`
	Conditions *CholesterolConditions `json:"conditions"`
	Code       *string                `json:"code"`
}

// CholesterolGuidelines slice
type CholesterolGuidelines []CholesterolGuideline

// Process function
func (b *CholesterolGuidelines) Process(cvd, age, chol float64, cholUnit, cholType string, medications map[string]bool) (Response, error) {
	cholesterol := tools.CalculateCholMMOLValue(chol, cholUnit)

	code := ""
	value := fmt.Sprintf("%.1f%s", chol, cholUnit)
	target := ""

	for _, g := range *b {
		for _, c := range *g.Conditions {
			ageFrom := 0.0
			ageTo := math.MaxFloat64

			cvdFrom := 0.0
			cvdTo := math.MaxFloat64

			cholFrom := 0.0
			cholTo := math.MaxFloat64

			if c.Age != nil {
				if c.Age.From != nil {
					ageFrom = *c.Age.From
				}
				if c.Age.To != nil {
					ageTo = *c.Age.To
				}
			}

			// fmt.Println("AGE =>", ageFrom, ageTo)

			if c.CVD != nil {
				if c.CVD.From != nil {
					cvdFrom = *c.CVD.From
				}
				if c.CVD.To != nil {
					cvdTo = *c.CVD.To
				}
			}

			// fmt.Println("CVD =>", cvdFrom, cvdTo)

			if c.Range != nil {
				if c.Range.From != nil {
					cholFrom = tools.CalculateCholMMOLValue(*c.Range.From, *c.Range.Unit)
				}
				if c.Range.To != nil {
					cholTo = tools.CalculateCholMMOLValue(*c.Range.To, *c.Range.Unit)
				}
			}

			conditionMedication := true
			if c.Medications != nil {
				if c.Medications.LipidLowering != nil && *c.Medications.LipidLowering != medications["lipid-lowering"] {
					conditionMedication = false
				}
			}
			// fmt.Println("CHOL =>", cholFrom, cholTo, cholUnit)

			if conditionMedication && (age >= ageFrom && age <= ageTo) && (cvd >= cvdFrom && cvd <= cvdTo) && (cholesterol >= cholFrom && cholesterol <= cholTo) {
				code = *g.Code
				target = *c.Target
				break
			}
		}
		if len(code) > 0 {
			break
		}
	}

	return GetResponse(cholType, code, value, target)
}
