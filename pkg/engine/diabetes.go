package engine

import (
	"math"

	"github.com/openhealthalgorithms/service/pkg/tools"
)

// BloodSugarCondition object
type BloodSugarCondition struct {
	Type *string  `json:"type"`
	From *float64 `json:"from"`
	To   *float64 `json:"to"`
	Unit *string  `json:"unit"`
}

// DiabetesCondition object
type DiabetesCondition struct {
	HxDiabetes *bool                `json:"hx_diabetes"`
	BloodSugar *BloodSugarCondition `json:"blood_sugar"`
	Target     *string              `json:"target"`
}

// DiabetesConditions slice
type DiabetesConditions []DiabetesCondition

// DiabetesGuideline object
type DiabetesGuideline struct {
	Category   *string             `json:"category"`
	Definition *string             `json:"definition"`
	Conditions *DiabetesConditions `json:"conditions"`
	Code       *string             `json:"code"`
}

// DiabetesGuidelines slice
type DiabetesGuidelines []DiabetesGuideline

// Process function
func (b *DiabetesGuidelines) Process(hxDiabetes bool, bsFromInput float64, bsType, unit string) (Response, error) {
	code := ""
	value := ""
	target := ""

	from := tools.CalculateMMOLValue(bsFromInput, unit)

	for _, g := range *b {
		for _, c := range *g.Conditions {
			bsFrom := 0.0
			bsTo := math.MaxFloat64

			if c.BloodSugar != nil {
				if c.BloodSugar.From != nil {
					bsFrom = tools.CalculateMMOLValue(*c.BloodSugar.From, *c.BloodSugar.Unit)
				}
				if c.BloodSugar.To != nil {
					bsTo = tools.CalculateMMOLValue(*c.BloodSugar.To, *c.BloodSugar.Unit)
				}
			}

			conditionHxDiabetes := true
			if c.HxDiabetes != nil && *c.HxDiabetes != hxDiabetes {
				conditionHxDiabetes = false
			}

			if conditionHxDiabetes && bsFrom <= from && bsTo >= from {
				code = *g.Code
				target = *c.Target
				value = *g.Category
				break
			}
		}
		if len(code) > 0 {
			break
		}
	}

	return GetResponse("diabetes", code, value, target)
}
