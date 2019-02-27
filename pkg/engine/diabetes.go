package engine

import (
	"fmt"
	"math"
	"strings"

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
	Medications *MedicationConditions `json:"medications"`
	HxDiabetes  *bool                 `json:"hx_diabetes"`
	BloodSugar  *BloodSugarCondition  `json:"blood_sugar"`
	Target      *string               `json:"target"`
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
func (b *DiabetesGuidelines) Process(hxDiabetes bool, bsFromInput float64, bsType, unit string, medications map[string]bool) (Response, error) {
	code := ""
	value := ""
	target := ""

	from := bsFromInput
	bslCheck := false

	if bsType == "fasting" || bsType == "random" || unit != "%" {
		from = tools.CalculateMMOLValue(bsFromInput, unit)
		bslCheck = true
	}

	for _, g := range *b {
		for _, c := range *g.Conditions {
			bsFrom := 0.0
			bsTo := math.MaxFloat64
			givenBsType := ""

			if c.BloodSugar != nil {
				if c.BloodSugar.From != nil {
					if bslCheck {
						bsFrom = tools.CalculateMMOLValue(*c.BloodSugar.From, *c.BloodSugar.Unit)
					} else {
						bsFrom = *c.BloodSugar.From
					}
				}
				if c.BloodSugar.To != nil {
					if bslCheck {
						bsTo = tools.CalculateMMOLValue(*c.BloodSugar.To, *c.BloodSugar.Unit)
					} else {
						bsTo = *c.BloodSugar.To
					}
				}
				givenBsType = *c.BloodSugar.Type
			}

			conditionMedication := true
			if c.Medications != nil {
				if (c.Medications.OralHypoglycaemic != nil && *c.Medications.OralHypoglycaemic != medications["oral-hypoglycaemic"]) || (c.Medications.Insulin != nil && *c.Medications.Insulin != medications["insulin"]) {
					conditionMedication = false
				}
			}

			conditionHxDiabetes := true
			if c.HxDiabetes != nil && *c.HxDiabetes != hxDiabetes {
				conditionHxDiabetes = false
			}

			if conditionHxDiabetes && conditionMedication && bsFrom <= from && bsTo >= from && strings.ToLower(bsType) == strings.ToLower(givenBsType) {
				code = *g.Code
				target = *c.Target
				value = fmt.Sprintf("%.1f%s", bsFromInput, unit)
				break
			}
		}
		if len(code) > 0 {
			break
		}
	}

	return GetResponse("diabetes", code, value, target)
}
