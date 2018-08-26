package tools

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/openhealthalgorithms/service/pkg/types"
)

type Demographics struct {
	Age         float64
	DateOfBirth time.Time
	Gender      string
	Ethnicity   string
	Region      string
}

type MedicalHistory struct {
	CurrentSmoker   bool
	AlcoholDaily    int
	AlcoholMax      int
	AlcoholFreeDays int
	Diabetes        bool
	Hypertension    bool
	HighCholesterol bool
	HighBsl         bool
	Ckd             bool
	Cvd             bool
	Pvd             bool
	Arrhythmia      bool
	Pregnant        bool
}

type FamilyHistory struct {
	FamilyCvd bool
	FamilyCkd bool
}

type Medications struct {
	Antihypertensives bool
}

type Measurements struct {
	Waist    float64
	Hip      float64
	Height   float64
	Weight   float64
	Sbp      int
	Dbp      int
	Pulse    int
	Bsl      float64
	BslUnit  string
	TChol    float64
	Hdl      float64
	Ldl      float64
	Tg       float64
	CholUnit string
}

type Params struct {
	Demographics
	MedicalHistory
	FamilyHistory
	Medications
	Measurements
}

func ParseParams(ctx context.Context) Params {
	v := ctx.Value(types.KeyValuesCtx).(*types.ValuesCtx)

	patternsRaw, ok := v.Params.Get("params")
	if !ok {
		return Params{}
	}

	pts, ok := patternsRaw.(string)
	if !ok {
		return Params{}
	}

	inputs := getInputs(pts)

	return inputs
}

func getInputs(input string) Params {
	out := Params{}

	tmp := strings.Split(input, ",")
	for _, t := range tmp {
		v := strings.Split(t, ":")
		key := strings.ToLower(v[0])
		value := strings.ToLower(v[1])

		switch key {
		case "gender":
			out.Gender = "f"
			if value == "male" || value == "m" {
				out.Gender = "m"
			}
		case "age":
			out.Age, _ = strconv.ParseFloat(value, 64)
		case "sbp":
			sbps := strings.Split(value, "|")
			total := 0
			for _, s := range sbps {
				c, _ := strconv.Atoi(s)
				total += c
			}
			out.Sbp = int(total / len(sbps))
		case "dbp":
			dbps := strings.Split(value, "|")
			total := 0
			for _, d := range dbps {
				c, _ := strconv.Atoi(d)
				total += c
			}
			out.Dbp = int(total / len(dbps))
		case "smoker":
			out.CurrentSmoker = false
			if value == "true" || value == "1" || value == "yes" {
				out.CurrentSmoker = true
			}
		case "diabetic":
			out.Diabetes = false
			if value == "true" || value == "1" || value == "yes" {
				out.Diabetes = true
			}
		case "region":
			out.Region = strings.ToUpper(value)
		case "tchol":
			out.TChol, _ = strconv.ParseFloat(value, 64)
		case "cholunit":
			out.CholUnit = "mmol"
			if strings.HasPrefix(value, "mg") {
				out.CholUnit = "mgdl"
			}
		case "waist":
			out.Waist, _ = strconv.ParseFloat(value, 64)
		case "hip":
			out.Hip, _ = strconv.ParseFloat(value, 64)
		case "height":
			out.Height, _ = strconv.ParseFloat(value, 64)
		case "weight":
			out.Weight, _ = strconv.ParseFloat(value, 64)
		}
	}

	return out
}
