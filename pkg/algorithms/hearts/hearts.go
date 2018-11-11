package hearts

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/fatih/structs"
	"github.com/pkg/errors"

	"github.com/openhealthalgorithms/service/pkg/datastructure"
	"github.com/openhealthalgorithms/service/pkg/engine"
	"github.com/openhealthalgorithms/service/pkg/tools"
	"github.com/openhealthalgorithms/service/pkg/types"
)

// Data holds results of plugin.
type Data struct {
	Algorithm datastructure.Result `json:"algorithm"`
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

	guide, err := ioutil.ReadFile(guideFile)
	if err != nil {
		return err
	}

	guideContent, err := ioutil.ReadFile(guideContentFile)
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

	// engineGuide.Body.Lifestyle.Smoking
	// res2B, _ := json.Marshal(engineContent)
	// fmt.Println(string(res2B))
	fmt.Printf("%+v\n", p)

	assessment := datastructure.NewResult("Hearts Algorithm")
	lifestyleActions := make([]datastructure.Action, 0)
	medicationsActions := make([]datastructure.Action, 0)
	followupActions := make([]datastructure.Action, 0)

	var res datastructure.Assessment

	// Smoking
	sm, err := engineGuide.Body.Lifestyle.Smoking.Process(p.Smoker.CurrentSmoker, p.Smoker.ExSmoker, p.Smoker.QuitWithinYear)
	if err != nil {
		return err
	}
	res, lifestyleActions = GetResults(sm, *engineContent.Body.Contents, lifestyleActions)
	assessment.AssessmentsAttributes.Lifestyle.Smoking = res

	// Alcohol
	alc, err := engineGuide.Body.Lifestyle.Alcohol.Process(p.Alcohol)
	if err != nil {
		return err
	}
	res, lifestyleActions = GetResults(alc, *engineContent.Body.Contents, lifestyleActions)
	assessment.AssessmentsAttributes.Lifestyle.Alcohol = res

	// Physical Activity
	ph, err := engineGuide.Body.Lifestyle.PhysicalActivity.Process(p.PhysicalActivity)
	if err != nil {
		return err
	}
	res, lifestyleActions = GetResults(ph, *engineContent.Body.Contents, lifestyleActions)
	assessment.AssessmentsAttributes.Lifestyle.PhysicalActivity = res

	// Fruits (Diet)
	frt, err := engineGuide.Body.Lifestyle.Diet.Fruit.Process(p.Fruits)
	if err != nil {
		return err
	}
	res, lifestyleActions = GetResults(frt, *engineContent.Body.Contents, lifestyleActions)
	assessment.AssessmentsAttributes.Lifestyle.Diet.Fruit = res

	// Vegetables (Diet)
	veg, err := engineGuide.Body.Lifestyle.Diet.Vegetables.Process(p.Vegetables)
	if err != nil {
		return err
	}
	res, lifestyleActions = GetResults(veg, *engineContent.Body.Contents, lifestyleActions)
	assessment.AssessmentsAttributes.Lifestyle.Diet.Vegetable = res

	// BMI
	bmi, err := engineGuide.Body.BodyComposition.BMI.Process(p.Height, p.Weight)
	if err != nil {
		return err
	}
	res, lifestyleActions = GetResults(bmi, *engineContent.Body.Contents, lifestyleActions)
	assessment.AssessmentsAttributes.BodyComposition.BMI = res

	// Waist Circumference
	waistCirc, err := engineGuide.Body.BodyComposition.WaistCirc.Process(p.Gender, p.Waist, p.WaistUnit)
	if err != nil {
		return err
	}
	res, lifestyleActions = GetResults(waistCirc, *engineContent.Body.Contents, lifestyleActions)
	assessment.AssessmentsAttributes.BodyComposition.WaistCirc = res

	// WHR
	whr, err := engineGuide.Body.BodyComposition.WHR.Process(p.Gender, p.Waist, p.Hip)
	if err != nil {
		return err
	}
	res, lifestyleActions = GetResults(whr, *engineContent.Body.Contents, lifestyleActions)
	assessment.AssessmentsAttributes.BodyComposition.WHR = res

	// BodyFat
	bodyFat, err := engineGuide.Body.BodyComposition.BodyFat.Process(p.Gender, p.Age, p.BodyFat)
	if err != nil {
		return err
	}
	res, lifestyleActions = GetResults(bodyFat, *engineContent.Body.Contents, lifestyleActions)
	assessment.AssessmentsAttributes.BodyComposition.BodyFat = res

	// Diabetes
	diabetes, err := engineGuide.Body.Diabetes.Process(p.Diabetes, p.Bsl, p.BslType, p.BslUnit)
	if err != nil {
		return err
	}
	res, followupActions = GetResults(diabetes, *engineContent.Body.Contents, followupActions)
	assessment.AssessmentsAttributes.Diabetes = res

	// Blood Pressure
	diab := false
	if diabetes.Value == "diabetes" {
		diab = true
	}
	bp, err := engineGuide.Body.BloodPressure.Process(diab, p.Sbp, p.Dbp)
	if err != nil {
		return err
	}
	res, followupActions = GetResults(bp, *engineContent.Body.Contents, followupActions)
	assessment.AssessmentsAttributes.BloodPressure = res

	// CVD
	cvdScore := -1.0
	cvd, err := engineGuide.Body.CVD.Guidelines.Process(ctx, p.AMI, p.Cvd, p.Pvd, p.Ckd, p.Age, *engineGuide.Body.CVD.PreProcessing)
	if err == nil {
		cvdScoreFromAssessment, cerr := strconv.ParseFloat(cvd.Value, 64)
		if cerr == nil {
			cvdScore = cvdScoreFromAssessment
		}
		res, followupActions = GetResults(cvd, *engineContent.Body.Contents, followupActions)
		assessment.AssessmentsAttributes.CVD = res
	} else {
		return err
	}

	// Cholesterol
	if cvdScore > 0 {
		chol, err := engineGuide.Body.Cholesterol.TotalCholesterol.Process(cvdScore, p.Age, p.TChol, p.CholUnit, "total cholesterol")
		if err != nil {
			return err
		}
		res, medicationsActions = GetResults(chol, *engineContent.Body.Contents, medicationsActions)
		assessment.AssessmentsAttributes.Cholesterol.TotalCholesterol = res
	}

	assessment.RecommendationsAttributes.Lifestyle.Actions = lifestyleActions
	assessment.RecommendationsAttributes.Medications.Actions = medicationsActions
	assessment.RecommendationsAttributes.Followup.Actions = followupActions

	d.Algorithm = assessment

	return nil
}

// GetResults from response
func GetResults(response engine.Response, contents engine.Contents, advices datastructure.Actions) (datastructure.Assessment, datastructure.Actions) {
	assessment := datastructure.Assessment{}

	assessment.Code = response.Code
	assessment.Value = response.Value
	assessment.Target = response.Target

	if output, ok := contents[response.Code]; ok {
		assessment.Output.Code = *output.Code
		assessment.Output.Type = *output.Type
		assessment.Output.Color = *output.Color
		advice := datastructure.Action{}
		advice.Goal = *output.Type
		advice.Messages = append(advice.Messages, *output.Advice)
		advices = append(advices, advice)
	}

	return assessment, advices
}
