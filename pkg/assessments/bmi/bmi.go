package bmi

import (
	"context"
	"fmt"

	"github.com/fatih/structs"
	"github.com/pkg/errors"

	"github.com/openhealthalgorithms/service/pkg/tools"
)

// Data holds results of plugin.
type Data struct {
	BMI `structs:"BMI_ASSESSMENT"`
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
	resultCode := ""

	bmi := inputs.Weight / (inputs.Height * inputs.Height)
	currentBmi := fmt.Sprintf("%.2f", bmi)
	target := "18.5 - 24.9"

	if bmi < 18.5 {
		resultCode = "BMI-1"
	} else if bmi < 25 {
		resultCode = "BMI-0"
	} else if bmi < 30 {
		resultCode = "BMI-2"
	} else {
		resultCode = "BMI-3"
	}

	d.BMI = NewBMI(currentBmi, resultCode, target)

	return nil
}

type BMI struct {
	BMI    string `structs:"bmi"`
	Code   string `structs:"code"`
	Target string `structs:"target"`
}

// NewBMI returns a BP object.
func NewBMI(bmi, code, target string) BMI {
	return BMI{
		BMI:    bmi,
		Code:   code,
		Target: target,
	}
}
