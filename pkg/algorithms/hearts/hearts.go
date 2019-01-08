package hearts

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/fatih/structs"
	"github.com/pkg/errors"

	ds "github.com/openhealthalgorithms/service/pkg/datastructure"
	"github.com/openhealthalgorithms/service/pkg/engine"
	"github.com/openhealthalgorithms/service/pkg/tools"
	"github.com/openhealthalgorithms/service/pkg/types"
)

// Data holds results of plugin.
type Data struct {
	Algorithm ds.Result `json:"algorithm"`
	Errors    []string  `json:"errors"`
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
	var err error

	v := ctx.Value(types.KeyValuesCtx).(*types.ValuesCtx)
	raw, ok := v.Params.Get("params")
	if !ok {
		return nil
	}

	p, ok := raw.(tools.Params)
	if !ok {
		return nil
	}

	raw, ok = v.Params.Get("guide")
	if !ok {
		return nil
	}

	guideFile, ok := raw.(string)
	if !ok {
		return nil
	}

	raw, ok = v.Params.Get("guidecontent")
	if !ok {
		return nil
	}

	guideContentFile, ok := raw.(string)
	if !ok {
		return nil
	}

	raw, ok = v.Params.Get("goal")
	if !ok {
		return nil
	}

	goalFile, ok := raw.(string)
	if !ok {
		return nil
	}

	raw, ok = v.Params.Get("goalcontent")
	if !ok {
		return nil
	}

	goalContentFile, ok := raw.(string)
	if !ok {
		return nil
	}

	guide, err := ioutil.ReadFile(guideFile)
	if err != nil {
		return err
	}

	guideContent, err := ioutil.ReadFile(guideContentFile)
	if err != nil {
		return err
	}

	goal, err := ioutil.ReadFile(goalFile)
	if err != nil {
		return err
	}

	goalContent, err := ioutil.ReadFile(goalContentFile)
	if err != nil {
		return err
	}

	engineGuide := engine.Guidelines{}
	if err := json.Unmarshal(guide, &engineGuide); err != nil {
		return err
	}

	engineContent := engine.GuideContents{}
	if err := json.Unmarshal(guideContent, &engineContent); err != nil {
		return err
	}

	engineGoal := engine.GoalGuidelines{}
	if err := json.Unmarshal(goal, &engineGoal); err != nil {
		return err
	}

	engineGoalContent := engine.GoalGuideContents{}
	if err := json.Unmarshal(goalContent, &engineGoalContent); err != nil {
		return err
	}

	// engineGuide.Body.Lifestyle.Smoking
	// res2B, _ := json.Marshal(engineGoal)
	// fmt.Println(string(res2B))

	// res2C, _ := json.Marshal(engineGoalContent)
	// fmt.Println(string(res2C))
	// fmt.Printf("%+v\n", p)

	assessment := ds.NewResult("Hearts Algorithm")
	lifestyleActions := make([]ds.Action, 0)
	medicationsActions := make([]ds.Action, 0)
	followupActions := make([]ds.Action, 0)

	var res ds.Assessment
	errs := make([]string, 0)

	lifestyleGrading := 0
	bodyCompositionGrading := 0
	cholesterolGrading := 0

	referral := false
	referralUrgent := false
	referralReasons := make([]ds.ReferralsResponse, 0)

	// Smoking
	sm, err := engineGuide.Body.Lifestyle.Smoking.Process(p.Smoker.CurrentSmoker, p.Smoker.ExSmoker, p.Smoker.QuitWithinYear)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res, lifestyleActions = GetResults(sm, *engineContent.Body.Contents, lifestyleActions)
		assessment.AssessmentsAttributes.Lifestyle.Components.Smoking = res
		lifestyleGrading += res.Grading
		if res.Refer != "no" {
			referral = referral || true
			ref := ds.ReferralsResponse{}
			if res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				ref.RUrgent = true
			}
			ref.RType = "smoking"
			referralReasons = append(referralReasons, ref)
		}
	}

	// Alcohol
	alc, err := engineGuide.Body.Lifestyle.Alcohol.Process(p.Alcohol)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res, lifestyleActions = GetResults(alc, *engineContent.Body.Contents, lifestyleActions)
		assessment.AssessmentsAttributes.Lifestyle.Components.Alcohol = res
		lifestyleGrading += res.Grading
		if res.Refer != "no" {
			referral = referral || true
			ref := ds.ReferralsResponse{}
			if res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				ref.RUrgent = true
			}
			ref.RType = "alcohol"
			referralReasons = append(referralReasons, ref)
		}
	}

	// Physical Activity
	ph, err := engineGuide.Body.Lifestyle.PhysicalActivity.Process(p.PhysicalActivity)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res, lifestyleActions = GetResults(ph, *engineContent.Body.Contents, lifestyleActions)
		assessment.AssessmentsAttributes.Lifestyle.Components.PhysicalActivity = res
		lifestyleGrading += res.Grading
		if res.Refer != "no" {
			referral = referral || true
			ref := ds.ReferralsResponse{}
			if res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				ref.RUrgent = true
			}
			ref.RType = "physical activity"
			referralReasons = append(referralReasons, ref)
		}
	}

	dietGrading := 0

	// Fruits (Diet)
	frt, err := engineGuide.Body.Lifestyle.Diet.Fruit.Process(p.Fruits)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res, lifestyleActions = GetResults(frt, *engineContent.Body.Contents, lifestyleActions)
		assessment.AssessmentsAttributes.Lifestyle.Components.Diet.Components.Fruit = res
		dietGrading += res.Grading
		if res.Refer != "no" {
			referral = referral || true
			ref := ds.ReferralsResponse{}
			if res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				ref.RUrgent = true
			}
			ref.RType = "fruit"
			referralReasons = append(referralReasons, ref)
		}
	}

	// Vegetables (Diet)
	veg, err := engineGuide.Body.Lifestyle.Diet.Vegetables.Process(p.Vegetables)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res, lifestyleActions = GetResults(veg, *engineContent.Body.Contents, lifestyleActions)
		assessment.AssessmentsAttributes.Lifestyle.Components.Diet.Components.Vegetable = res
		dietGrading += res.Grading
		if res.Refer != "no" {
			referral = referral || true
			ref := ds.ReferralsResponse{}
			if res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				ref.RUrgent = true
			}
			ref.RType = "vegetable"
			referralReasons = append(referralReasons, ref)
		}
	}

	// BMI
	bmi, err := engineGuide.Body.BodyComposition.BMI.Process(p.Height, p.Weight)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res, lifestyleActions = GetResults(bmi, *engineContent.Body.Contents, lifestyleActions)
		assessment.AssessmentsAttributes.BodyComposition.Components.BMI = res
		bodyCompositionGrading += res.Grading
		if res.Refer != "no" {
			referral = referral || true
			ref := ds.ReferralsResponse{}
			if res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				ref.RUrgent = true
			}
			ref.RType = "bmi"
			referralReasons = append(referralReasons, ref)
		}
	}

	// Waist Circumference
	waistCirc, err := engineGuide.Body.BodyComposition.WaistCirc.Process(p.Gender, p.Waist, p.WaistUnit)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res, lifestyleActions = GetResults(waistCirc, *engineContent.Body.Contents, lifestyleActions)
		assessment.AssessmentsAttributes.BodyComposition.Components.WaistCirc = res
		bodyCompositionGrading += res.Grading
		if res.Refer != "no" {
			referral = referral || true
			ref := ds.ReferralsResponse{}
			if res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				ref.RUrgent = true
			}
			ref.RType = "waist circumference"
			referralReasons = append(referralReasons, ref)
		}
	}

	// WHR
	whr, err := engineGuide.Body.BodyComposition.WHR.Process(p.Gender, p.Waist, p.Hip)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res, lifestyleActions = GetResults(whr, *engineContent.Body.Contents, lifestyleActions)
		assessment.AssessmentsAttributes.BodyComposition.Components.WHR = res
		bodyCompositionGrading += res.Grading
		if res.Refer != "no" {
			referral = referral || true
			ref := ds.ReferralsResponse{}
			if res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				ref.RUrgent = true
			}
			ref.RType = "whr"
			referralReasons = append(referralReasons, ref)
		}
	}

	// BodyFat
	bodyFat, err := engineGuide.Body.BodyComposition.BodyFat.Process(p.Gender, p.Age, p.BodyFat)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res, lifestyleActions = GetResults(bodyFat, *engineContent.Body.Contents, lifestyleActions)
		assessment.AssessmentsAttributes.BodyComposition.Components.BodyFat = res
		bodyCompositionGrading += res.Grading
		if res.Refer != "no" {
			referral = referral || true
			ref := ds.ReferralsResponse{}
			if res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				ref.RUrgent = true
			}
			ref.RType = "body fat"
			referralReasons = append(referralReasons, ref)
		}
	}

	bslOrA1c := 0.0
	bslOrA1cType := "HbA1C"
	bslOrA1cUnit := "%"
	if p.A1C > 0.0 {
		bslOrA1c = p.A1C
	} else {
		bslOrA1c = p.Bsl
		bslOrA1cType = p.BslType
		bslOrA1cUnit = p.BslUnit
	}

	// Diabetes
	diabetes, err := engineGuide.Body.Diabetes.Process(p.Diabetes, bslOrA1c, bslOrA1cType, bslOrA1cUnit)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res, followupActions = GetResults(diabetes, *engineContent.Body.Contents, followupActions)
		assessment.AssessmentsAttributes.Diabetes = res
		if res.Refer != "no" {
			referral = referral || true
			ref := ds.ReferralsResponse{}
			if res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				ref.RUrgent = true
			}
			ref.RType = "diabetes"
			referralReasons = append(referralReasons, ref)
		}
	}

	// Blood Pressure
	diab := false
	if diabetes.Value == "diabetes" {
		diab = true
	}
	bp, err := engineGuide.Body.BloodPressure.Process(diab, p.Sbp, p.Dbp)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res, followupActions = GetResults(bp, *engineContent.Body.Contents, followupActions)
		assessment.AssessmentsAttributes.BloodPressure = res
		if res.Refer != "no" {
			referral = referral || true
			ref := ds.ReferralsResponse{}
			if res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				ref.RUrgent = true
			}
			ref.RType = "blood pressure"
			referralReasons = append(referralReasons, ref)
		}
	}

	// CVD
	cvdScore := ""
	cvd, err := engineGuide.Body.CVD.Guidelines.Process(ctx, p.AMI, p.Cvd, p.Pvd, p.Ckd, p.Age, *engineGuide.Body.CVD.PreProcessing)
	if err == nil {
		cvdScore = cvd.Value
		res, followupActions = GetResults(cvd, *engineContent.Body.Contents, followupActions)
		assessment.AssessmentsAttributes.CVD = res
		if res.Refer != "no" {
			referral = referral || true
			ref := ds.ReferralsResponse{}
			if res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				ref.RUrgent = true
			}
			ref.RType = "cvd"
			referralReasons = append(referralReasons, ref)
		}
	} else {
		errs = append(errs, err.Error())
	}

	// fmt.Println("CVD Score: ", cvdScore)
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
		// fmt.Println("CVD for Chol: ", cvdForChol)
		chol, err := engineGuide.Body.Cholesterol.TotalCholesterol.Process(cvdForChol, p.Age, p.TChol, p.CholUnit, "total cholesterol")
		if err != nil {
			errs = append(errs, err.Error())
		} else {
			res, medicationsActions = GetResults(chol, *engineContent.Body.Contents, medicationsActions)
			assessment.AssessmentsAttributes.Cholesterol.Components.TotalCholesterol = res
			cholesterolGrading += res.Grading
			if res.Refer != "no" {
				referral = referral || true
				ref := ds.ReferralsResponse{}
				if res.Refer == "urgent" {
					referralUrgent = referralUrgent || true
					ref.RUrgent = true
				}
				ref.RType = "total cholesterol"
				referralReasons = append(referralReasons, ref)
			}
		}
	}

	// assessment.RecommendationsAttributes.Lifestyle.Actions = lifestyleActions
	// assessment.RecommendationsAttributes.Medications.Actions = medicationsActions
	// assessment.RecommendationsAttributes.Followup.Actions = followupActions

	// Assessment message calculation
	if engineContent.Body.Gradings.Lifestyle != nil {
		for _, bc := range *engineContent.Body.Gradings.Lifestyle {
			if lifestyleGrading >= *bc.Grading.From && lifestyleGrading <= *bc.Grading.To {
				assessment.AssessmentsAttributes.Lifestyle.Message = *bc.Message
			}
		}
	}

	if engineContent.Body.Gradings.Diet != nil {
		for _, bc := range *engineContent.Body.Gradings.Diet {
			if dietGrading >= *bc.Grading.From && dietGrading <= *bc.Grading.To {
				assessment.AssessmentsAttributes.Lifestyle.Components.Diet.Message = *bc.Message
			}
		}
	}

	if engineContent.Body.Gradings.BodyComposition != nil {
		for _, bc := range *engineContent.Body.Gradings.BodyComposition {
			if bodyCompositionGrading >= *bc.Grading.From && bodyCompositionGrading <= *bc.Grading.To {
				assessment.AssessmentsAttributes.BodyComposition.Message = *bc.Message
			}
		}
	}

	if engineContent.Body.Gradings.Cholesterol != nil {
		for _, bc := range *engineContent.Body.Gradings.Cholesterol {
			if cholesterolGrading >= *bc.Grading.From && cholesterolGrading <= *bc.Grading.To {
				assessment.AssessmentsAttributes.Cholesterol.Message = *bc.Message
			}
		}
	}

	if referral {
		assessment.AssessmentReferralAttibutes.Refer = true
		if referralUrgent {
			assessment.AssessmentReferralAttibutes.Urgent = true
		}
		assessment.AssessmentReferralAttibutes.Reasons = referralReasons
	}

	/***** GOALS *****/
	codes := engineGoal.GenerateGoals(
		assessment.AssessmentsAttributes.Lifestyle.Components.Smoking,
		assessment.AssessmentsAttributes.Lifestyle.Components.Alcohol,
		assessment.AssessmentsAttributes.Lifestyle.Components.PhysicalActivity,
		assessment.AssessmentsAttributes.Lifestyle.Components.Diet.Components.Fruit,
		assessment.AssessmentsAttributes.Lifestyle.Components.Diet.Components.Vegetable,
		assessment.AssessmentsAttributes.BodyComposition.Components.BMI,
		assessment.AssessmentsAttributes.BodyComposition.Components.WaistCirc,
		assessment.AssessmentsAttributes.BodyComposition.Components.WHR,
		assessment.AssessmentsAttributes.BodyComposition.Components.BodyFat,
		assessment.AssessmentsAttributes.BloodPressure,
		assessment.AssessmentsAttributes.Diabetes,
		assessment.AssessmentsAttributes.Cholesterol.Components.TotalCholesterol,
		assessment.AssessmentsAttributes.CVD,
	)

	goals := engineGoalContent.GenerateGoalsGuideline(codes...)
	assessment.GoalsAttributes = goals

	d.Algorithm = assessment
	d.Errors = errs

	return nil
}

// GetResults from response
func GetResults(response engine.Response, contents engine.Contents, advices ds.Actions) (ds.Assessment, ds.Actions) {
	assessment := ds.Assessment{}

	assessment.Code = response.Code
	assessment.Value = response.Value
	assessment.Target = response.Target

	if output, ok := contents[response.Code]; ok {
		assessment.Eval = *output.Eval
		assessment.TFL = *output.TFL
		assessment.Message = *output.Message
		assessment.Refer = *output.Refer
		assessment.Grading = *output.Grading

		advice := ds.Action{}
		advice.Goal = *output.Eval
		advice.Messages = append(advice.Messages, *output.Message)
		advices = append(advices, advice)
	}

	return assessment, advices
}
