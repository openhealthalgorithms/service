package riskmodels

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/pkg/errors"

	"github.com/openhealthalgorithms/service/riskmodels/config"
	"github.com/openhealthalgorithms/service/tools"
)

var (
	riskConfig config.WHOColorChart
)

// Calculate risk score
func Calculate(region, gender string, ageValue float64, sbpValue int, tChol float64, cholUnit string, diabetes, currentSmoker, debug bool, colorChartPath string) (map[string]string, map[string]interface{}, error) {
	if !(len(region) > 0 && (gender == "m" || gender == "f")) {
		return nil, nil, errors.New("invalid input for risk model, valid region needed")
	}
	age := tools.ConvertAge(ageValue)
	sbp := tools.ConvertSbp(sbpValue)
	cholValue := -1
	chol := false
	if tChol > 0 && len(cholUnit) > 3 {
		chol = true
		cholValue = tools.ConvertCholesterol(tChol, cholUnit)
	}

	riskConfigFile, err := ioutil.ReadFile(colorChartPath)
	if err != nil {
		return nil, nil, errors.New("Read file error: " + err.Error())
	}

	if err := json.Unmarshal(riskConfigFile, &riskConfig); err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	contents := getContents(region, gender, chol, diabetes, currentSmoker, age)

	riskScore := contents[sbp][0]
	if chol {
		riskScore = contents[sbp][cholValue]
	}

	calculatedValues := make(map[string]string)
	calculatedValues["risk"] = strconv.Itoa(riskScore)
	calculatedValues["risk_range"] = cvdRiskValue(riskScore)

	debugValues := make(map[string]interface{})

	if debug {
		debugValues["matrix"] = contents
		debugValues["index"] = fmt.Sprintf("%d, %d", sbp, cholValue)
	}

	return calculatedValues, debugValues, nil
}

func getContents(region, gender string, cholesterol, diabetes, smoker bool, age int) [][]int {
	values := riskConfig.ColorCharts[region]
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
