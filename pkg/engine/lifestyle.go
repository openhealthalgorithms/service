package engine

import (
	"fmt"
	"math"

	"github.com/openhealthalgorithms/service/pkg/tools"
)

// LifestyleGuideline object
type LifestyleGuideline struct {
	Smoking          *SmokingGuidelines          `json:"smoking"`
	Alcohol          *AlcoholGuidelines          `json:"alcohol"`
	PhysicalActivity *PhysicalActivityGuidelines `json:"physical_activity"`
	Diet             *DietGuideline              `json:"diet"`
}

// SmokingCondition object
type SmokingCondition struct {
	Smoker             *bool   `json:"smoker"`
	ExSmoker           *bool   `json:"ex_smoker"`
	QuitWithin12Months *bool   `json:"quit_within_12_months"`
	Target             *string `json:"target"`
}

// SmokingConditions slice
type SmokingConditions []SmokingCondition

// SmokingGuideline object
type SmokingGuideline struct {
	Category   *string            `json:"category"`
	Definition *string            `json:"definition"`
	Conditions *SmokingConditions `json:"conditions"`
	Code       *string            `json:"code"`
}

// SmokingGuidelines slice
type SmokingGuidelines []SmokingGuideline

// Process function
func (s *SmokingGuidelines) Process(smoker bool, exSmoker bool, quitWithin12Months bool) (Response, error) {
	code := ""
	value := ""
	target := ""

	for _, g := range *s {
		for _, c := range *g.Conditions {
			conditionSmoker := true
			if c.Smoker != nil && *c.Smoker != smoker {
				conditionSmoker = false
			}

			conditionExSmoker := true
			if c.ExSmoker != nil && *c.ExSmoker != exSmoker {
				conditionExSmoker = false
			}

			conditionQuitWithin12Months := true
			if c.QuitWithin12Months != nil && *c.QuitWithin12Months != quitWithin12Months {
				conditionQuitWithin12Months = false
			}

			if conditionSmoker && conditionExSmoker && conditionQuitWithin12Months {
				code = *g.Code
				value = *g.Category
				target = *c.Target
				break
			}
		}
		if len(code) > 0 {
			break
		}
	}

	return GetResponse("smoking", code, value, target)
}

// AlcoholCondition object
type AlcoholCondition struct {
	From     *float64 `json:"from"`
	To       *float64 `json:"to"`
	Unit     *string  `json:"unit"`
	Duration *string  `json:"duration"`
	Target   *string  `json:"target"`
}

// AlcoholConditions slice
type AlcoholConditions []AlcoholCondition

// AlcoholGuide Object
type AlcoholGuide struct {
	Category   *string            `json:"category"`
	Definition *string            `json:"definition"`
	Conditions *AlcoholConditions `json:"conditions"`
	Code       *string            `json:"code"`
}

// AlcoholGuidelines slice
type AlcoholGuidelines []AlcoholGuide

// Process function
func (a *AlcoholGuidelines) Process(units float64) (Response, error) {
	code := ""
	value := fmt.Sprintf("%.1f units", units)
	target := ""

	for _, g := range *a {
		for _, c := range *g.Conditions {
			from := 0.0
			to := math.MaxFloat64
			if c.From != nil {
				from = tools.CalculateAlcoholConsumption(*c.From, *c.Duration)
			}
			if c.To != nil {
				to = tools.CalculateAlcoholConsumption(*c.To, *c.Duration)
			}
			if units >= from && units <= to {
				code = *g.Code
				target = *c.Target
				break
			}
		}
		if len(code) > 0 {
			break
		}
	}

	return GetResponse("alcohol", code, value, target)
}

// PhysicalActivityCondition object
type PhysicalActivityCondition struct {
	From   *int    `json:"from"`
	To     *int    `json:"to"`
	Unit   *string `json:"unit"`
	Type   *string `json:"type"`
	Target *string `json:"target"`
}

// PhysicalActivityConditions slice
type PhysicalActivityConditions []PhysicalActivityCondition

// PhysicalActivityGuideline object
type PhysicalActivityGuideline struct {
	Category   *string                     `json:"category"`
	Definition *string                     `json:"definition"`
	Conditions *PhysicalActivityConditions `json:"conditions"`
	Code       *string                     `json:"code"`
}

// PhysicalActivityGuidelines slice
type PhysicalActivityGuidelines []PhysicalActivityGuideline

// Process function
func (p *PhysicalActivityGuidelines) Process(duration int) (Response, error) {
	code := ""
	value := fmt.Sprintf("%d minutes", duration)
	target := ""

	for _, g := range *p {
		for _, c := range *g.Conditions {
			from := 0
			to := math.MaxInt32
			if c.From != nil {
				from = tools.CalculateExercise(*c.From, *c.Unit, "weekly", *c.Type)
			}
			if c.To != nil {
				to = tools.CalculateExercise(*c.To, *c.Unit, "weekly", *c.Type)
			}
			if duration >= from && duration <= to {
				code = *g.Code
				target = *c.Target
				break
			}
		}
		if len(code) > 0 {
			break
		}
	}

	return GetResponse("physical activity", code, value, target)
}

// DietCondition object
type DietCondition struct {
	From     *int    `json:"from"`
	To       *int    `json:"to"`
	Unit     *string `json:"unit"`
	Duration *string `json:"duration"`
	Target   *string `json:"target"`
}

// DietConditions slice
type DietConditions []DietCondition

// DietGuide Object
type DietGuide struct {
	Category   *string         `json:"category"`
	Definition *string         `json:"definition"`
	Conditions *DietConditions `json:"conditions"`
	Code       *string         `json:"code"`
}

// FruitGuidelines slice
type FruitGuidelines []DietGuide

// VegetableGuidelines slice
type VegetableGuidelines []DietGuide

// DietGuideline object
type DietGuideline struct {
	Fruit      *FruitGuidelines     `json:"fruit"`
	Vegetables *VegetableGuidelines `json:"vegetables"`
}

// Process function
func (f *FruitGuidelines) Process(servings int) (Response, error) {
	resp, err := dietCalculation(*f, servings, "fruit")

	return resp, err
}

// Process function
func (v *VegetableGuidelines) Process(servings int) (Response, error) {
	resp, err := dietCalculation(*v, servings, "vegetable")

	return resp, err
}

func dietCalculation(guides []DietGuide, servings int, dietType string) (Response, error) {
	code := ""
	value := fmt.Sprintf("%d servings", servings)
	target := ""

	for _, g := range guides {
		for _, c := range *g.Conditions {
			from := 0
			to := math.MaxInt32
			if c.From != nil {
				from = tools.CalculateDietConsumption(*c.From, *c.Duration)
			}
			if c.To != nil {
				to = tools.CalculateDietConsumption(*c.To, *c.Duration)
			}
			if servings >= from && servings <= to {
				code = *g.Code
				target = *c.Target
				break
			}
		}
		if len(code) > 0 {
			break
		}
	}

	return GetResponse(dietType, code, value, target)
}
