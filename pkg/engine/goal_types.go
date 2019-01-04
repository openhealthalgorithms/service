package engine

import ds "github.com/openhealthalgorithms/service/pkg/datastructure"

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
	WeightLoss []GoalCondition `json:"weight_loss"`
	// Lifestyle       *LifestyleGuideline        `json:"lifestyle"`
	// BodyComposition *BodyCompositionGuideline  `json:"body_composition"`
	// BloodPressure   *BloodPressureGuidelines   `json:"blood_pressure"`
	// Diabetes        *DiabetesGuidelines        `json:"diabetes"`
	// Cholesterol     *CholesterolGuidelinesFull `json:"cholesterol"`
	// CVD             *CVDGuidelinesFull         `json:"cvd"`
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
	WaistCirc        []string `json:"waist_circ"`
	Whr              []string `json:"whr"`
	BodyFat          []string `json:"body_fat"`
	BloodPressure    []string `json:"bloodPressure"`
	Diabetes         []string `json:"diabetes"`
	Cholesterol      []string `json:"cholesterol"`
	Cvd              []string `json:"cvd"`
}

// GenerateGoals function
func (g *GoalGuidelines) GenerateGoals(
	smoking ds.Assessment,
	alcohol ds.Assessment,
	physicalActivity ds.Assessment,
	fruit ds.Assessment,
	vegetables ds.Assessment,
	bmi ds.Assessment,
	waistCirc ds.Assessment,
	whr ds.Assessment,
	bodyFat ds.Assessment,
	bloodPressure ds.Assessment,
	diabetes ds.Assessment,
	cholesterol ds.Assessment,
	cvd ds.Assessment,
) {

}
