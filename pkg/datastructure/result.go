package datastructure

import (
	"github.com/pborman/uuid"
)

// Result object
type Result struct {
	MetaAttributes              Meta               `structs:"meta" json:"meta"`
	AssessmentsAttributes       Assessments        `structs:"assessments" json:"assessments"`
	GoalsAttributes             Goals              `structs:"goals" json:"goals"`
	RecommendationsAttributes   Recommendations    `structs:"recommendations" json:"recommendations"`
	AssessmentReferralAttibutes AssessmentReferral `structs:"referrals" json:"referrals"`
}

// Meta object
type Meta struct {
	AlgorithmName string    `structs:"algorithm" json:"algorithm"`
	RequestID     uuid.UUID `structs:"request_id" json:"request_id"`
}

/* * * * * Assessments * * * * */

// Assessments object
type Assessments struct {
	Lifestyle       LifestyleAssessment       `structs:"lifestyle" json:"lifestyle"`
	BodyComposition BodyCompositionAssessment `structs:"body_composition" json:"body_composition"`
	BloodPressure   Assessment                `structs:"blood_pressure" json:"blood_pressure"`
	Diabetes        Assessment                `structs:"diabetes" json:"diabetes"`
	Cholesterol     CholesterolAssessment     `structs:"cholesterol" json:"cholesterol"`
	CVD             Assessment                `structs:"cvd" json:"cvd"`
}

// LifestyleAssessment object
type LifestyleAssessment struct {
	Components LifestyleComponents `structs:"components" json:"components"`
	Message    string              `structs:"message" json:"message"`
}

// LifestyleComponents object
type LifestyleComponents struct {
	Smoking          Assessment     `structs:"smoking" json:"smoking"`
	Alcohol          Assessment     `structs:"alcohol" json:"alcohol"`
	PhysicalActivity Assessment     `structs:"physical_activity" json:"physical_activity"`
	Diet             DietAssessment `structs:"diet" json:"diet"`
}

// DietAssessment object
type DietAssessment struct {
	Components DietComponents `structs:"components" json:"components"`
	Message    string         `structs:"message" json:"message"`
}

// DietComponents object
type DietComponents struct {
	Fruit     Assessment `structs:"fruit" json:"fruit"`
	Vegetable Assessment `structs:"vegetable" json:"vegetable"`
}

// BodyCompositionAssessment object
type BodyCompositionAssessment struct {
	Components BodyCompositionComponents `structs:"components" json:"components"`
	Message    string                    `structs:"message" json:"message"`
}

// BodyCompositionComponents object
type BodyCompositionComponents struct {
	BMI       Assessment `structs:"bmi" json:"bmi"`
	WaistCirc Assessment `structs:"waist_circ" json:"waist_circ"`
	WHR       Assessment `structs:"whr" json:"whr"`
	BodyFat   Assessment `structs:"body_fat" json:"body_fat"`
}

// CholesterolAssessment object
type CholesterolAssessment struct {
	Components CholesterolComponents `structs:"components" json:"components"`
	Message    string                `structs:"message" json:"message"`
}

// CholesterolComponents object
type CholesterolComponents struct {
	TotalCholesterol Assessment `structs:"total_cholesterol" json:"total_cholesterol"`
	HDL              Assessment `structs:"hdl" json:"hdl"`
	LDL              Assessment `structs:"ldl" json:"ldl"`
	TG               Assessment `structs:"tg" json:"tg"`
}

// Assessment object
type Assessment struct {
	Code    string `structs:"code" json:"code"`
	Eval    string `structs:"eval" json:"eval"`
	Grading int    `structs:"grading" json:"-"`
	TFL     string `structs:"tfl" json:"tfl"`
	Value   string `structs:"value" json:"value"`
	Target  string `structs:"target" json:"target"`
	Message string `structs:"message" json:"message"`
	Refer   string `structs:"refer" json:"refer"`
}

/* * * * * Goals * * * * */

// Goals object
type Goals interface{}

// Goal object
// type Goal struct {
// 	Code    string   `structs:"code" json:"code"`
// 	Name    string   `structs:"name" json:"name"`
// 	Reasons []string `structs:"reasons" json:"reasons"`
// }

/* * * * * Recommendations * * * * */

// Recommendations object
type Recommendations struct {
	Lifestyle   LifestyleRecommendation   `structs:"lifestyle" json:"lifestyle"`
	Medications MedicationsRecommendation `structs:"medications" json:"medications"`
	Followup    FollowupRecommendation    `structs:"followup" json:"followup"`
}

// LifestyleRecommendation object
type LifestyleRecommendation struct {
	Actions `structs:"actions" json:"actions"`
}

// MedicationsRecommendation object
type MedicationsRecommendation struct {
	Actions `structs:"actions" json:"actions"`
}

// FollowupRecommendation object
type FollowupRecommendation struct {
	Actions `structs:"actions" json:"actions"`
}

// Actions is a collection of actions
type Actions []Action

// Action object
type Action struct {
	Goal     string   `structs:"goal" json:"goal"`
	Messages []string `structs:"messages" json:"messages"`
}

/* * * * * Recommendations * * * * */

// AssessmentReferral object
type AssessmentReferral struct {
	Refer   bool     `structs:"refer" json:"refer"`
	Urgent  bool     `structs:"urgent" json:"urgent"`
	Reasons []string `structs:"reasons" json:"reasons"`
}

// NewResult returns a Result object with meta information
func NewResult(algorithmName string) Result {
	result := Result{}

	result.MetaAttributes.AlgorithmName = algorithmName
	result.MetaAttributes.RequestID = uuid.NewRandom()

	return result
}
