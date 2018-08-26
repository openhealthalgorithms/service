package hearts

import (
	"context"
	"github.com/openhealthalgorithms/service/pkg/assessments/bmi"
	"github.com/openhealthalgorithms/service/pkg/assessments/bp"
	"github.com/openhealthalgorithms/service/pkg/assessments/whr"

	"github.com/fatih/structs"
	"github.com/pkg/errors"

	who "github.com/openhealthalgorithms/service/pkg/riskmodels/whocvd"
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
	var err error

	whocvd := who.New()
	err = whocvd.Get(ctx)
	if err != nil {
		return nil
	}

	bpa := bp.New()
	err = bpa.Get(ctx)
	if err != nil {
		return nil
	}

	bmis := bmi.New()
	err = bmis.Get(ctx)
	if err != nil {
		return nil
	}

	whrs := whr.New()
	err = whrs.Get(ctx)
	if err != nil {
		return nil
	}

	d.Hearts = NewHearts(whocvd.WHOCVD.Output, bpa.BP, bmis.BMI, whrs.WHR)

	return nil
}

// Hearts represents hostname.
type Hearts struct {
	CVDRisk       map[string]string `structs:"cvd_risk"`
	BPAssessment  bp.BP             `structs:"bp"`
	BMIAssessment bmi.BMI           `structs:"bmi"`
	WHRAssessment whr.WHR           `structs:"whr"`
}

// NewHearts returns a Hostname from a string.
func NewHearts(cvd map[string]string, bp bp.BP, bmi bmi.BMI, whr whr.WHR) Hearts {
	return Hearts{
		CVDRisk:       cvd,
		BPAssessment:  bp,
		BMIAssessment: bmi,
		WHRAssessment: whr,
	}
}
