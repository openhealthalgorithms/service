package assessments

import (
	"context"

	"github.com/pkg/errors"

	bpAssess "github.com/openhealthalgorithms/service/pkg/assessments/bp"
)

var (
	// ErrAssessmentNotFound returned when algorithm not found.
	ErrAssessmentNotFound = errors.New("assessment not found")
)

// Assess is an interface which any of assessment must implement.
type Assess interface {
	Get(context.Context) error
	Output() (map[string]interface{}, error)
}

// Assessment contains risk model object.
type Assessment struct {
	Name string
	Assess
}

// DefaultRiskModels returns a set of default algorithms.
func DefaultAssessments() map[string]*Assessment {
	dp := make(map[string]*Assessment)
	dp["BPAssessment"] = &Assessment{Name: "BPAssessment", Assess: bpAssess.New()}

	return dp
}
