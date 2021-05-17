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

// Calculate risk score
func Calculate(
	version,
	region,
	gender string,
	ageValue float64,
	sbpValue int,
	tChol float64,
	cholUnit string,
	diabetes,
	currentSmoker,
	debug bool,
	colorChartPath string,
	labBased bool,
	bmi float64,
) (map[string]string, map[string]interface{}, error) {
	var err error
	fmt.Println("region:", region, "gender:", gender, "ageValue:", ageValue, "sbpValue:", sbpValue, "tChol:", tChol, "cholUnit:", cholUnit, "diabetes:", diabetes, "currentSmoker:", currentSmoker, "labBased:", labBased)
	if !(len(region) > 0 && (gender == "m" || gender == "f")) {
		return nil, nil, errors.New("invalid input for risk model, valid region needed")
	}

	group, row, column, riskValues := [][]int{}, -1, -1, []config.RiskRange{}
	if version == "who_ish_2019" {
		group, row, column, riskValues, err = processV2(colorChartPath, region, labBased, gender, ageValue, sbpValue, tChol, cholUnit, diabetes, currentSmoker, bmi)
	} else {
		group, row, column, riskValues, err = processV1(colorChartPath, region, labBased, gender, ageValue, sbpValue, tChol, cholUnit, diabetes, currentSmoker)
	}

	if err != nil {
		return nil, nil, err
	}

	if len(group) == 0 {
		return nil, nil, errors.New("invalid version given")
	}

	riskScore := group[row][column]
	calculatedValues := make(map[string]string)
	calculatedValues["risk"] = strconv.Itoa(riskScore)
	calculatedValues["risk_range"] = cvdRiskValue(riskScore, riskValues)

	debugValues := make(map[string]interface{})

	if debug {
		debugValues["matrix"] = group
		debugValues["index"] = fmt.Sprintf("%d, %d", row, column)
	}

	return calculatedValues, debugValues, nil
}

func processV1(colorChartPath, region string, labBased bool, gender string, ageValue float64, sbpValue int, tChol float64, cholUnit string, diabetes, currentSmoker bool) ([][]int, int, int, []config.RiskRange, error) {
	riskConfigFile, err := ioutil.ReadFile(colorChartPath)
	if err != nil {
		return nil, -1, -1, nil, errors.New("Read file error: " + err.Error())
	}

	riskConfig := config.WHOColorChartV1{}
	if err := json.Unmarshal(riskConfigFile, &riskConfig); err != nil {
		return nil, -1, -1, nil, err
	}

	regionData := config.Regions{}
	if val, ok := riskConfig.ColorCharts[region]; ok {
		regionData = val
	} else {
		return nil, -1, -1, nil, errors.New("invalid region given")
	}

	chol := regionData.NonChol
	cholValue := 0
	if labBased && tChol > 0 && len(cholUnit) > 3 {
		chol = regionData.Chol
		// get chol index
		cholV := tools.ConvertCholesterol(tChol, cholUnit)
		cholValue = cholIndex(cholV, riskConfig.Meta.Configuration.Cholesterols)
	}

	diab := chol.NonDiabetic
	if diabetes {
		diab = chol.Diabetic
	}

	gend := diab.Female
	if gender == "m" {
		gend = diab.Male
	}

	smoking := gend.NonSmoker
	if currentSmoker {
		smoking = gend.Smoker
	}

	// get age index
	ageRange := ageKey(int(ageValue), riskConfig.Meta.Configuration.Ages)
	var ageGroup config.AgeRanges
	if val, ok := smoking[ageRange]; ok {
		ageGroup = val
	} else {
		return nil, -1, -1, nil, errors.New("cannot perform calculation for the given age")
	}

	// get sbp index
	sbp := sbpIndex(sbpValue, riskConfig.Meta.Configuration.Systolics)

	return ageGroup, sbp, cholValue, riskConfig.Meta.Configuration.RiskValues, nil
}

func processV2(colorChartPath, region string, labBased bool, gender string, ageValue float64, sbpValue int, tChol float64, cholUnit string, diabetes, currentSmoker bool, bmi float64) ([][]int, int, int, []config.RiskRange, error) {
	riskConfigFile, err := ioutil.ReadFile(colorChartPath)
	if err != nil {
		return nil, -1, -1, nil, errors.New("Read file error: " + err.Error())
	}

	riskConfig := config.WHOColorChartV2{}
	if err := json.Unmarshal(riskConfigFile, &riskConfig); err != nil {
		return nil, -1, -1, nil, err
	}

	regionData := config.RegionsV2{}
	if val, ok := riskConfig.ColorCharts[region]; ok {
		regionData = val
	} else {
		return nil, -1, -1, nil, errors.New("invalid region given")
	}

	var gend config.Gender
	columnIndex := -1

	if labBased && tChol > 0 && len(cholUnit) > 3 {
		chol := regionData.LabBased
		cholV := tools.ConvertCholesterol(tChol, cholUnit)
		columnIndex = cholIndex(cholV, riskConfig.Meta.Configuration.Cholesterols)
		diab := chol.NonDiabetic
		if diabetes {
			diab = chol.Diabetic
		}
		gend = diab.Female
		if gender == "m" {
			gend = diab.Male
		}
	} else {
		nLab := regionData.NonLabBased
		columnIndex = cholIndex(bmi, riskConfig.Meta.Configuration.BMIs)
		gend = nLab.Female
		if gender == "m" {
			gend = nLab.Male
		}
	}

	smoking := gend.NonSmoker
	if currentSmoker {
		smoking = gend.Smoker
	}

	// get age index
	ageRange := ageKey(int(ageValue), riskConfig.Meta.Configuration.Ages)
	var ageGroup config.AgeRanges
	if val, ok := smoking[ageRange]; ok {
		ageGroup = val
	} else {
		return nil, -1, -1, nil, errors.New("cannot perform calculation for the given age")
	}

	// get sbp index
	sbp := sbpIndex(sbpValue, riskConfig.Meta.Configuration.Systolics)

	return ageGroup, sbp, columnIndex, riskConfig.Meta.Configuration.RiskValues, nil
}

func cvdRiskValue(risk int, risks []config.RiskRange) string {
	for _, r := range risks {
		if risk >= r.From && risk <= r.To {
			return r.Value
		}
	}

	return ""
}

func ageKey(age int, ages []config.AgeRange) string {
	for _, r := range ages {
		if age >= r.From && age <= r.To {
			return r.Key
		}
	}

	return ""
}

func sbpIndex(sbp int, systolics []config.SBPRange) int {
	for i, r := range systolics {
		if sbp >= r.From && sbp <= r.To {
			return i
		}
	}

	return -1
}

func cholIndex(chol float64, chols []config.CholRange) int {
	for i, r := range chols {
		if chol >= r.From && chol <= r.To {
			return i
		}
	}

	return -1
}
