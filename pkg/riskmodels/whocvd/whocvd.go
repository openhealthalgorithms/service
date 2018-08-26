package whocvd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fatih/structs"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"

	"github.com/openhealthalgorithms/service/pkg/riskmodels/common/config"
	"github.com/openhealthalgorithms/service/pkg/tools"
	"github.com/openhealthalgorithms/service/pkg/types"
)

var (
	riskConfig = config.NewSettings()
)

// Data holds results of plugin.
type Data struct {
	WHOCVD `structs:"WHO_ISH_RISK_SCORE"`
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
	v := ctx.Value(types.KeyValuesCtx).(*types.ValuesCtx)
	inputs := tools.ParseParams(ctx)

	debugOutput := false
	_, oks := v.Params.Get("debug")
	if oks {
		debugOutput = true
	}

	output, debug := calculate(inputs, debugOutput)

	d.WHOCVD = NewWHOCVD(inputs, output, debug)

	return nil
}

// WHOCVD represents hostname.
type WHOCVD struct {
	RequestId uuid.UUID
	Input     tools.Params
	Output    map[string]string
	Debug     map[string]interface{} `structs:"debug,omitempty"`
}

// NewWHOCVD returns a Hostname from a string.
func NewWHOCVD(i tools.Params, o map[string]string, d map[string]interface{}) WHOCVD {
	return WHOCVD{
		RequestId: uuid.NewRandom(),
		Input:  i,
		Output: o,
		Debug:  d,
	}
}

func calculate(params tools.Params, debug bool) (map[string]string, map[string]interface{}) {
	age := tools.ConvertAge(params.Age)
	sbp := tools.ConvertSbp(params.Sbp)
	cholValue := -1
	chol := false
	if params.TChol > 0 && len(params.CholUnit) > 3 {
		chol = true
		cholValue = tools.ConvertCholesterol(params.TChol, params.CholUnit)
	}

	contents := getContents(params.Region, params.Gender, chol, params.Diabetes, params.CurrentSmoker, age)

	riskScore := contents[sbp][0]
	if chol {
		riskScore = contents[sbp][cholValue]
	}

	calculatedValues := make(map[string]string)
	calculatedValues["risk"] = strconv.Itoa(riskScore)
	calculatedValues["risk_range"] = cvdRiskValue(riskScore)

	if debug {
		debugValues := make(map[string]interface{})
		debugValues["matrix"] = contents
		debugValues["index"] = fmt.Sprintf("%d, %d", sbp, cholValue)

		return calculatedValues, debugValues
	}

	return calculatedValues, nil
}

func getContents(region, gender string, cholesterol, diabetes, smoker bool, age int) [][]int {
	values := riskConfig.RegionColorChart[region]
	for _, v := range values {
		if v.Cholesterol == cholesterol && v.Diabetes == diabetes && v.Gender == gender && v.Smoker == smoker && v.Age == age {
			return v.Chart
		}
	}

	return nil
}

func cvdRiskValue(risk int) string {
	riskValue := ">40"

	switch risk {
	case 10:
		riskValue = "<10"
	case 20:
		riskValue = "10-20"
	case 30:
		riskValue = "20-30"
	case 40:
		riskValue = "30-40"
	}

	return riskValue
}
