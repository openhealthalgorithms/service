package algorithms

import (
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
        case "physical-activity":
            // Physical Activity
            // TODO: MULTIPLE INPUTS
            ph, err := h.Guideline.Body.Lifestyle.PhysicalActivity.Process(tools.CalculateExercise((*ls.Value).(int), *ls.Units, *ls.Frequency, *ls.Intensity), gender, age)
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
                    ref.Type = ls.Name
                    referralReasons = append(referralReasons, ref)
                }
            }
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
//     fmt.Println("---- PA ----")
//
//     dietGrading := 0
//
//     fmt.Println("---- DIET (FRUIT) ----")
//     // Fruits (Diet)
//     frt, err := engineGuide.Body.Lifestyle.Diet.Fruit.Process(p.Fruits)
//     if err != nil {
//         errs = append(errs, err.Error())
//     } else {
//         res, lifestyleActions = GetResults(frt, *engineContent.Body.Contents, lifestyleActions)
//         assessment.AssessmentsAttributes.Lifestyle.Components.Diet.Components.Fruit = res
//         dietGrading += res.Grading
//         if res.Refer != "no" {
//             referral = referral || true
//             ref := ds.ReferralsResponse{}
//             if res.Refer == "urgent" {
//                 referralUrgent = referralUrgent || true
//                 ref.RUrgent = true
//             }
//             ref.RType = "fruit"
//             referralReasons = append(referralReasons, ref)
//         }
//     }
//
//     fmt.Println("---- DIET (VEGETABLE) ----")
//     // Vegetables (Diet)
//     veg, err := engineGuide.Body.Lifestyle.Diet.Vegetables.Process(p.Vegetables)
//     if err != nil {
//         errs = append(errs, err.Error())
//     } else {
//         res, lifestyleActions = GetResults(veg, *engineContent.Body.Contents, lifestyleActions)
//         assessment.AssessmentsAttributes.Lifestyle.Components.Diet.Components.Vegetable = res
//         dietGrading += res.Grading
//         if res.Refer != "no" {
//             referral = referral || true
//             ref := ds.ReferralsResponse{}
//             if res.Refer == "urgent" {
//                 referralUrgent = referralUrgent || true
//                 ref.RUrgent = true
//             }
//             ref.RType = "vegetable"
//             referralReasons = append(referralReasons, ref)
//         }
//     }
//
//     fmt.Println("---- DIET (COMBINED) ----")
//     // Fruit_Vegetables (Diet)
//     fveg, err := engineGuide.Body.Lifestyle.Diet.FruitVegetables.Process(p.Fruits + p.Vegetables)
//     if err != nil {
//         errs = append(errs, err.Error())
//     } else {
//         res, lifestyleActions = GetResults(fveg, *engineContent.Body.Contents, lifestyleActions)
//         assessment.AssessmentsAttributes.Lifestyle.Components.Diet.Components.FruitVegetable = res
//         dietGrading += res.Grading
//         if res.Refer != "no" {
//             referral = referral || true
//             ref := ds.ReferralsResponse{}
//             if res.Refer == "urgent" {
//                 referralUrgent = referralUrgent || true
//                 ref.RUrgent = true
//             }
//             ref.RType = "fruit_vegetable"
//             referralReasons = append(referralReasons, ref)
//         }
//     }
//
//     fmt.Println("---- BMI ----")
//     // BMI
//     bmi, err := engineGuide.Body.BodyComposition.BMI.Process(p.Height, p.Weight)
//     if err != nil {
//         errs = append(errs, err.Error())
//     } else {
//         res, lifestyleActions = GetResults(bmi, *engineContent.Body.Contents, lifestyleActions)
//         assessment.AssessmentsAttributes.BodyComposition.Components.BMI = res
//         bodyCompositionGrading += res.Grading
//         if res.Refer != "no" {
//             referral = referral || true
//             ref := ds.ReferralsResponse{}
//             if res.Refer == "urgent" {
//                 referralUrgent = referralUrgent || true
//                 ref.RUrgent = true
//             }
//             ref.RType = "bmi"
//             referralReasons = append(referralReasons, ref)
//         }
//     }
//
//     fmt.Println("---- WAIST ----")
//     // Waist Circumference
//     waistCirc, err := engineGuide.Body.BodyComposition.WaistCirc.Process(p.Gender, p.Waist, p.WaistUnit)
//     if err != nil {
//         errs = append(errs, err.Error())
//     } else {
//         res, lifestyleActions = GetResults(waistCirc, *engineContent.Body.Contents, lifestyleActions)
//         assessment.AssessmentsAttributes.BodyComposition.Components.WaistCirc = res
//         bodyCompositionGrading += res.Grading
//         if res.Refer != "no" {
//             referral = referral || true
//             ref := ds.ReferralsResponse{}
//             if res.Refer == "urgent" {
//                 referralUrgent = referralUrgent || true
//                 ref.RUrgent = true
//             }
//             ref.RType = "waist circumference"
//             referralReasons = append(referralReasons, ref)
//         }
//     }
//
//     fmt.Println("---- WHR ----")
//     // WHR
//     whr, err := engineGuide.Body.BodyComposition.WHR.Process(p.Gender, p.Waist, p.Hip)
//     if err != nil {
//         errs = append(errs, err.Error())
//     } else {
//         res, lifestyleActions = GetResults(whr, *engineContent.Body.Contents, lifestyleActions)
//         assessment.AssessmentsAttributes.BodyComposition.Components.WHR = res
//         bodyCompositionGrading += res.Grading
//         if res.Refer != "no" {
//             referral = referral || true
//             ref := ds.ReferralsResponse{}
//             if res.Refer == "urgent" {
//                 referralUrgent = referralUrgent || true
//                 ref.RUrgent = true
//             }
//             ref.RType = "whr"
//             referralReasons = append(referralReasons, ref)
//         }
//     }
//
//     fmt.Println("---- BODY FAT ----")
//     // BodyFat
//     bodyFat, err := engineGuide.Body.BodyComposition.BodyFat.Process(p.Gender, p.Age, p.BodyFat)
//     if err != nil {
//         errs = append(errs, err.Error())
//     } else {
//         res, lifestyleActions = GetResults(bodyFat, *engineContent.Body.Contents, lifestyleActions)
//         assessment.AssessmentsAttributes.BodyComposition.Components.BodyFat = res
//         bodyCompositionGrading += res.Grading
//         if res.Refer != "no" {
//             referral = referral || true
//             ref := ds.ReferralsResponse{}
//             if res.Refer == "urgent" {
//                 referralUrgent = referralUrgent || true
//                 ref.RUrgent = true
//             }
//             ref.RType = "body fat"
//             referralReasons = append(referralReasons, ref)
//         }
//     }
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
