package assessments

import (
	"fmt"
	"math"

	"github.com/openhealthalgorithms/service/tools"
)

// BodyCompositionGuideline object
type BodyCompositionGuideline struct {
	BMI       *BMIGuidelines       `json:"bmi"`
	WaistCirc *WaistCircGuidelines `json:"waist_circ"`
	WHR       *WHRGuidelines       `json:"whr"`
	BodyFat   *BodyFatGuidelines   `json:"body_fat"`
}

// RangeInt object
type RangeInt struct {
	From   *int    `json:"from"`
	To     *int    `json:"to"`
	Unit   *string `json:"unit"`
	Target *string `json:"target"`
}

// RangeFloat object
type RangeFloat struct {
	From   *float64 `json:"from"`
	To     *float64 `json:"to"`
	Unit   *string  `json:"unit"`
	Target *string  `json:"target"`
}

// BMIConditions slice
type BMIConditions []RangeFloat

// BMIGuideline object
type BMIGuideline struct {
	Category   *string        `json:"category"`
	Definition *string        `json:"definition"`
	Conditions *BMIConditions `json:"conditions"`
	Code       *string        `json:"code"`
}

// BMIGuidelines slice
type BMIGuidelines []BMIGuideline

// Process function
func (b *BMIGuidelines) Process(height, weight float64) (Response, error) {
	bmi := weight / (height * height)

	code := ""
	value := fmt.Sprintf("%.2f", bmi)
	target := ""

	for _, g := range *b {
		for _, c := range *g.Conditions {
			from := 0.0
			to := math.MaxFloat64
			if c.From != nil {
				from = *c.From
			}
			if c.To != nil {
				to = *c.To
			}
			if bmi >= from && bmi <= to {
				code = *g.Code
				target = *c.Target
				break
			}
		}
		if len(code) > 0 {
			break
		}
	}

	return GetResponse("bmi", code, value, target)
}

// WaistCircCondition object
type WaistCircCondition struct {
	Gender *string  `json:"gender"`
	From   *float64 `json:"from"`
	To     *float64 `json:"to"`
	Unit   *string  `json:"unit"`
	Target *string  `json:"target"`
}

// WaistCircConditions slice
type WaistCircConditions []WaistCircCondition

// WaistCircGuideline object
type WaistCircGuideline struct {
	Category   *string              `json:"category"`
	Definition *string              `json:"definition"`
	Conditions *WaistCircConditions `json:"conditions"`
	Code       *string              `json:"code"`
}

// WaistCircGuidelines slice
type WaistCircGuidelines []WaistCircGuideline

// Process function
func (b *WaistCircGuidelines) Process(gender string, waist float64, waistUnit string) (Response, error) {
	waist = tools.CalculateLength(waist, waistUnit, "cm")

	code := ""
	value := fmt.Sprintf("%.1fcm", waist)
	target := ""

	gender = tools.GetFullGenderText(gender)

	for _, g := range *b {
		for _, c := range *g.Conditions {
			from := 0.0
			to := math.MaxFloat64
			if c.From != nil {
				from = tools.CalculateLength(*c.From, *c.Unit, "cm")
			}
			if c.To != nil {
				to = tools.CalculateLength(*c.To, *c.Unit, "cm")
			}

			conditionGender := true
			if c.Gender != nil && *c.Gender != gender {
				conditionGender = false
			}

			if conditionGender && waist >= from && waist <= to {
				code = *g.Code
				target = *c.Target
				break
			}
		}
		if len(code) > 0 {
			break
		}
	}

	return GetResponse("waist circumference", code, value, target)
}

// WHRCondition object
type WHRCondition struct {
	Gender *string  `json:"gender"`
	From   *float64 `json:"from"`
	To     *float64 `json:"to"`
	Target *string  `json:"target"`
}

// WHRConditions slice
type WHRConditions []WHRCondition

// WHRGuideline object
type WHRGuideline struct {
	Category   *string        `json:"category"`
	Definition *string        `json:"definition"`
	Conditions *WHRConditions `json:"conditions"`
	Code       *string        `json:"code"`
}

// WHRGuidelines slice
type WHRGuidelines []WHRGuideline

// Process function
func (b *WHRGuidelines) Process(gender string, waist, hip float64) (Response, error) {
	whr := waist / hip

	code := ""
	value := fmt.Sprintf("%.2f", whr)
	target := ""

	gender = tools.GetFullGenderText(gender)

	for _, g := range *b {
		for _, c := range *g.Conditions {
			from := 0.0
			to := math.MaxFloat64
			if c.From != nil {
				from = *c.From
			}
			if c.To != nil {
				to = *c.To
			}

			conditionGender := true
			if c.Gender != nil && *c.Gender != gender {
				conditionGender = false
			}

			if conditionGender && whr >= from && whr <= to {
				code = *g.Code
				target = *c.Target
				break
			}
		}
		if len(code) > 0 {
			break
		}
	}

	return GetResponse("whr", code, value, target)
}

// BodyFatCondition object
type BodyFatCondition struct {
	Gender *string     `json:"gender"`
	Age    *RangeFloat `json:"age"`
	Range  *RangeFloat `json:"range"`
	Target *string     `json:"target"`
}

// BodyFatConditions slice
type BodyFatConditions []BodyFatCondition

// BodyFatGuideline object
type BodyFatGuideline struct {
	Category   *string            `json:"category"`
	Definition *string            `json:"definition"`
	Conditions *BodyFatConditions `json:"conditions"`
	Code       *string            `json:"code"`
}

// BodyFatGuidelines slice
type BodyFatGuidelines []BodyFatGuideline

// Process function
func (b *BodyFatGuidelines) Process(gender string, age, bodyFat float64) (Response, error) {
	code := ""
	value := fmt.Sprintf("%.1f%%", bodyFat)
	target := ""

	gender = tools.GetFullGenderText(gender)

	for _, g := range *b {
		for _, c := range *g.Conditions {
			ageFrom := 0.0
			ageTo := math.MaxFloat64
			bodyFatFrom := 0.0
			bodyFatTo := math.MaxFloat64

			if c.Age.From != nil {
				ageFrom = *c.Age.From
			}
			if c.Age.To != nil {
				ageTo = *c.Age.To
			}

			if c.Range.From != nil {
				bodyFatFrom = *c.Range.From
			}
			if c.Range.To != nil {
				bodyFatTo = *c.Range.To
			}

			conditionGender := true
			if c.Gender != nil && *c.Gender != gender {
				conditionGender = false
			}

			if conditionGender && (age >= ageFrom && age <= ageTo) && (bodyFat >= bodyFatFrom && bodyFat <= bodyFatTo) {
				code = *g.Code
				target = *c.Target
				break
			}
		}
		if len(code) > 0 {
			break
		}
	}

	return GetResponse("body fat", code, value, target)
}
