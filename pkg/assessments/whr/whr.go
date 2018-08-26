package whr

import (
	"context"
	"fmt"

	"github.com/fatih/structs"
	"github.com/pkg/errors"

	"github.com/openhealthalgorithms/service/pkg/tools"
)

// Data holds results of plugin.
type Data struct {
	WHR `structs:"WHR_ASSESSMENT"`
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
	inputs := tools.ParseParams(ctx)
	resultCode := "WHR-0"

	whr := inputs.Waist / inputs.Hip
	currentWhr := fmt.Sprintf("%.2f", whr)
	target := "0.85"
	if inputs.Gender == "m" {
		target = "0.9"
	}

	if whr >= 0.85 && inputs.Gender == "f" {
		resultCode = "WHR-1"
	} else if whr >= 0.9 && inputs.Gender == "m" {
		resultCode = "WHR-2"
	}

	d.WHR = NewWHR(currentWhr, resultCode, target)

	return nil
}

type WHR struct {
	WHR    string `structs:"whr"`
	Code   string `structs:"code"`
	Target string `structs:"target"`
}

// NewWHR returns a BP object.
func NewWHR(whr, code, target string) WHR {
	return WHR{
		WHR:    whr,
		Code:   code,
		Target: target,
	}
}
