package hearts

import (
	"context"
	jp "github.com/buger/jsonparser"
	"github.com/fatih/structs"
	"github.com/openhealthalgorithms/service/pkg/assessments"
	who "github.com/openhealthalgorithms/service/pkg/riskmodels/whocvd"
	"github.com/openhealthalgorithms/service/pkg/tools"
	"github.com/openhealthalgorithms/service/pkg/types"
	"github.com/pkg/errors"
	"io/ioutil"
)

// Data holds results of plugin.
type Data struct {
	Hearts `structs:"Hearts"`
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
		return nil
	}

	guideContent, err := ioutil.ReadFile(guideContentFile)
	if err != nil {
		return nil
	}

	// Diabetes
	output := getOutput(guideContent, "body", "messages", "diabetes", p.DiabetesAssessment.Code)
	db := Diabetes{p.DiabetesAssessment.BSL, p.DiabetesAssessment.Code, p.DiabetesAssessment.Status, output}

	// BP
	output = getOutput(guideContent, "body", "messages", "blood_pressure", p.BPAssessment.Code)
	bp := BloodPressure{p.BPAssessment.BP, p.BPAssessment.Code, p.BPAssessment.Target, output}

	// CVD Assessments
	whocvd := who.New()
	err = whocvd.Get(ctx)
	if err != nil {
		return nil
	}

	highRisks := getOutput(guide, "body", "high_risk_conditions")
	hra, _ := assessments.GetHighRisksWithHrc(p.Sbp, p.Dbp, p.Age, p.Conditions, highRisks)

	ot, _, _, _ := jp.Get(guide, "body", "cvd_risk", whocvd.WHOCVD.Output["risk_range"])

	o := make(map[string]interface{})
	o["score"], _ = jp.GetString(ot, "score")
	o["label"], _ = jp.GetString(ot, "label")
	o["bp_target"], _ = jp.GetString(ot, "bp_target")
	o["follow_up_interval"], _ = jp.GetString(ot, "follow_up_interval")
	o["follow_up_message"], _ = jp.GetString(ot, "follow_up_message")
	adv := getOutput(ot, "advice")
	o["advice"] = adv

	advices := make(map[string]interface{})
	for _, a := range adv {
		advices[a], _ = jp.GetString(guideContent, "body", "messages", "advice", a)
	}
	o["management"] = advices

	cv := CvdAssessment{hra, whocvd.WHOCVD.Output, o}

	// Lifestyle
	output = getOutput(guideContent, "body", "messages", "anthro", p.BMIAssessment.Code)
	bmi := Bmi{p.BMIAssessment.BMI, p.BMIAssessment.Code, p.BMIAssessment.Target, output}

	output = getOutput(guideContent, "body", "messages", "anthro", p.WHRAssessment.Code)
	whr := Whr{p.WHRAssessment.WHR, p.WHRAssessment.Code, p.WHRAssessment.Target, output}

	output = getOutput(guideContent, "body", "messages", "smoking", p.SmokingAssessment.Code)
	smoking := Smoking{p.SmokingAssessment.Code, p.SmokingAssessment.Status, p.SmokingAssessment.SmokingCalc, output}

	exTarget, err := jp.GetInt(guide, "body", "targets", "general", "physical_activity", "active_time")
	ex := p.ExerciseAssessment
	if err == nil {
		ex, _ = assessments.GetExerciseWithTarget(p.PhysicalActivity, int(exTarget))
	}
	output = getOutput(guideContent, "body", "messages", "nutrition", ex.Code)
	exa := Exercise{ex.Current, ex.Code, ex.Target, output}

	dtf, err := jp.GetInt(guide, "body", "targets", "general", "diet", "fruit")
	if err != nil {
		dtf = 2
	}
	dtv, err := jp.GetInt(guide, "body", "targets", "general", "diet", "vegetables")
	if err == nil {
		dtv = 5
	}
	dt, _ := assessments.GetDietWithTarget(p.Fruits, p.Vegetables, int(dtf), int(dtv))
	output = getOutput(guideContent, "body", "messages", "physical_activity", dt.Code)
	val := Values{p.Fruits, p.Vegetables}
	diet := Diet{val, dt.Code, output}

	lf := Lifestyle{
		Bmi:      bmi,
		Whr:      whr,
		Diet:     diet,
		Exercise: exa,
		Smoking:  smoking,
	}

	d.Hearts = NewHearts(db, bp, lf, cv)

	return nil
}

func getOutput(guideContent []byte, keys ...string) []string {
	output := make([]string, 0)
	jp.ArrayEach(guideContent, func(value []byte, dataType jp.ValueType, offset int, err error) {
		output = append(output, string(value))
	}, keys...)

	return output
}

// Hearts represents hostname.
type Hearts struct {
	Diabetes      `structs:"diabetes"`
	BloodPressure `structs:"blood_pressure"`
	Lifestyle     `structs:"lifestyle"`
	CvdAssessment `structs:"cvd_assessment"`
}

type Diabetes struct {
	BSL    float64  `structs:"value"`
	Code   string   `structs:"code"`
	Status bool     `structs:"status"`
	Output []string `structs:"output"`
}

type BloodPressure struct {
	BP     string   `structs:"bp"`
	Code   string   `structs:"code"`
	Target string   `structs:"target"`
	Output []string `structs:"output"`
}

type Bmi struct {
	BMI    string   `structs:"value"`
	Code   string   `structs:"code"`
	Target string   `structs:"target"`
	Output []string `structs:"output"`
}

type Whr struct {
	WHR    string   `structs:"value"`
	Code   string   `structs:"code"`
	Target string   `structs:"target"`
	Output []string `structs:"output"`
}

type Values struct {
	Fruit      int `structs:"fruit"`
	Vegetables int `structs:"vegetables"`
}

type Diet struct {
	Values `structs:"value"`
	Code   string   `structs:"code"`
	Output []string `structs:"output"`
}

type Exercise struct {
	Current int      `structs:"value"`
	Code    string   `structs:"code"`
	Target  string   `structs:"target"`
	Output  []string `structs:"output"`
}

type Smoking struct {
	Code        string   `structs:"code"`
	Status      bool     `structs:"status"`
	SmokingCalc bool     `structs:"smoking_calc"`
	Output      []string `structs:"output"`
}

type Lifestyle struct {
	Bmi      `structs:"bmi"`
	Whr      `structs:"whr"`
	Diet     `structs:"diet"`
	Exercise `structs:"exercise"`
	Smoking  `structs:"smoking"`
}

type CvdAssessment struct {
	assessments.HighRisksAssessment `structs:"high_risk_condition"`
	CvdRisk                         map[string]string      `structs:"cvd_risk_result"`
	Guidelines                      map[string]interface{} `structs:"guidelines"`
}

// NewHearts returns a Hostname from a string.
func NewHearts(
	diab Diabetes,
	bp BloodPressure,
	lifestyle Lifestyle,
	cvd CvdAssessment) Hearts {
	return Hearts{
		Diabetes:      diab,
		BloodPressure: bp,
		Lifestyle:     lifestyle,
		CvdAssessment: cvd,
	}
}
