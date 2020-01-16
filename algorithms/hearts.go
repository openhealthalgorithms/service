package algorithms

import (
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

    referral := false
    referralUrgent := false
    referralReasons := make([]m.ORRReferralReason, 0)

    TRUE := true
    // FALSE := false
    gend := *o.Params.Demographics.Gender
    gender := strings.ToLower(gend[:1])
    age := tools.CalculateAge(float64(*o.Params.Demographics.Age.Value), *o.Params.Demographics.Age.Unit)

    fv := 0
    fvt := false
    pa := 0
    pat := false
    for _, ls := range o.Params.Components.Lifestyle {
        switch *ls.Name {
        case "smoking":
            // Smoking
            cSm := false
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
                        ref.Urgent = &TRUE
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
                    if *res.Refer == "refer" {
                        referralUrgent = referralUrgent || true
                        ref.Urgent = &TRUE
                    }
                    ref.Type = ls.Name
                    referralReasons = append(referralReasons, ref)
                }
            }
        case "exercise":
            // Physical Activity
            // TODO: MULTIPLE INPUTS
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
                        ref.Urgent = &TRUE
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
                        ref.Urgent = &TRUE
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
                    ref.Urgent = &TRUE
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
                    ref.Urgent = &TRUE
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
                ref.Urgent = &TRUE
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
                ref.Urgent = &TRUE
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
                ref.Urgent = &TRUE
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
                ref.Urgent = &TRUE
            }
            patype := "body_fat"
            ref.Type = &patype
            referralReasons = append(referralReasons, ref)
        }

    }

    referrals.Refer = &referral
    referrals.Urgent = &referralUrgent
    referrals.Reasons = referralReasons

    return assessments, goals, referrals, debug, errs, err
}

// func (d *Data) get(ctx context.Context) error {
//     var err error
//
//     bslOrA1c := 0.0
//     bslOrA1cType := "HbA1C"
//     bslOrA1cUnit := "%"
//     if p.A1C > 0.0 {
//         bslOrA1c = p.A1C
//     } else {
//         bslOrA1c = p.Bsl
//         bslOrA1cType = p.BslType
//         bslOrA1cUnit = p.BslUnit
//     }
//
//     fmt.Println("---- DIABETES ----")
//     // Diabetes
//     diabetes, err := engineGuide.Body.Diabetes.Process(p.Diabetes, bslOrA1c, bslOrA1cType, bslOrA1cUnit, p.Medications)
//     if err != nil {
//         errs = append(errs, err.Error())
//     } else {
//         res, followupActions = GetResults(diabetes, *engineContent.Body.Contents, followupActions)
//         assessment.AssessmentsAttributes.Diabetes = res
//         if res.Refer != "no" {
//             referral = referral || true
//             ref := ds.ReferralsResponse{}
//             if res.Refer == "urgent" {
//                 referralUrgent = referralUrgent || true
//                 ref.RUrgent = true
//             }
//             ref.RType = "diabetes"
//             referralReasons = append(referralReasons, ref)
//         }
//     }
//
//     fmt.Println("---- BP ----")
//     // Blood Pressure
//     diab := false
//     if diabetes.Value == "diabetes" {
//         diab = true
//     }
//     bp, err := engineGuide.Body.BloodPressure.Process(diab, p.Sbp, p.Dbp, p.Age, p.Medications)
//     if err != nil {
//         errs = append(errs, err.Error())
//     } else {
//         res, followupActions = GetResults(bp, *engineContent.Body.Contents, followupActions)
//         assessment.AssessmentsAttributes.BloodPressure = res
//         if res.Refer != "no" {
//             referral = referral || true
//             ref := ds.ReferralsResponse{}
//             if res.Refer == "urgent" {
//                 referralUrgent = referralUrgent || true
//                 ref.RUrgent = true
//             }
//             ref.RType = "blood pressure"
//             referralReasons = append(referralReasons, ref)
//         }
//     }
//
//     fmt.Println("---- CVD ----")
//     // CVD
//     cvdScore := ""
//     cvd, err := engineGuide.Body.CVD.Guidelines.Process(ctx, p.ConditionNames, p.Age, *engineGuide.Body.CVD.PreProcessing, p.Medications)
//     if err == nil {
//         cvdScore = cvd.Value
//         res, followupActions = GetResults(cvd, *engineContent.Body.Contents, followupActions)
//         assessment.AssessmentsAttributes.CVD = res
//         if res.Refer != "no" {
//             referral = referral || true
//             ref := ds.ReferralsResponse{}
//             if res.Refer == "urgent" {
//                 referralUrgent = referralUrgent || true
//                 ref.RUrgent = true
//             }
//             ref.RType = "cvd"
//             referralReasons = append(referralReasons, ref)
//         }
//     } else {
//         errs = append(errs, err.Error())
//     }
//
//     fmt.Println("---- CHOLESTEROL ----")
//     // fmt.Println("CVD Score: ", cvdScore)
//     // Cholesterol
//     if len(cvdScore) > 0 {
//         cvdForChol := 1.0
//         if cvdScore == "10-20%" {
//             cvdForChol = 20.0
//         } else if cvdScore == "20-30%" {
//             cvdForChol = 30.0
//         } else if cvdScore == "30-40%" {
//             cvdForChol = 40.0
//         } else if cvdScore == ">40%" {
//             cvdForChol = 50.0
//         } else if cvdScore == "<10%" {
//             cvdForChol = 10.0
//         }
//         // fmt.Println("CVD for Chol: ", cvdForChol)
//         chol, err := engineGuide.Body.Cholesterol.TotalCholesterol.Process(cvdForChol, p.Age, p.TChol, p.CholUnit, "total cholesterol", p.Medications)
//         if err != nil {
//             errs = append(errs, err.Error())
//         } else {
//             res, medicationsActions = GetResults(chol, *engineContent.Body.Contents, medicationsActions)
//             assessment.AssessmentsAttributes.Cholesterol.Components.TotalCholesterol = res
//             cholesterolGrading += res.Grading
//             if res.Refer != "no" {
//                 referral = referral || true
//                 ref := ds.ReferralsResponse{}
//                 if res.Refer == "urgent" {
//                     referralUrgent = referralUrgent || true
//                     ref.RUrgent = true
//                 }
//                 ref.RType = "total cholesterol"
//                 referralReasons = append(referralReasons, ref)
//             }
//         }
//     }
//
//     // assessment.RecommendationsAttributes.Lifestyle.Actions = lifestyleActions
//     // assessment.RecommendationsAttributes.Medications.Actions = medicationsActions
//     // assessment.RecommendationsAttributes.Followup.Actions = followupActions
//
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
//
//     if referral {
//         assessment.AssessmentReferralAttibutes.Refer = true
//         if referralUrgent {
//             assessment.AssessmentReferralAttibutes.Urgent = true
//         }
//         assessment.AssessmentReferralAttibutes.Reasons = referralReasons
//     }
//
//     fmt.Println("---- GOALS ----")
//     /***** GOALS *****/
//     codes := engineGoal.GenerateGoals(
//         assessment.AssessmentsAttributes.Lifestyle.Components.Smoking,
//         assessment.AssessmentsAttributes.Lifestyle.Components.Alcohol,
//         assessment.AssessmentsAttributes.Lifestyle.Components.PhysicalActivity,
//         assessment.AssessmentsAttributes.Lifestyle.Components.Diet.Components.Fruit,
//         assessment.AssessmentsAttributes.Lifestyle.Components.Diet.Components.Vegetable,
//         assessment.AssessmentsAttributes.BodyComposition.Components.BMI,
//         assessment.AssessmentsAttributes.BodyComposition.Components.WaistCirc,
//         assessment.AssessmentsAttributes.BodyComposition.Components.WHR,
//         assessment.AssessmentsAttributes.BodyComposition.Components.BodyFat,
//         assessment.AssessmentsAttributes.BloodPressure,
//         assessment.AssessmentsAttributes.Diabetes,
//         assessment.AssessmentsAttributes.Cholesterol.Components.TotalCholesterol,
//         assessment.AssessmentsAttributes.CVD,
//     )
//
//     goals := engineGoalContent.GenerateGoalsGuideline(codes...)
//     assessment.GoalsAttributes = goals
//
//     fmt.Println("---- DEBUG ----")
//     if p.Debug {
//         m := make(map[string]interface{})
//         err := json.Unmarshal(p.Input, &m)
//         if err != nil {
//             assessment.Input = map[string]interface{}{"error": "Cannot preview inputs"}
//         } else {
//             assessment.Input = m
//         }
//     }
//
//     d.Algorithm = assessment
//     d.Errors = errs
//     fmt.Println("---- COMPLETE ----")
//
//     return nil
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
