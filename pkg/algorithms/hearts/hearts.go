package hearts

import (
	"context"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
)

// Data holds results of plugin.
type Data struct {
	Hearts `structs:"Hearts"`
}

// New returns a ready to use instance of the plugin.
func New() *Data {
	return &Data{}
}

// Get fills the Data and returns error.
func (d *Data) Get(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = errors.Errorf("%v", r)
			}
		}
	}()

	return d.get(ctx)
}

// Output returns information gathered by the plugin.
func (d *Data) Output() (map[string]interface{}, error) {
	return structs.Map(d), nil
}

// get does all the job.
func (d *Data) get(ctx context.Context) error {

	return nil
}

// Hearts represents hostname.
type Hearts struct {
	Input  map[string]string
	Output map[string]string
}

// NewHearts returns a Hostname from a string.
func NewHearts(i map[string]string, o map[string]string) Hearts {
	return Hearts{
		Input:  i,
		Output: o,
	}
}
