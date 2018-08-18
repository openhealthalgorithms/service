package whocvd

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/fatih/structs"
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

	patternsRaw, ok := v.Params.Get("params")
	if !ok {
		return nil
	}

	pts, ok := patternsRaw.(string)
	if !ok {
		return nil
	}

	inputs := getInputs(pts)

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
	Input  map[string]string
	Output map[string]string
	Debug  map[string]interface{} `structs:"debug,omitempty"`
}

// NewWHOCVD returns a Hostname from a string.
func NewWHOCVD(i, o map[string]string, d map[string]interface{}) WHOCVD {
	return WHOCVD{
		Input:  i,
		Output: o,
		Debug:  d,
	}
}

func getInputs(input string) map[string]string {
	out := make(map[string]string)

	tmp := strings.Split(input, ",")
	for _, t := range tmp {
		v := strings.Split(t, ":")
		out[v[0]] = v[1]
	}

	return out
}

func calculate(params map[string]string, debug bool) (map[string]string, map[string]interface{}) {
	diabetes := tools.TernaryString(params["diabetic"] == "true", "d", "ud")
	gender := strings.ToLower(string(params["gender"][0]))
	smoker := tools.TernaryString(params["smoker"] == "true", "s", "ns")
	age := convertAge(params["age"])
	sbp1, _ := strconv.Atoi(params["systolic1"])
	sbp2, _ := strconv.Atoi(params["systolic2"])
	sbp := convertSbp((sbp1 + sbp2) / 2)
	region := strings.ToUpper(params["region"])
	cholesterol, _ := strconv.Atoi(params["cholesterol"])
	cholesterolUnit := params["cholesterolUnit"]
	cholValue := convertCholesterol(float64(cholesterol), cholesterolUnit)
	
	contents := getContents(region, "c", diabetes, gender, smoker, age)

	riskScore := contents[sbp][cholValue]

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

func getContents(region, cholesterol, diabetes, gender, smoker string, age int) [][]int {
	values := riskConfig.RegionColorChart[strings.ToUpper(region)]
	for _, v := range values {
		if v.Cholesterol == cholesterol && v.Diabetes == diabetes && v.Gender == gender && v.Smoker == smoker && v.Age == age {
			return v.Chart
		}
	}

	return nil
}

func convertAge(a string) int {
	age, _ := strconv.Atoi(a)

	if age <= 18 {
		return 0
	} else if age < 50 {
		return 40
	} else if age < 60 {
		return 50
	} else if age < 70 {
		return 60
	}

	return 70
}

func convertSbp(sbp int) int {
	if sbp < 140 {
		return 3
	} else if sbp >= 140 && sbp < 160 {
		return 2
	} else if sbp >= 160 && sbp < 180 {
		return 1
	}

	return 0
}

func convertCholesterol(cholesterol float64, unit string) int {
	if unit == "mgdl" {
		cholesterol = cholesterol * 0.02586
	}

	tmp := int(math.Floor(cholesterol)) - 4

	if tmp < 1 {
		return 0
	} else if tmp <= 4 {
		return tmp
	}

	return 4
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
