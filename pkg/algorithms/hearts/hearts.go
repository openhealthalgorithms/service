package hearts

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"

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

	guide, guideContent, goal, goalContent, err := parseGuidesFiles(ctx)
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

	fmt.Println("---- SMOKING ----")
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

	fmt.Println("---- ALCOHOL ----")
	// Alcohol
	alc, err := engineGuide.Body.Lifestyle.Alcohol.Process(p.Alcohol, p.Gender)
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

	fmt.Println("---- PA ----")
	// Physical Activity
	ph, err := engineGuide.Body.Lifestyle.PhysicalActivity.Process(p.PhysicalActivity, p.Gender, p.Age)
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

	fmt.Println("---- DIET (FRUIT) ----")
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

	fmt.Println("---- DIET (VEGETABLE) ----")
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

	fmt.Println("---- DIET (COMBINED) ----")
	// Fruit_Vegetables (Diet)
	fveg, err := engineGuide.Body.Lifestyle.Diet.FruitVegetables.Process(p.Fruits + p.Vegetables)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		res, lifestyleActions = GetResults(fveg, *engineContent.Body.Contents, lifestyleActions)
		assessment.AssessmentsAttributes.Lifestyle.Components.Diet.Components.FruitVegetable = res
		dietGrading += res.Grading
		if res.Refer != "no" {
			referral = referral || true
			ref := ds.ReferralsResponse{}
			if res.Refer == "urgent" {
				referralUrgent = referralUrgent || true
				ref.RUrgent = true
			}
			ref.RType = "fruit_vegetable"
			referralReasons = append(referralReasons, ref)
		}
	}

	fmt.Println("---- BMI ----")
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

	fmt.Println("---- WAIST ----")
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

	fmt.Println("---- WHR ----")
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

	fmt.Println("---- BODY FAT ----")
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

	fmt.Println("---- DIABETES ----")
	// Diabetes
	diabetes, err := engineGuide.Body.Diabetes.Process(p.Diabetes, bslOrA1c, bslOrA1cType, bslOrA1cUnit, p.Medications)
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

	fmt.Println("---- BP ----")
	// Blood Pressure
	diab := false
	if diabetes.Value == "diabetes" {
		diab = true
	}
	bp, err := engineGuide.Body.BloodPressure.Process(diab, p.Sbp, p.Dbp, p.Age, p.Medications)
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

	fmt.Println("---- CVD ----")
	// CVD
	cvdScore := ""
	cvd, err := engineGuide.Body.CVD.Guidelines.Process(ctx, p.ConditionNames, p.Age, *engineGuide.Body.CVD.PreProcessing, p.Medications)
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

	fmt.Println("---- CHOLESTEROL ----")
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
		chol, err := engineGuide.Body.Cholesterol.TotalCholesterol.Process(cvdForChol, p.Age, p.TChol, p.CholUnit, "total cholesterol", p.Medications)
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

	fmt.Println("---- ASSESSMENTS MESSAGES ----")

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

	fmt.Println("---- GOALS ----")
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

	fmt.Println("---- DEBUG ----")
	if p.Debug {
		m := make(map[string]interface{})
		err := json.Unmarshal(p.Input, &m)
		if err != nil {
			assessment.Input = map[string]interface{}{"error": "Cannot preview inputs"}
		} else {
			assessment.Input = m
		}
	}

	d.Algorithm = assessment
	d.Errors = errs
	fmt.Println("---- COMPLETE ----")

	return nil
}

func parseGuidesFiles(ctx context.Context) ([]byte, []byte, []byte, []byte, error) {
	v := ctx.Value(types.KeyValuesCtx).(*types.ValuesCtx)

	var raw interface{}
	var ok bool
	var guide, guideContent, goal, goalContent []byte
	var err error

	raw, ok = v.Params.Get("cloud")
	if !ok {
		return nil, nil, nil, nil, nil
	}

	cloudEnable, ok := raw.(string)
	if !ok {
		return nil, nil, nil, nil, nil
	}

	if cloudEnable == "yes" {
		raw, ok = v.Params.Get("project")
		if !ok {
			return nil, nil, nil, nil, nil
		}

		projectName, ok := raw.(string)
		if !ok {
			return nil, nil, nil, nil, nil
		}
		if len(projectName) == 0 {
			projectName = "default-json"
		}
		fmt.Println("PROJECT: " + projectName)

		raw, ok = v.Params.Get("bucket")
		if !ok {
			return nil, nil, nil, nil, nil
		}

		bucketName, ok := raw.(string)
		if !ok {
			return nil, nil, nil, nil, nil
		}

		raw, ok = v.Params.Get("configfile")
		if !ok {
			return nil, nil, nil, nil, nil
		}

		configFile, ok := raw.(string)
		if !ok {
			return nil, nil, nil, nil, nil
		}

		err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", configFile)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		ctxBack := context.Background()

		// Creates a client.
		client, err := storage.NewClient(ctxBack)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		bucket := client.Bucket(bucketName)
		objs := bucket.Objects(ctx, &storage.Query{
			Prefix:    projectName + "/",
			Delimiter: "/",
		})
		i := 0
		for {
			attrs, err := objs.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, nil, nil, nil, err
			}
			// fmt.Println(attrs.Name)
			name := strings.ToLower(attrs.Name)
			if strings.Contains(name, "guideline_hearts.json") {
				i++
				guide, err = readStorageContent(ctxBack, bucket, attrs.Name)
				if err != nil {
					return nil, nil, nil, nil, err
				}
			} else if strings.Contains(name, "guideline_hearts_content.json") {
				i++
				guideContent, err = readStorageContent(ctxBack, bucket, attrs.Name)
				if err != nil {
					return nil, nil, nil, nil, err
				}
			} else if strings.Contains(name, "goals_hearts.json") {
				i++
				goal, err = readStorageContent(ctxBack, bucket, attrs.Name)
				if err != nil {
					return nil, nil, nil, nil, err
				}
			} else if strings.Contains(name, "goals_hearts_content.json") {
				i++
				goalContent, err = readStorageContent(ctxBack, bucket, attrs.Name)
				if err != nil {
					return nil, nil, nil, nil, err
				}
			}
		}

		if !(i == 4) {
			return nil, nil, nil, nil, errors.New("guideline files for the project are missing")
		}
	} else {
		raw, ok = v.Params.Get("guide")
		if !ok {
			return nil, nil, nil, nil, nil
		}

		guideFile, ok := raw.(string)
		if !ok {
			return nil, nil, nil, nil, nil
		}

		raw, ok = v.Params.Get("guidecontent")
		if !ok {
			return nil, nil, nil, nil, nil
		}

		guideContentFile, ok := raw.(string)
		if !ok {
			return nil, nil, nil, nil, nil
		}

		raw, ok = v.Params.Get("goal")
		if !ok {
			return nil, nil, nil, nil, nil
		}

		goalFile, ok := raw.(string)
		if !ok {
			return nil, nil, nil, nil, nil
		}

		raw, ok = v.Params.Get("goalcontent")
		if !ok {
			return nil, nil, nil, nil, nil
		}

		goalContentFile, ok := raw.(string)
		if !ok {
			return nil, nil, nil, nil, nil
		}

		guide, err = ioutil.ReadFile(guideFile)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		guideContent, err = ioutil.ReadFile(guideContentFile)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		goal, err = ioutil.ReadFile(goalFile)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		goalContent, err = ioutil.ReadFile(goalContentFile)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}

	return guide, guideContent, goal, goalContent, nil
}

func readStorageContent(ctx context.Context, bucket *storage.BucketHandle, name string) ([]byte, error) {
	rc, err := bucket.Object(name).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
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
