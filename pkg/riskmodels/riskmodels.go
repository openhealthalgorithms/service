// Package riskmodels contains data structures describing risk score.
package riskmodels

import (
	"context"

	"github.com/pkg/errors"

	freRm "github.com/openhealthalgorithms/service/pkg/riskmodels/framingham"
	whoCvdRm "github.com/openhealthalgorithms/service/pkg/riskmodels/whocvd"
)

var (
	// ErrRiskModelNotFound returned when algorithm not found.
	ErrRiskModelNotFound = errors.New("risk model not found")
)

// RiskModeler is an interface which any of risk model must implement.
//
// TaskExecutor relies on its methods
type RiskModeler interface {
	Get(context.Context) error
	Output() (map[string]interface{}, error)
}

// RiskModel contains risk model object.
type RiskModel struct {
	Name string
	RiskModeler
}

// DefaultRiskModels returns a set of default algorithms.
func DefaultRiskModels() map[string]*RiskModel {
	dp := make(map[string]*RiskModel)
	dp["WhoCVDRiskModel"] = &RiskModel{Name: "WhoCVDRiskModel", RiskModeler: whoCvdRm.New()}
	dp["FraminghamRiskModel"] = &RiskModel{Name: "FraminghamRiskModel", RiskModeler: freRm.New()}

	return dp
}
