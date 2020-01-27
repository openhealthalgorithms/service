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
func (h *Hearts) Process(o m.OHARequest) (*m.ORRAssessments, []m.ORRGoal, *m.ORRReferrals, map[string]interface{}, []string, error) {
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
    if _, ok := conditions["diabetes"]; !ok {
        diab = true
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
    diab = false
    if diabetes.Value == "diabetes" {
        diab = true
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

    countries := tools.Countries()
    region := ""
    if code, ok := countries[*o.Params.Demographics.BirthCountryCode]; ok {
        if code.Region != "#N/A" {
            region = code.Region
        } else {
            errr := errors.New("unsupported country/region")
            errs = append(errs, errr.Error())

            return assessments, goals, referrals, debug, errs, err
        }
    }

    // CVD
    cvdScore := ""
    cvd, dbg, err := h.Guideline.Body.CVD.Guidelines.Process(
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
    )
    if err == nil {
        cvdScore = cvd.Value
        res := GetResults(cvd, *h.GuidelineContent.Body.Contents)
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
    dietCodes := []string{
        *assessments.Lifestyle.Components.Diet.Components.Fruit.Code,
        *assessments.Lifestyle.Components.Diet.Components.Vegetable.Code,
        *assessments.Lifestyle.Components.Diet.Components.FruitVegetable.Code,
    }

    dietMessageFromPool := h.GuidelineContent.Body.MessagePool.Process(dietCodes, "diet")
    assessments.Lifestyle.Components.Diet.Message = &dietMessageFromPool

    bodyCompositionCodes := []string{
        *assessments.BodyComposition.Components.BMI.Code,
        *assessments.BodyComposition.Components.WaistCirc.Code,
        *assessments.BodyComposition.Components.WHR.Code,
        *assessments.BodyComposition.Components.BodyFat.Code,
    }

    bcMessageFromPool := h.GuidelineContent.Body.MessagePool.Process(bodyCompositionCodes, "body-composition")
    assessments.BodyComposition.Message = &bcMessageFromPool

    /***** GOALS *****/
    codes := h.Goal.GenerateGoals(
        *assessments.Lifestyle.Components.Smoking,
        *assessments.Lifestyle.Components.Alcohol,
        *assessments.Lifestyle.Components.PhysicalActivity,
        *assessments.Lifestyle.Components.Diet.Components.Fruit,
        *assessments.Lifestyle.Components.Diet.Components.Vegetable,
        *assessments.BodyComposition.Components.BMI,
        *assessments.BodyComposition.Components.WaistCirc,
        *assessments.BodyComposition.Components.WHR,
        *assessments.BodyComposition.Components.BodyFat,
        *assessments.BloodPressure,
        *assessments.Diabetes,
        *assessments.Cholesterol.Components.TChol,
        *assessments.CVD,
    )

    goals = h.GoalContent.GenerateGoalsGuideline(codes...)

    if debugInputValue {
        for k, v := range dbg {
            debug[k] = v
        }
        debug["input"] = o
    }

    return assessments, goals, referrals, debug, errs, err
}

// func (d *Data) get(ctx context.Context) error {
//     fmt.Println("---- ASSESSMENTS MESSAGES ----")
//
//     // Assessment message calculation
//     if engineContent.Body.Gradings.Lifestyle != nil {
//         for _, bc := range *engineContent.Body.Gradings.Lifestyle {
//             if lifestyleGrading >= *bc.Grading.From && lifestyleGrading <= *bc.Grading.To {
//                 assessment.AssessmentsAttributes.Lifestyle.Message = *bc.Message
//             }
//         }
//     }
//
//     if engineContent.Body.Gradings.Diet != nil {
//         for _, bc := range *engineContent.Body.Gradings.Diet {
//             if dietGrading >= *bc.Grading.From && dietGrading <= *bc.Grading.To {
//                 assessment.AssessmentsAttributes.Lifestyle.Components.Diet.Message = *bc.Message
//             }
//         }
//     }
//
//     if engineContent.Body.Gradings.BodyComposition != nil {
//         for _, bc := range *engineContent.Body.Gradings.BodyComposition {
//             if bodyCompositionGrading >= *bc.Grading.From && bodyCompositionGrading <= *bc.Grading.To {
//                 assessment.AssessmentsAttributes.BodyComposition.Message = *bc.Message
//             }
//         }
//     }
//
//     if engineContent.Body.Gradings.Cholesterol != nil {
//         for _, bc := range *engineContent.Body.Gradings.Cholesterol {
//             if cholesterolGrading >= *bc.Grading.From && cholesterolGrading <= *bc.Grading.To {
//                 assessment.AssessmentsAttributes.Cholesterol.Message = *bc.Message
//             }
//         }
//     }
// }

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
