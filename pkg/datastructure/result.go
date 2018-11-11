package datastructure

import (
	"github.com/pborman/uuid"
)

// Result object
type Result struct {
	MetaAttributes            Meta            `json:"meta"`
	AssessmentsAttributes     Assessments     `json:"assessments"`
	GoalsAttributes           Goals           `json:"goals"`
	RecommendationsAttributes Recommendations `json:"recommendations"`
}

// Meta object
type Meta struct {
	AlgorithmName string    `json:"algorithm"`
	RequestID     uuid.UUID `json:"request_id"`
}

/* * * * * Assessments * * * * */

// Assessments object
type Assessments struct {
	Lifestyle       LifestyleAssessment       `json:"lifestyle"`
	BodyComposition BodyCompositionAssessment `json:"body_composition"`
	BloodPressure   Assessment                `json:"blood_pressure"`
	Diabetes        Assessment                `json:"diabetes"`
	Cholesterol     CholesterolAssessment     `json:"cholesterol"`
	CVD             Assessment                `json:"cvd"`
}

// LifestyleAssessment object
type LifestyleAssessment struct {
	Smoking          Assessment     `json:"smoking"`
	Alcohol          Assessment     `json:"alcohol"`
	PhysicalActivity Assessment     `json:"physical_activity"`
	Diet             DietAssessment `json:"diet"`
}

// DietAssessment object
type DietAssessment struct {
	Fruit     Assessment `json:"fruit"`
	Vegetable Assessment `json:"vegetable"`
}

// BodyCompositionAssessment object
type BodyCompositionAssessment struct {
	BMI       Assessment `json:"bmi"`
	WaistCirc Assessment `json:"waist_circ"`
	WHR       Assessment `json:"whr"`
	BodyFat   Assessment `json:"body_fat"`
}

// CholesterolAssessment object
type CholesterolAssessment struct {
	TotalCholesterol Assessment `json:"total_cholesterol"`
	HDL              Assessment `json:"hdl"`
	LDL              Assessment `json:"ldl"`
	TG               Assessment `json:"tg"`
}

// Assessment object
type Assessment struct {
	Code   string `json:"code"`
	Value  string `json:"value"`
	Target string `json:"target"`
	Output Output `json:"output"`
}

// Output object
type Output struct {
	Code  string `json:"code"`
	Type  string `json:"type"`
	Color string `json:"color"`
}

/* * * * * Goals * * * * */

// Goals object
type Goals []Goal

// Goal object
type Goal struct {
	Code    string   `json:"code"`
	Name    string   `json:"name"`
	Reasons []string `json:"reasons"`
}

/* * * * * Recommendations * * * * */

// Recommendations object
type Recommendations struct {
	Lifestyle   LifestyleRecommendation   `json:"lifestyle"`
	Medications MedicationsRecommendation `json:"medications"`
	Followup    FollowupRecommendation    `json:"followup"`
}

// LifestyleRecommendation object
type LifestyleRecommendation struct {
	Actions `json:"actions"`
}

// MedicationsRecommendation object
type MedicationsRecommendation struct {
	Actions `json:"actions"`
}

// FollowupRecommendation object
type FollowupRecommendation struct {
	Actions `json:"actions"`
}

// Actions is a collection of actions
type Actions []Action

// Action object
type Action struct {
	Goal     string   `json:"goal"`
	Messages []string `json:"messages"`
}

// NewResult returns a Result object with meta information
func NewResult(algorithmName string) Result {
	result := Result{}

	result.MetaAttributes.AlgorithmName = algorithmName
	result.MetaAttributes.RequestID = uuid.NewRandom()

	return result
}
