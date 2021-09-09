package algorithms

import (
	"errors"
	"strconv"
	"strings"

	a "github.com/openhealthalgorithms/service/engines/assessments"
	m "github.com/openhealthalgorithms/service/models"
	"github.com/openhealthalgorithms/service/tools"
)

// Hearts object
type Hearts struct {
	Guideline        a.Guidelines
	GuidelineContent a.GuideContents
	Goal             a.GoalGuidelines
	GoalContent      a.GoalGuideContents
}

// Process function
func (h *Hearts) Process(o m.OHARequest, colorChartPath, countriesPath string) (*m.ORRAssessments, []m.ORRGoal, *m.ORRReferrals, map[string]interface{}, []string, error) {
	var err error
	assessments := m.NewORRAssessments()
	goals := make([]m.ORRGoal, 0)
	referrals := m.NewORRReferrals()
	debug := make(map[string]interface{})
	errs := make([]string, 0)

	debugInputValue := false
	if o.Config.Debug != nil {
		debugInputValue = *o.Config.Debug
	}

	referral := false
	referralUrgent := false
	referralReasons := make([]m.ORRReferralReason, 0)

	gend := *o.Params.Demographics.Gender
	gender := strings.ToLower(gend[:1])
	age := tools.CalculateAge(float64(*o.Params.Demographics.Age.Value), *o.Params.Demographics.Age.Unit)

	fv := 0
	fvt := false
	pa := 0
	pat := false
	cSm := false
	for _, ls := range o.Params.Components.Lifestyle {
		switch *ls.Name {
		case "smoking":
			// Smoking
			exSm := false
			q := false
			if *ls.Value == "smoker" {
				cSm = true
			} else if *ls.Value == "ex-smoker" {
				exSm = true
			}
			if ls.QuitWithinYear != nil && *ls.QuitWithinYear {
				q = true
			}
			sm, err := h.Guideline.Body.Lifestyle.Smoking.Process(cSm, exSm, q)
			if err != nil {
				errs = append(errs, err.Error())
			} else {
				res := GetResults(sm, *h.GuidelineContent.Body.Contents)
				assessments.Lifestyle.Components.Smoking = &res
				if res.Refer != nil && *res.Refer != "no" {
					referral = referral || true
					ref := m.ORRReferralReason{}
					if *res.Refer == "urgent" {
						referralUrgent = referralUrgent || true
						val := true
						ref.Urgent = &val
					} else {
						val := false
						ref.Urgent = &val
					}
					ref.Type = ls.Name
					referralReasons = append(referralReasons, ref)
				}
			}
		case "alcohol_history":
			// Alcohol
			alc, err := h.Guideline.Body.Lifestyle.Alcohol.Process(tools.CalculateAlcoholConsumption((*ls.Value).(float64), *ls.Frequency), gender)
			if err != nil {
				errs = append(errs, err.Error())
			} else {
				res := GetResults(alc, *h.GuidelineContent.Body.Contents)
				assessments.Lifestyle.Components.Alcohol = &res
				if res.Refer != nil && *res.Refer != "no" {
					referral = referral || true
					ref := m.ORRReferralReason{}
					if *res.Refer == "urgent" {
						referralUrgent = referralUrgent || true
						val := true
						ref.Urgent = &val
					} else {
						val := false
						ref.Urgent = &val
					}
					ref.Type = ls.Name
					referralReasons = append(referralReasons, ref)
				}
			}
		case "exercise":
			// Physical Activity
			pat = true
			pa += tools.CalculateExercise(int((*ls.Value).(float64)), *ls.Units, *ls.Frequency, *ls.Intensity)
		case "fruit":
			// Fruits (Diet)
			fvt = true
			fv += tools.CalculateDietConsumption(int((*ls.Value).(float64)), *ls.Frequency)
			frt, err := h.Guideline.Body.Lifestyle.Diet.Fruit.Process(tools.CalculateDietConsumption(int((*ls.Value).(float64)), *ls.Frequency))
			if err != nil {
				errs = append(errs, err.Error())
			} else {
				res := GetResults(frt, *h.GuidelineContent.Body.Contents)
				assessments.Lifestyle.Components.Diet.Components.Fruit = &res
				if res.Refer != nil && *res.Refer != "no" {
					referral = referral || true
					ref := m.ORRReferralReason{}
					if *res.Refer == "urgent" {
						referralUrgent = referralUrgent || true
						val := true
						ref.Urgent = &val
					} else {
						val := false
						ref.Urgent = &val
					}
					ref.Type = ls.Name
					referralReasons = append(referralReasons, ref)
				}
			}
		case "vegetables":
			// Vegetables (Diet)
			fvt = true
			fv += tools.CalculateDietConsumption(int((*ls.Value).(float64)), *ls.Frequency)
			veg, err := h.Guideline.Body.Lifestyle.Diet.Vegetables.Process(tools.CalculateDietConsumption(int((*ls.Value).(float64)), *ls.Frequency))
			if err != nil {
				errs = append(errs, err.Error())
			} else {
				res := GetResults(veg, *h.GuidelineContent.Body.Contents)
				assessments.Lifestyle.Components.Diet.Components.Vegetable = &res
				if res.Refer != nil && *res.Refer != "no" {
					referral = referral || true
					ref := m.ORRReferralReason{}
					if *res.Refer == "urgent" {
						referralUrgent = referralUrgent || true
						val := true
						ref.Urgent = &val
					} else {
						val := false
						ref.Urgent = &val
					}
					ref.Type = ls.Name
					referralReasons = append(referralReasons, ref)
				}
			}
		}
	}
	if fvt {
		// Fruit_Vegetables (Diet)
		fveg, err := h.Guideline.Body.Lifestyle.Diet.FruitVegetables.Process(fv)
		if err != nil {
			errs = append(errs, err.Error())
		} else {
			res := GetResults(fveg, *h.GuidelineContent.Body.Contents)
			assessments.Lifestyle.Components.Diet.Components.FruitVegetable = &res
			if res.Refer != nil && *res.Refer != "no" {
				referral = referral || true
				ref := m.ORRReferralReason{}
				if *res.Refer == "urgent" {
					referralUrgent = referralUrgent || true
					val := true
					ref.Urgent = &val
				} else {
					val := false
					ref.Urgent = &val
				}
				fvtype := "fruit_vegetable"
				ref.Type = &fvtype
				referralReasons = append(referralReasons, ref)
			}
		}
	}
	if pat {
		ph, err := h.Guideline.Body.Lifestyle.PhysicalActivity.Process(pa, gender, age)
		if err != nil {
			errs = append(errs, err.Error())
		} else {
			res := GetResults(ph, *h.GuidelineContent.Body.Contents)
			assessments.Lifestyle.Components.PhysicalActivity = &res
			if res.Refer != nil && *res.Refer != "no" {
				referral = referral || true
				ref := m.ORRReferralReason{}
				if *res.Refer == "urgent" {
					referralUrgent = referralUrgent || true
					val := true
					ref.Urgent = &val
				} else {
					val := false
					ref.Urgent = &val
				}
				patype := "physical_exercise"
				ref.Type = &patype
				referralReasons = append(referralReasons, ref)
			}
		}
	}

	height, weight, hip, waist, bft := 0.0, 0.0, 0.0, 0.0, 0.0
	sbpTotal, dbpTotal, bpCount := 0, 0, 0

	for _, bm := range o.Params.Components.BodyMeasurements {
		switch *bm.Name {
		case "height":
			height = tools.ConvertLength((*bm.Value).(float64), *bm.Units)
		case "weight":
			weight = tools.ConvertWeight((*bm.Value).(float64), *bm.Units)
		case "hip":
			hip = tools.ConvertLength((*bm.Value).(float64), *bm.Units)
		case "waist":
			waist = tools.ConvertLength((*bm.Value).(float64), *bm.Units)
		case "body-fat":
			bft = tools.ConvertLength((*bm.Value).(float64), *bm.Units)
		case "blood_pressure":
			bp := (*bm.Value).(string)
			bps := strings.Split(bp, "/")
			sbp, _ := strconv.Atoi(bps[0])
			dbp, _ := strconv.Atoi(bps[1])
			sbpTotal += sbp
			dbpTotal += dbp
			bpCount++
		}
	}

	// BMI
	bmi, err := h.Guideline.Body.BodyComposition.BMI.Process(height, weight)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res := GetResults(bmi, *h.GuidelineContent.Body.Contents)
		assessments.BodyComposition.Components.BMI = &res
		if res.Refer != nil && *res.Refer != "no" {
			referral = referral || true
			ref := m.ORRReferralReason{}
			if *res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				val := true
				ref.Urgent = &val
			} else {
				val := false
				ref.Urgent = &val
			}
			patype := "bmi"
			ref.Type = &patype
			referralReasons = append(referralReasons, ref)
		}
	}

	// Waist Circumference
	waistCirc, err := h.Guideline.Body.BodyComposition.WaistCirc.Process(gender, waist, "m")
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res := GetResults(waistCirc, *h.GuidelineContent.Body.Contents)
		assessments.BodyComposition.Components.WaistCirc = &res
		if res.Refer != nil && *res.Refer != "no" {
			referral = referral || true
			ref := m.ORRReferralReason{}
			if *res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				val := true
				ref.Urgent = &val
			} else {
				val := false
				ref.Urgent = &val
			}
			patype := "waist circumference"
			ref.Type = &patype
			referralReasons = append(referralReasons, ref)
		}
	}

	// WHR
	whr, err := h.Guideline.Body.BodyComposition.WHR.Process(gender, waist, hip)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res := GetResults(whr, *h.GuidelineContent.Body.Contents)
		assessments.BodyComposition.Components.WHR = &res
		if res.Refer != nil && *res.Refer != "no" {
			referral = referral || true
			ref := m.ORRReferralReason{}
			if *res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				val := true
				ref.Urgent = &val
			} else {
				val := false
				ref.Urgent = &val
			}
			patype := "whr"
			ref.Type = &patype
			referralReasons = append(referralReasons, ref)
		}

	}

	// BodyFat
	bodyFat, err := h.Guideline.Body.BodyComposition.BodyFat.Process(gender, age, bft)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res := GetResults(bodyFat, *h.GuidelineContent.Body.Contents)
		assessments.BodyComposition.Components.BodyFat = &res
		if res.Refer != nil && *res.Refer != "no" {
			referral = referral || true
			ref := m.ORRReferralReason{}
			if *res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				val := true
				ref.Urgent = &val
			} else {
				val := false
				ref.Urgent = &val
			}
			patype := "body_fat"
			ref.Type = &patype
			referralReasons = append(referralReasons, ref)
		}
	}

	sbp, dbp := 0, 0
	sbp = sbpTotal / bpCount
	dbp = dbpTotal / bpCount

	bloodTestType, bslUnit, cholUnit := "", "", ""
	bslValue, cholValue := 0.0, 0.0

	for _, bs := range o.Params.Components.BiologicalSamples {
		switch *bs.Name {
		case "blood_sugar":
			if bs.Type != nil {
				bloodTestType = *bs.Type
			}
			bslUnit = *bs.Units
			bslValue = (*bs.Value).(float64)
		case "a1c":
			bslUnit = "%"
			bslValue = (*bs.Value).(float64)
			bloodTestType = "HbA1c"
		case "total_cholesterol":
			cholUnit = *bs.Units
			cholValue = (*bs.Value).(float64)
		}
	}

	conditions := make(map[string]bool)
	allergies := make(map[string]string)
	for _, cnd := range o.Params.Components.MedicalHistory {
		switch *cnd.Category {
		case "condition":
			conditions[strings.ToLower(*cnd.Name)] = *cnd.IsActive
		case "allergy":
			allergies[strings.ToLower(*cnd.Type)] = *cnd.Allergen
		}
	}

	medications := make(map[string]bool)
	for _, cnd := range o.Params.Components.Medications {
		medications[*cnd.Class] = true
		medications[*cnd.Category] = true
	}

	diab := false
	if d, ok := conditions["diabetes"]; ok {
		diab = d
	}

	// Diabetes
	diabetes, err := h.Guideline.Body.Diabetes.Process(diab, bslValue, bloodTestType, bslUnit, medications)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res := GetResults(diabetes, *h.GuidelineContent.Body.Contents)
		assessments.Diabetes = &res
		if res.Refer != nil && *res.Refer != "no" {
			referral = referral || true
			ref := m.ORRReferralReason{}
			if *res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				val := true
				ref.Urgent = &val
			} else {
				val := false
				ref.Urgent = &val
			}
			patype := "diabetes"
			ref.Type = &patype
			referralReasons = append(referralReasons, ref)
		}
	}

	// Blood Pressure
	diab = true
	if diabetes.Code == "DM-NONE" || diabetes.Code == "DM-PRE-DIABETES" {
		diab = false
	}
	bp, err := h.Guideline.Body.BloodPressure.Process(diab, sbp, dbp, age, medications)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res := GetResults(bp, *h.GuidelineContent.Body.Contents)
		assessments.BloodPressure = &res
		if res.Refer != nil && *res.Refer != "no" {
			referral = referral || true
			ref := m.ORRReferralReason{}
			if *res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				val := true
				ref.Urgent = &val
			} else {
				val := false
				ref.Urgent = &val
			}
			patype := "blood_pressure"
			ref.Type = &patype
			referralReasons = append(referralReasons, ref)
		}
	}

	countries := tools.Countries(countriesPath)
	region := ""
	if code, ok := countries.Countries[*o.Params.Demographics.BirthCountryCode]; ok {
		if code.Region != "#N/A" {
			region = code.Region
		} else {
			errr := errors.New("unsupported country/region")
			errs = append(errs, errr.Error())

			return assessments, goals, referrals, debug, errs, err
		}
	} else {
		errr := errors.New("invalid country/region")
		errs = append(errs, errr.Error())

		return assessments, goals, referrals, debug, errs, err
	}

	// CVD
	bmiValue, err := strconv.ParseFloat(bmi.Value, 64)
	if err != nil {
		errr := errors.New("invalid BMI value")
		errs = append(errs, errr.Error())
		return assessments, goals, referrals, debug, errs, err
	}

	cvdScore := ""
	cvd, dbg, err := h.Guideline.Body.CVD.Guidelines.Process(
		*o.Config.RiskModelVersion,
		conditions,
		age,
		*h.Guideline.Body.CVD.PreProcessing,
		medications,
		region,
		gender,
		sbp,
		cholValue,
		cholUnit,
		diab,
		cSm,
		debugInputValue,
		colorChartPath,
		*o.Config.LabBased,
		bmiValue,
	)
	if err == nil {
		cvdScore = cvd.Value
		res := GetResultWithVersion(cvd, *h.GuidelineContent.Body.Contents, *o.Config.RiskModelVersion)
		assessments.CVD = &res
		if res.Refer != nil && *res.Refer != "no" {
			referral = referral || true
			ref := m.ORRReferralReason{}
			if *res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				val := true
				ref.Urgent = &val
			} else {
				val := false
				ref.Urgent = &val
			}
			patype := "cvd"
			ref.Type = &patype
			referralReasons = append(referralReasons, ref)
		}
	} else {
		errs = append(errs, err.Error())
	}

	// Cholesterol
	if len(cvdScore) > 0 {
		cvdForChol := 1.0
		if cvdScore == "10-20%" {
			cvdForChol = 20.0
		} else if cvdScore == "20-30%" {
			cvdForChol = 30.0
		} else if cvdScore == "30-40%" {
			cvdForChol = 40.0
		} else if cvdScore == ">40%" {
			cvdForChol = 50.0
		} else if cvdScore == "<10%" {
			cvdForChol = 10.0
		}

		chol, err := h.Guideline.Body.Cholesterol.TotalCholesterol.Process(cvdForChol, age, cholValue, cholUnit, "total cholesterol", medications)
		if err != nil {
			errs = append(errs, err.Error())
		} else {
			res := GetResults(chol, *h.GuidelineContent.Body.Contents)
			assessments.Cholesterol.Components.TChol = &res
			if res.Refer != nil && *res.Refer != "no" {
				referral = referral || true
				ref := m.ORRReferralReason{}
				if *res.Refer == "urgent" {
					referralUrgent = referralUrgent || true
					val := true
					ref.Urgent = &val
				} else {
					val := false
					ref.Urgent = &val
				}
				patype := "total_cholesterol"
				ref.Type = &patype
				referralReasons = append(referralReasons, ref)
			}
		}
	}

	referrals.Refer = &referral
	referrals.Urgent = &referralUrgent
	referrals.Reasons = referralReasons

	/***** Utilizing Message Pool *****/
	dietCodes := make([]string, 0)
	if assessments.Lifestyle != nil && assessments.Lifestyle.Components.Diet != nil {
		ab := assessments.Lifestyle.Components.Diet.Components
		if ab.Fruit != nil {
			dietCodes = append(dietCodes, *ab.Fruit.Code)
		}
		if ab.Vegetable != nil {
			dietCodes = append(dietCodes, *ab.Vegetable.Code)
		}
		if ab.FruitVegetable != nil {
			dietCodes = append(dietCodes, *ab.FruitVegetable.Code)
		}
	}

	dietMessageFromPool := ""
	if len(dietCodes) > 0 {
		dietMessageFromPool = h.GuidelineContent.Body.MessagePool.Process(dietCodes, "diet")
	}
	assessments.Lifestyle.Components.Diet.Message = &dietMessageFromPool

	bodyCompositionCodes := make([]string, 0)
	if assessments.BodyComposition != nil {
		ab := assessments.BodyComposition.Components
		if ab.BMI != nil {
			bodyCompositionCodes = append(bodyCompositionCodes, *ab.BMI.Code)
		}
		if ab.WaistCirc != nil {
			bodyCompositionCodes = append(bodyCompositionCodes, *ab.WaistCirc.Code)
		}
		if ab.WHR != nil {
			bodyCompositionCodes = append(bodyCompositionCodes, *ab.WHR.Code)
		}
		if ab.BodyFat != nil {
			bodyCompositionCodes = append(bodyCompositionCodes, *ab.BodyFat.Code)
		}
	}

	bcMessageFromPool := ""
	if len(bodyCompositionCodes) > 0 {
		bcMessageFromPool = h.GuidelineContent.Body.MessagePool.Process(bodyCompositionCodes, "body-composition")
	}
	assessments.BodyComposition.Message = &bcMessageFromPool

	/***** GOALS *****/
	lSmoking, lAlcohol, lPhysicalActivity, lFruit, lVegetable, lBMI, lWaistCirc, lWHR, lBodyFat, lBloodPressure, lDiabetes, lTChol, lCVD := "", "", "", "", "", "", "", "", "", "", "", "", ""
	if assessments.Lifestyle != nil {
		if assessments.Lifestyle.Components.Smoking != nil {
			if assessments.Lifestyle.Components.Smoking.Code != nil {
				lSmoking = *assessments.Lifestyle.Components.Smoking.Code
			}
		}
		if assessments.Lifestyle.Components.Alcohol != nil {
			if assessments.Lifestyle.Components.Alcohol.Code != nil {
				lAlcohol = *assessments.Lifestyle.Components.Alcohol.Code
			}
		}
		if assessments.Lifestyle.Components.PhysicalActivity != nil {
			if assessments.Lifestyle.Components.PhysicalActivity.Code != nil {
				lPhysicalActivity = *assessments.Lifestyle.Components.PhysicalActivity.Code
			}
		}
		if assessments.Lifestyle.Components.Diet != nil {
			if assessments.Lifestyle.Components.Diet.Components.Fruit != nil {
				if assessments.Lifestyle.Components.Diet.Components.Fruit.Code != nil {
					lFruit = *assessments.Lifestyle.Components.Diet.Components.Fruit.Code
				}
			}
			if assessments.Lifestyle.Components.Diet.Components.Vegetable != nil {
				if assessments.Lifestyle.Components.Diet.Components.Vegetable.Code != nil {
					lVegetable = *assessments.Lifestyle.Components.Diet.Components.Vegetable.Code
				}
			}
		}
	}
	if assessments.BodyComposition != nil {
		if assessments.BodyComposition.Components.BMI != nil {
			if assessments.BodyComposition.Components.BMI.Code != nil {
				lBMI = *assessments.BodyComposition.Components.BMI.Code
			}
		}
		if assessments.BodyComposition.Components.WaistCirc != nil {
			if assessments.BodyComposition.Components.WaistCirc.Code != nil {
				lWaistCirc = *assessments.BodyComposition.Components.WaistCirc.Code
			}
		}
		if assessments.BodyComposition.Components.WHR != nil {
			if assessments.BodyComposition.Components.WHR.Code != nil {
				lWHR = *assessments.BodyComposition.Components.WHR.Code
			}
		}
		if assessments.BodyComposition.Components.BodyFat != nil {
			if assessments.BodyComposition.Components.BodyFat.Code != nil {
				lBodyFat = *assessments.BodyComposition.Components.BodyFat.Code
			}
		}
	}
	if assessments.BloodPressure != nil {
		if assessments.BloodPressure.Code != nil {
			lBloodPressure = *assessments.BloodPressure.Code
		}
	}
	if assessments.Diabetes != nil {
		if assessments.Diabetes.Code != nil {
			lDiabetes = *assessments.Diabetes.Code
		}
	}
	if assessments.Cholesterol != nil && assessments.Cholesterol.Components.TChol != nil {
		if assessments.Cholesterol.Components.TChol.Code != nil {
			lTChol = *assessments.Cholesterol.Components.TChol.Code
		}
	}
	if assessments.CVD != nil {
		if assessments.CVD.Code != nil {
			lCVD = *assessments.CVD.Code
		}
	}

	codes := h.Goal.GenerateGoals(lSmoking, lAlcohol, lPhysicalActivity, lFruit, lVegetable, lBMI, lWaistCirc, lWHR, lBodyFat, lBloodPressure, lDiabetes, lTChol, lCVD)
	goals = h.GoalContent.GenerateGoalsGuideline(codes...)

	if debugInputValue {
		for k, v := range dbg {
			debug[k] = v
		}
		debug["input"] = o
	}

	return assessments, goals, referrals, debug, errs, err
}

// GetResults from response
func GetResults(response a.Response, contents a.Contents) m.ORRAssessment {
	assessment := m.ORRAssessment{}

	assessment.Code = &response.Code
	assessment.Value = &response.Value
	assessment.Target = &response.Target

	if output, ok := contents[response.Code]; ok {
		assessment.Eval = output.Eval
		assessment.TFL = output.TFL
		assessment.Message = output.Message
		assessment.Refer = output.Refer
	}

	return assessment
}

// GetResultWithVersion from response
func GetResultWithVersion(response a.Response, contents a.Contents, version string) m.ORRAssessment {
	assessment := m.ORRAssessment{}

	code := response.Code + "-" + strings.ToUpper(version)

	if output, ok := contents[code]; ok {
		assessment.Code = &response.Code
		assessment.Value = &response.Value
		assessment.Target = &response.Target

		assessment.Eval = output.Eval
		assessment.TFL = output.TFL
		assessment.Message = output.Message
		assessment.Refer = output.Refer
	}

	return assessment
}
