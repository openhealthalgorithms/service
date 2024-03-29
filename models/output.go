package models

import (
	"github.com/google/uuid"

	"github.com/openhealthalgorithms/service/pkg"
)

// OutputError object
type OutputError struct {
	Category string   `json:"category"`
	Key      string   `json:"key"`
	Messages []string `json:"messages"`
}

// Output object
type Output struct {
	Errors      []OutputError          `json:"errors"`
	Meta        Meta                   `json:"meta"`
	Assessments *ORRAssessments        `json:"assessments"`
	Goals       []ORRGoal              `json:"goals"`
	Referrals   *ORRReferrals          `json:"referrals"`
	CarePlan    *CarePlanOutput        `json:"careplan,omitempty"`
	Debug       map[string]interface{} `json:"debug,omitempty"`
}

// NewOutput function
func NewOutput(algorithmName string) *Output {
	output := &Output{}

	output.Meta.AlgorithmName = algorithmName
	output.Meta.APIVersion = pkg.GetVersion()
	output.Meta.RequestID = uuid.New()
	output.Meta.Comments = []string{}
	output.Errors = []OutputError{}

	return output
}

// Meta object
type Meta struct {
	AlgorithmName    string    `json:"algorithm"`
	RequestID        uuid.UUID `json:"request_id"`
	APIVersion       string    `json:"api_version"`
	Debug            bool      `json:"debug,omitempty"`
	CarePlan         bool      `json:"careplan,omitempty"`
	RiskModel        string    `json:"risk_model,omitempty"`
	RiskModelVersion string    `json:"risk_model_version,omitempty"`
	LabBased         bool      `json:"lab_based"`
	Comments         []string  `json:"comments"`
}
