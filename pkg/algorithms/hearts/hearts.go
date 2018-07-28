package hearts

import (
	"bufio"
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/structs"
	"github.com/openhealthalgorithms/service/pkg/tools"
	"github.com/openhealthalgorithms/service/pkg/types"
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
	output := calculate(inputs)

	d.Hearts = NewHearts(inputs, output)

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

func getInputs(input string) map[string]string {
	out := make(map[string]string)

	tmp := strings.Split(input, ",")
	for _, t := range tmp {
		v := strings.Split(t, ":")
		out[v[0]] = v[1]
	}

	return out
}

func calculate(params map[string]string) map[string]string {
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

	filename := fmt.Sprintf("%s_%s_%s_%s_%d.txt", "c", diabetes, gender, smoker, age)
	location, _ := fileLocation(filename, region)

	contents := getFileContents(location)

	riskScore := contents[sbp][cholValue]

	calculatedValues := make(map[string]string)
	calculatedValues["risk"] = riskScore
	calculatedValues["risk_range"] = cvdRiskValue(riskScore)

	return calculatedValues
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

func fileLocation(filename, region string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	fLocation := filepath.Join(pwd, "pkg", "algorithms", "hearts", "color_charts", region, filename)
	if _, err := os.Stat(fLocation); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
	}

	return fLocation, nil
}

func getFileContents(filename string) map[int][]string {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	lines := make(map[int][]string)
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) > 0 {
			elements := strings.Split(line, ",")
			lines[i] = elements
			i++
		}
	}

	return lines
}

func cvdRiskValue(risk string) string {
	cvdRisk, _ := strconv.Atoi(risk)

	riskValue := ">40"

	switch cvdRisk {
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
