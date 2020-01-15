package models

import (
    "github.com/google/uuid"

    "github.com/openhealthalgorithms/service/pkg"
)

// Output object
type Output struct {
    Errors      []string               `json:"errors"`
    Meta        Meta                   `json:"meta"`
    Assessments *ORRAssessments        `json:"assessments"`
    Goals       []ORRGoal              `json:"goals"`
    Referrals   *ORRReferrals          `json:"referrals"`
    CarePlan    *CarePlanOutput        `json:"careplan,omitempty"`
    Debug       map[string]interface{} `json:"input,omitempty"`
}

// NewOutput function
func NewOutput(algorithmName string) *Output {
    output := &Output{}

    output.Meta.AlgorithmName = algorithmName
    output.Meta.APIVersion = pkg.GetVersion()
    output.Meta.RequestID = uuid.New()

    return output
}

// Meta object
type Meta struct {
    AlgorithmName string    `json:"algorithm"`
    RequestID     uuid.UUID `json:"request_id"`
    APIVersion    string    `json:"api_version"`
}
