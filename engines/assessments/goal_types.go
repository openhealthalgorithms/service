package assessments

import (
	"github.com/openhealthalgorithms/service/tools"
)

// GoalGuidelines object
type GoalGuidelines struct {
	Meta *GoalGuidelinesMetaStructure `json:"meta"`
	Body *GoalGuidelinesBodyStructure `json:"body"`
}

// GoalGuidelinesMetaStructure object
type GoalGuidelinesMetaStructure struct {
	GoalGuidelineName *string `json:"goal_name"`
	Publisher         *string `json:"publisher"`
	PublicationDate   *string `json:"publication_date"`
}

// GoalGuidelinesBodyStructure object
type GoalGuidelinesBodyStructure struct {
	WeightControl        []GoalCondition `json:"weight_control"`
	SmokingControl       []GoalCondition `json:"smoking_control"`
	AlcoholControl       []GoalCondition `json:"alcohol_control"`
	MedicalControl       []GoalCondition `json:"medical_control"`
	BloodPressureControl []GoalCondition `json:"blood_pressure_control"`
	BloodSugarControl    []GoalCondition `json:"blood_sugar_control"`
	CholesterolControl   []GoalCondition `json:"cholesterol_control"`
}

// GoalCondition object
type GoalCondition struct {
	Category   *string    `json:"category"`
	Definition *string    `json:"definition"`
	Conditions Conditions `json:"conditions"`
	Code       *string    `json:"code"`
}

// Conditions slice
type Conditions []Condition

// Condition object
type Condition struct {
	Smoking          []string `json:"smoking"`
	Alcohol          []string `json:"alcohol"`
	PhysicalActivity []string `json:"physical_activity"`
	Fruit            []string `json:"fruit"`
	Vegetables       []string `json:"vegetables"`
	Bmi              []string `json:"bmi"`
	WaistCirc        []string `json:"wst"`
	Whr              []string `json:"whr"`
	BodyFat          []string `json:"body_fat"`
	BloodPressure    []string `json:"blood_pressure"`
	Diabetes         []string `json:"diabetes"`
	Cholesterol      []string `json:"cholesterol"`
	Cvd              []string `json:"cvd"`
}

// GenerateGoals function
func (g *GoalGuidelines) GenerateGoals(smoking, alcohol, physicalActivity, fruit, vegetables, bmi, waistCirc, whr, bodyFat, bloodPressure, diabetes, cholesterol, cvd string) []string {
	codes := make([]string, 0)

	codes = checkConditions(g.Body.WeightControl, smoking, alcohol, physicalActivity, fruit, vegetables, bmi, waistCirc, whr, bodyFat, bloodPressure, diabetes, cholesterol, cvd, codes)
	codes = checkConditions(g.Body.SmokingControl, smoking, alcohol, physicalActivity, fruit, vegetables, bmi, waistCirc, whr, bodyFat, bloodPressure, diabetes, cholesterol, cvd, codes)
	codes = checkConditions(g.Body.AlcoholControl, smoking, alcohol, physicalActivity, fruit, vegetables, bmi, waistCirc, whr, bodyFat, bloodPressure, diabetes, cholesterol, cvd, codes)
	codes = checkConditions(g.Body.MedicalControl, smoking, alcohol, physicalActivity, fruit, vegetables, bmi, waistCirc, whr, bodyFat, bloodPressure, diabetes, cholesterol, cvd, codes)
	codes = checkConditions(g.Body.BloodPressureControl, smoking, alcohol, physicalActivity, fruit, vegetables, bmi, waistCirc, whr, bodyFat, bloodPressure, diabetes, cholesterol, cvd, codes)
	codes = checkConditions(g.Body.BloodSugarControl, smoking, alcohol, physicalActivity, fruit, vegetables, bmi, waistCirc, whr, bodyFat, bloodPressure, diabetes, cholesterol, cvd, codes)
	codes = checkConditions(g.Body.CholesterolControl, smoking, alcohol, physicalActivity, fruit, vegetables, bmi, waistCirc, whr, bodyFat, bloodPressure, diabetes, cholesterol, cvd, codes)

	return codes
}

func checkConditions(goalConditions []GoalCondition, smoking, alcohol, physicalActivity, fruit, vegetables, bmi, waistCirc, whr, bodyFat, bloodPressure, diabetes, cholesterol, cvd string, codes []string) []string {
	code := ""
	for _, wc := range goalConditions {
		for _, c := range wc.Conditions {
			match := true

			match = checkCode(c.Smoking, smoking, match)
			match = checkCode(c.Alcohol, alcohol, match)
			match = checkCode(c.PhysicalActivity, physicalActivity, match)
			match = checkCode(c.Fruit, fruit, match)
			match = checkCode(c.Vegetables, vegetables, match)
			match = checkCode(c.Bmi, bmi, match)
			match = checkCode(c.WaistCirc, waistCirc, match)
			match = checkCode(c.Whr, whr, match)
			match = checkCode(c.BodyFat, bodyFat, match)
			match = checkCode(c.BloodPressure, bloodPressure, match)
			match = checkCode(c.Diabetes, diabetes, match)
			match = checkCode(c.Cholesterol, cholesterol, match)
			match = checkCode(c.Cvd, cvd, match)

			if match {
				code = *wc.Code
				break
			}
		}
		if len(code) > 0 {
			codes = append(codes, code)
			break
		}
	}

	return codes
}

func checkCode(conditions []string, code string, current bool) bool {
	if len(conditions) > 0 {
		_, found := tools.SliceContainsString(conditions, code)
		current = current && found
	}

	return current
}
