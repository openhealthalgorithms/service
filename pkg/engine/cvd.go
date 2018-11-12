package engine

import (
	"context"
	"fmt"
	"math"
	"strconv"

	who "github.com/openhealthalgorithms/service/pkg/riskmodels/whocvd"
)

// CVDGuidelinesFull object
type CVDGuidelinesFull struct {
	PreProcessing *PreProcessing `json:"pre_processing"`
	Guidelines    *CVDGuidelines `json:"guidelines"`
}

// PreProcessing object
type PreProcessing struct {
	ExistingCVD       *PreProcessingGuidelines `json:"existing_cvd"`
	HighRiskCondition *PreProcessingGuidelines `json:"high_risk_conditions"`
	AgeCheckForCVD    *PreProcessingGuidelines `json:"age_check_for_cvd"`
}

// PreProcessingCondition object
type PreProcessingCondition struct {
	AMI   *bool       `json:"ami"`
	HxCVD *bool       `json:"hx_cvd"`
	HxPVD *bool       `json:"hx_pvd"`
	HxCKD *bool       `json:"hx_ckd"`
	Age   *RangeFloat `json:"age"`
}

// PreProcessingConditions slice
type PreProcessingConditions []PreProcessingCondition

// PreProcessingGuideline object
type PreProcessingGuideline struct {
	Conditions *PreProcessingConditions `json:"conditions"`
	Return     *bool                    `json:"return"`
}

// PreProcessingGuidelines slice
type PreProcessingGuidelines []PreProcessingGuideline

// PreProcess function
func (p *PreProcessingGuidelines) PreProcess(ami, hxCVD, hxPVD, hxCKD bool, age float64) bool {
	returnValue := false

	for _, g := range *p {
		for _, c := range *g.Conditions {
			ageFrom := 0.0
			ageTo := math.MaxFloat64

			if c.Age != nil {
				if c.Age.From != nil {
					ageFrom = *c.Age.From
				}
				if c.Age.To != nil {
					ageTo = *c.Age.To
				}
			}

			conditionAMI := true
			if c.AMI != nil && *c.AMI != ami {
				conditionAMI = false
			}

			conditionCVD := true
			if c.HxCVD != nil && *c.HxCVD != hxCVD {
				conditionCVD = false
			}

			conditionPVD := true
			if c.HxPVD != nil && *c.HxPVD != hxPVD {
				conditionPVD = false
			}

			conditionCKD := true
			if c.HxCKD != nil && *c.HxCKD != hxCKD {
				conditionCKD = false
			}

			if conditionAMI && conditionCVD && conditionPVD && conditionCKD && (age >= ageFrom && age <= ageTo) {
				returnValue = *g.Return
				break
			}
		}
		if returnValue {
			break
		}
	}

	return returnValue
}

// CVDCondition object
type CVDCondition struct {
	ExistingCVD       *bool       `json:"existing_cvd"`
	HighRiskCondition *bool       `json:"high_risk_conditions"`
	AgeCheckForCVD    *bool       `json:"age_check_for_cvd"`
	Range             *RangeFloat `json:"range"`
	Target            *string     `json:"target"`
}

// CVDConditions slice
type CVDConditions []CVDCondition

// CVDGuideline object
type CVDGuideline struct {
	Category   *string        `json:"category"`
	Definition *string        `json:"definition"`
	Conditions *CVDConditions `json:"conditions"`
	Code       *string        `json:"code"`
}

// CVDGuidelines slice
type CVDGuidelines []CVDGuideline

// Process function
func (b *CVDGuidelines) Process(ctx context.Context, ami, hxCVD, hxPVD, hxCKD bool, age float64, preProcessing PreProcessing) (Response, error) {
	code := ""
	value := ""
	target := ""

	// res2B, _ := json.Marshal(preProcessing)
	// fmt.Println(string(res2B))

	ageCheckForCVD := preProcessing.AgeCheckForCVD.PreProcess(ami, hxCVD, hxPVD, hxCKD, age)
	existingCVD := preProcessing.ExistingCVD.PreProcess(ami, hxCVD, hxPVD, hxCKD, age)
	highRiskCondition := preProcessing.HighRiskCondition.PreProcess(ami, hxCVD, hxPVD, hxCKD, age)

	// CVD Assessments
	whocvd := who.New()
	err := whocvd.Get(ctx)
	if err != nil {
		return Response{}, err
	}

	riskScore, err := strconv.ParseFloat(whocvd.WHOCVD.Output["risk"], 64)
	if err != nil {
		return Response{}, err
	}

	// fmt.Printf("%+v\n", ageCheckForCVD)
	// fmt.Printf("%+v\n", existingCVD)
	// fmt.Printf("%+v\n", highRiskCondition)
	// fmt.Printf("%+v\n", riskScore)

	for _, g := range *b {
		for _, c := range *g.Conditions {
			rangeFrom := 0.0
			rangeTo := math.MaxFloat64

			if c.Range != nil {
				if c.Range.From != nil {
					rangeFrom = *c.Range.From
				}
				if c.Range.To != nil {
					rangeTo = *c.Range.To
				}
			}

			conditionAge := true
			if c.AgeCheckForCVD != nil && *c.AgeCheckForCVD != ageCheckForCVD {
				conditionAge = false
			}

			conditionExistingCVD := true
			if c.ExistingCVD != nil && *c.ExistingCVD != existingCVD {
				conditionExistingCVD = false
			}

			conditionHighRisk := true
			if c.HighRiskCondition != nil && *c.HighRiskCondition != highRiskCondition {
				conditionHighRisk = false
			}

			// res2B, _ := json.Marshal(c)
			// fmt.Println(string(res2B))
			// fmt.Printf("%+v\n", conditionAge)
			// fmt.Printf("%+v\n", conditionExistingCVD)
			// fmt.Printf("%+v\n", conditionHighRisk)

			if conditionAge && conditionExistingCVD && conditionHighRisk && (riskScore >= rangeFrom && riskScore <= rangeTo) {
				code = *g.Code
				value = fmt.Sprintf("%d", int(riskScore))
				target = *c.Target
				break
			}
		}
		if len(code) > 0 {
			break
		}
	}

	return GetResponse("cvd", code, value, target)
}
