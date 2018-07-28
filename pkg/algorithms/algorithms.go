// Package algorithms contains data structures describing algorithms.
package algorithms

import (
	"context"

	"github.com/pkg/errors"

	heartsAlg "github.com/openhealthalgorithms/service/pkg/algorithms/hearts"
)

var (
	// ErrAlgorithmNotFound returned when algorithm not found.
	ErrAlgorithmNotFound = errors.New("algorithm not found")
)

// Algorithmer is an interface which any of algorithms must implement.
//
// TaskExecutor relies on its methods
type Algorithmer interface {
	Get(context.Context) error
	Output() (map[string]interface{}, error)
}

// Algorithm contains algorithm object.
type Algorithm struct {
	Name string
	Algorithmer
}

// DefaultAlgorithms returns a set of default algorithms.
func DefaultAlgorithms() map[string]*Algorithm {
	dp := make(map[string]*Algorithm)
	dp["HeartsAlgorithm"] = &Algorithm{Name: "HeartsAlgorithm", Algorithmer: heartsAlg.New()}

	return dp
}
