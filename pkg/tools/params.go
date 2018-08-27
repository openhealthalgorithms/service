package tools

import (
	"strconv"
	"strings"
	"time"

	jp "github.com/buger/jsonparser"

	"github.com/openhealthalgorithms/service/pkg/assessments"
)

type Demographics struct {
	Age         float64
	DateOfBirth time.Time
	Gender      string
	Ethnicity   string
	Region      string
}

type Smoker struct {
	CurrentSmoker  bool
	ExSmoker       bool
	QuitWithinYear bool
}

type MedicalHistory struct {
	Smoker
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
	Asthma          bool
	Tuberculosis    bool
	Conditions      map[string]bool
}

type FamilyHistory struct {
	FamilyCvd bool
	FamilyCkd bool
}

type Medications struct {
	Antihypertensives bool
	Statin            bool
	Antiplatelet      bool
	Bronchodilator    bool
}

type Measurements struct {
	Waist      float64
	WaistUnit  string
	Hip        float64
	HipUnit    string
	Height     float64
	HeightUnit string
	Weight     float64
	WeightUnit string
	Sbp        int
	Dbp        int
	Pulse      int
	Bsl        float64
	BslUnit    string
	BslType    string
	TChol      float64
	Hdl        float64
	Ldl        float64
	Tg         float64
	CholUnit   string
	CholType   string
}

type Activities struct {
	PhysicalActivity int
	Fruits           int
	Vegetables       int
	Rice             int
	Oil              string
}

type Assessments struct {
	assessments.DiabetesAssessment
	assessments.BPAssessment
	assessments.BMIAssessment
	assessments.WHRAssessment
	assessments.SmokingAssessment
	assessments.ExerciseAssessment
	assessments.DietAssessment
	assessments.HighRisksAssessment
}

type Params struct {
	Demographics
	MedicalHistory
	FamilyHistory
	Medications
	Measurements
	Activities
	Assessments
}

func ParseParams(content []byte) Params {
	inputs := getInputs(content)

	return inputs
}

func getInputs(data []byte) Params {
	out := Params{}

	if value, err := jp.GetString(data, "demographics", "gender"); err == nil {
		out.Gender = strings.ToLower(value)
	}

	if value, err := jp.GetFloat(data, "demographics", "age"); err == nil {
		out.Age = value
	}

	if value, err := jp.GetString(data, "region"); err == nil {
		out.Region = value
	}

	temp := 0
	count := 0
	jp.ArrayEach(data, func(value []byte, dataType jp.ValueType, offset int, err error) {
		c, e := strconv.Atoi(string(value))
		if e == nil {
			temp += c
			count++
		}
	}, "measurements", "sbp")

	out.Sbp = int(temp / count)

	temp = 0
	count = 0
	jp.ArrayEach(data, func(value []byte, dataType jp.ValueType, offset int, err error) {
		c, e := strconv.Atoi(string(value))
		if e == nil {
			temp += c
			count++
		}
	}, "measurements", "dbp")

	out.Dbp = int(temp / count)

	if value, err := jp.GetFloat(data, "measurements", "height", "[0]"); err == nil {
		out.Height = value
	}

	if value, err := jp.GetString(data, "measurements", "height", "[1]"); err == nil {
		out.HeightUnit = value
	}

	if value, err := jp.GetFloat(data, "measurements", "weight", "[0]"); err == nil {
		out.Weight = value
	}

	if value, err := jp.GetString(data, "measurements", "weight", "[1]"); err == nil {
		out.WeightUnit = value
	}

	if value, err := jp.GetFloat(data, "measurements", "waist", "[0]"); err == nil {
		out.Waist = value
	}

	if value, err := jp.GetString(data, "measurements", "waist", "[1]"); err == nil {
		out.WaistUnit = value
	}

	if value, err := jp.GetFloat(data, "measurements", "hip", "[0]"); err == nil {
		out.Hip = value
	}

	if value, err := jp.GetString(data, "measurements", "hip", "[1]"); err == nil {
		out.HipUnit = value
	}

	if value, err := jp.GetInt(data, "smoking", "current"); err == nil {
		out.CurrentSmoker = false
		if value == 1 {
			out.CurrentSmoker = true
		}
	}

	if value, err := jp.GetInt(data, "smoking", "ex_smoker"); err == nil {
		out.ExSmoker = false
		if value == 1 {
			out.ExSmoker = true
		}
	}

	if value, err := jp.GetInt(data, "smoking", "quit_within_year"); err == nil {
		out.QuitWithinYear = false
		if value == 1 {
			out.QuitWithinYear = true
		}
	}

	jp.ArrayEach(data, func(value []byte, dataType jp.ValueType, offset int, err error) {
		s := string(value)
		if s == "cvd" {
			out.FamilyCvd = true
		} else if s == "ckd" {
			out.FamilyCkd = true
		}
	}, "family_history")

	jp.ArrayEach(data, func(value []byte, dataType jp.ValueType, offset int, err error) {
		s := string(value)
		if s == "anti_hypertensive" {
			out.Antihypertensives = true
		} else if s == "statin" {
			out.Statin = true
		} else if s == "antiplatelet" {
			out.Antiplatelet = true
		} else if s == "bronchodilator" {
			out.Bronchodilator = true
		}
	}, "medications")

	out.Conditions = make(map[string]bool)
	jp.ArrayEach(data, func(value []byte, dataType jp.ValueType, offset int, err error) {
		s := string(value)
		if s == "asthma" {
			out.Asthma = true
		} else if s == "tuberculosis" {
			out.Tuberculosis = true
		} else if s == "diabetes" {
			out.Diabetes = true
		} else if s == "hypertension" {
			out.Hypertension = true
		} else if s == "ckd" {
			out.Ckd = true
		} else if s == "cvd" {
			out.Cvd = true
		} else if s == "pvd" {
			out.Pvd = true
		} else if s == "pregnant" {
			out.Pregnant = true
		} else if s == "arrhythmia" {
			out.Arrhythmia = true
		}

		out.Conditions[strings.ToUpper(s)] = true
	}, "medical_history", "conditions")

	if value, err := jp.GetString(data, "pathology", "bsl", "type"); err == nil {
		out.BslType = value
	}

	if value, err := jp.GetString(data, "pathology", "bsl", "units"); err == nil {
		out.BslUnit = value
	}

	if value, err := jp.GetFloat(data, "pathology", "bsl", "value"); err == nil {
		out.Bsl = value
	}

	if value, err := jp.GetString(data, "pathology", "cholesterol", "type"); err == nil {
		out.CholType = value
	}

	if value, err := jp.GetString(data, "pathology", "cholesterol", "units"); err == nil {
		out.CholUnit = value
	}

	if value, err := jp.GetFloat(data, "pathology", "cholesterol", "total_chol"); err == nil {
		out.TChol = value
	}

	if value, err := jp.GetFloat(data, "pathology", "cholesterol", "hdl"); err == nil {
		out.Hdl = value
	}

	if value, err := jp.GetFloat(data, "pathology", "cholesterol", "ldl"); err == nil {
		out.Ldl = value
	}

	if value, err := jp.GetFloat(data, "pathology", "cholesterol", "tg"); err == nil {
		out.Tg = value
	}

	if value, err := jp.GetString(data, "physical_activity"); err == nil {
		v, e := strconv.Atoi(value)
		if e != nil {
			out.PhysicalActivity = 0
		} else {
			out.PhysicalActivity = v
		}
	}

	if value, err := jp.GetInt(data, "diet_history", "fruit"); err == nil {
		out.Fruits = int(value)
	}

	if value, err := jp.GetInt(data, "diet_history", "veg"); err == nil {
		out.Vegetables = int(value)
	}

	if value, err := jp.GetInt(data, "diet_history", "rice"); err == nil {
		out.Rice = int(value)
	}

	if value, err := jp.GetString(data, "diet_history", "oil"); err == nil {
		out.Oil = value
	}

	// Calculations
	// Diabetes
	if value, err := assessments.GetDiabetes(out.Bsl, out.BslUnit, out.BslType, out.Diabetes); err == nil {
		out.DiabetesAssessment = value
		out.Diabetes = value.Status
	}

	// BMI
	if value, err := assessments.GetBMI(out.Weight, out.Height); err == nil {
		out.BMIAssessment = value
	}

	// WHR
	if value, err := assessments.GetWHR(out.Waist, out.Hip, out.Gender); err == nil {
		out.WHRAssessment = value
	}

	// BP
	if value, err := assessments.GetBP(out.Sbp, out.Dbp, out.Diabetes); err == nil {
		out.BPAssessment = value
	}

	// Smoking
	if value, err := assessments.GetSmoking(out.CurrentSmoker, out.ExSmoker, out.QuitWithinYear); err == nil {
		out.SmokingAssessment = value
	}

	// Exercise
	if value, err := assessments.GetExercise(out.PhysicalActivity); err == nil {
		out.ExerciseAssessment = value
	}

	// Diet
	if value, err := assessments.GetDiet(out.Fruits, out.Vegetables); err == nil {
		out.DietAssessment = value
	}

	// HighRisk
	if value, err := assessments.GetHighRisks(out.Sbp, out.Dbp, out.Age, out.Conditions); err == nil {
		out.HighRisksAssessment = value
	}

	return out
}
