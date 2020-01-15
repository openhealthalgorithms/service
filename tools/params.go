package tools

import (
	"strconv"
	"strings"
	"time"

	jp "github.com/buger/jsonparser"
	"github.com/pkg/errors"
)

// Configs object
type Configs struct {
	Algorithm string
	RiskModel string
	Debug     bool
}

// Demographics object
type Demographics struct {
	Age         float64
	DateOfBirth time.Time
	Gender      string
	Ethnicity   string
	Region      string
}

// Smoker object
type Smoker struct {
	CurrentSmoker  bool
	ExSmoker       bool
	QuitWithinYear bool
}

// MedicalHistory object
type MedicalHistory struct {
	Smoker
	Alcohol         float64
	AlcoholUnit     string
	AMI             bool
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
	ConditionNames  map[string]bool
}

// FamilyHistory object
type FamilyHistory struct {
	FamilyCvd bool
	FamilyCkd bool
}

// Medications object
type Medications map[string]bool

// Measurements object
type Measurements struct {
	Waist       float64
	WaistUnit   string
	Hip         float64
	HipUnit     string
	Height      float64
	HeightUnit  string
	Weight      float64
	WeightUnit  string
	BodyFat     float64
	BodyFatUnit string
	Sbp         int
	Dbp         int
	Pulse       int
	Bsl         float64
	BslUnit     string
	BslType     string
	A1C         float64
	TChol       float64
	Hdl         float64
	Ldl         float64
	Tg          float64
	CholUnit    string
	CholType    string
}

// Activities object
type Activities struct {
	PhysicalActivity int
	Fruits           int
	Vegetables       int
	Rice             int
	Oil              string
}

// type Assessments struct {
// 	assessments.DiabetesAssessment
// 	assessments.BPAssessment
// 	assessments.BMIAssessment
// 	assessments.WHRAssessment
// 	assessments.SmokingAssessment
// 	assessments.ExerciseAssessment
// 	assessments.DietAssessment
// 	assessments.HighRisksAssessment
// }

// Params object
type Params struct {
	Configs
	Demographics
	MedicalHistory
	FamilyHistory
	Medications
	Measurements
	Activities
	// Assessments
	Input []byte
}

// ParseParams function
func ParseParams(content []byte) (Params, error) {
	return getInputs(content)
}

func getInputs(data []byte) (Params, error) {
	var err error
	var intValue int64
	var stringValue string
	var floatValue float64

	out := Params{}
	mandatory := make([]string, 0)
	unsupported := make([]string, 0)

	out.Input = data

	// Config
	if stringValue, err = jp.GetString(data, "config", "algorithm"); err == nil {
		out.Algorithm = strings.ToLower(stringValue)
	} else {
		mandatory = append(mandatory, "algorithm")
	}

	if stringValue, err = jp.GetString(data, "config", "risk_model"); err == nil {
		out.RiskModel = strings.ToLower(stringValue)
	} else {
		mandatory = append(mandatory, "risk_model")
	}

	if val, err := jp.GetBoolean(data, "config", "debug"); err == nil {
		out.Debug = val
	} else {
		out.Debug = false
	}

	// Params
	// Demographics
	if stringValue, err = jp.GetString(data, "params", "demographics", "gender"); err == nil {
		out.Gender = strings.ToLower(stringValue[:1])
	} else {
		mandatory = append(mandatory, "gender")
	}

	age := 0.0
	if floatValue, err = jp.GetFloat(data, "params", "demographics", "age", "value"); err == nil {
		age = floatValue
	} else {
		mandatory = append(mandatory, "age")
	}

	unit := "year"
	if stringValue, err = jp.GetString(data, "params", "demographics", "age", "unit"); err == nil {
		unit = strings.ToLower(stringValue)
	} else {
		mandatory = append(mandatory, "age unit")
	}

	out.Age = CalculateAge(age, unit)

	if stringValue, err = jp.GetString(data, "params", "demographics", "birth_country_code"); err == nil {
		if code, ok := countries[stringValue]; ok {
			if code.Region != "#N/A" {
				out.Region = code.Region
			} else {
				unsupported = append(unsupported, stringValue)
			}
		} else {
			unsupported = append(unsupported, stringValue)
		}
	} else {
		mandatory = append(mandatory, "birth_country_code")
	}

	// Components
	// Lifestyle
	jp.ArrayEach(data, func(value []byte, dataType jp.ValueType, offset int, err error) {
		name := ""
		category := ""
		if stringValue, err = jp.GetString(value, "name"); err == nil {
			name = stringValue
		}

		if stringValue, err = jp.GetString(value, "category"); err == nil {
			category = stringValue
		}

		switch category {
		case "addiction":
			switch name {
			case "smoking":
				if val, err := jp.GetString(value, "value"); err == nil {
					if val == "smoker" {
						out.CurrentSmoker = true
					} else if val == "ex-smoker" {
						out.ExSmoker = true
					}
				}

				if val, err := jp.GetBoolean(value, "quit_within_year"); err == nil {
					out.QuitWithinYear = val
				}
			case "alcohol_history":
				frt := 0.0
				freq := ""
				if floatValue, err = jp.GetFloat(value, "value"); err == nil {
					frt = floatValue
				}
				if stringValue, err = jp.GetString(value, "frequency"); err == nil {
					freq = stringValue
				}
				out.Alcohol = CalculateAlcoholConsumption(frt, freq)
				out.AlcoholUnit = "weekly"
			}
		case "nutrition":
			switch name {
			case "fruit":
				frt := 0
				freq := ""
				if intValue, err = jp.GetInt(value, "value"); err == nil {
					frt = int(intValue)
				}
				if stringValue, err = jp.GetString(value, "frequency"); err == nil {
					freq = stringValue
				}
				out.Fruits = CalculateDietConsumption(frt, freq)
			case "vegetables":
				veg := 0
				freq := ""
				if intValue, err = jp.GetInt(value, "value"); err == nil {
					veg = int(intValue)
				}
				if stringValue, err = jp.GetString(value, "frequency"); err == nil {
					freq = stringValue
				}
				out.Vegetables = CalculateDietConsumption(veg, freq)
			case "rice":
				if intValue, err = jp.GetInt(value, "value"); err == nil {
					out.Rice = int(intValue)
				}
			case "oil":
				if stringValue, err = jp.GetString(value, "value"); err == nil {
					out.Oil = stringValue
				}
			}
		case "physical-activity":
			val := 0
			unit := ""
			freq := ""
			intensity := ""

			if intValue, err = jp.GetInt(value, "value"); err == nil {
				val = int(intValue)
			}
			if stringValue, err = jp.GetString(value, "units"); err == nil {
				unit = stringValue
			}
			if stringValue, err = jp.GetString(value, "frequency"); err == nil {
				freq = stringValue
			}
			if stringValue, err = jp.GetString(value, "intensity"); err == nil {
				intensity = stringValue
			}

			out.PhysicalActivity += CalculateExercise(val, unit, freq, intensity)
		}
	}, "params", "components", "lifestyle")

	// Body measurements
	sbpTotal := 0
	dbpTotal := 0
	bpCount := 0
	jp.ArrayEach(data, func(value []byte, dataType jp.ValueType, offset int, err error) {
		name := ""
		if stringValue, err = jp.GetString(value, "name"); err == nil {
			name = stringValue
		}

		switch name {
		case "height":
			units := ""
			if val, err := jp.GetString(value, "units"); err == nil {
				units = val
			}

			values := 0.0
			if val, err := jp.GetFloat(value, "value"); err == nil {
				values = val
			}

			out.Height = ConvertLength(values, units)
			out.HeightUnit = "m"
		case "weight":
			units := ""
			if val, err := jp.GetString(value, "units"); err == nil {
				units = val
			}

			values := 0.0
			if val, err := jp.GetFloat(value, "value"); err == nil {
				values = val
			}

			out.Weight = ConvertWeight(values, units)
			out.WeightUnit = "kg"
		case "hip":
			units := ""
			if val, err := jp.GetString(value, "units"); err == nil {
				units = val
			}

			values := 0.0
			if val, err := jp.GetFloat(value, "value"); err == nil {
				values = val
			}

			out.Hip = ConvertLength(values, units)
			out.HipUnit = "m"
		case "waist":
			units := ""
			if val, err := jp.GetString(value, "units"); err == nil {
				units = val
			}

			values := 0.0
			if val, err := jp.GetFloat(value, "value"); err == nil {
				values = val
			}

			out.Waist = ConvertLength(values, units)
			out.WaistUnit = "m"
		case "body-fat":
			units := ""
			if val, err := jp.GetString(value, "units"); err == nil {
				units = val
			}

			values := 0.0
			if val, err := jp.GetFloat(value, "value"); err == nil {
				values = val
			}

			out.BodyFat = ConvertLength(float64(values), units)
			out.BodyFatUnit = units
		case "blood_pressure":
			if stringValue, err = jp.GetString(value, "value"); err == nil {
				bps := strings.Split(stringValue, "/")
				sbp, _ := strconv.Atoi(bps[0])
				dbp, _ := strconv.Atoi(bps[1])
				sbpTotal += sbp
				dbpTotal += dbp
				bpCount++
			}
		}
	}, "params", "components", "body-measurements")

	if bpCount > 0 {
		out.Sbp = int(sbpTotal / bpCount)
		out.Dbp = int(dbpTotal / bpCount)
	} else {
		mandatory = append(mandatory, "blood_pressure")
	}

	// Biological Samples
	jp.ArrayEach(data, func(value []byte, dataType jp.ValueType, offset int, err error) {
		name := ""
		if stringValue, err = jp.GetString(value, "name"); err == nil {
			name = stringValue
		}

		switch name {
		case "blood_sugar":
			if val, err := jp.GetString(value, "units"); err == nil {
				out.BslUnit = val
			}

			if val, err := jp.GetFloat(value, "value"); err == nil {
				out.Bsl = val
			}

			if val, err := jp.GetString(value, "type"); err == nil {
				out.BslType = val
				out.CholType = val
			}
		case "a1c":
			if val, err := jp.GetFloat(value, "value"); err == nil {
				out.A1C = val
			}
		case "total_cholesterol":
			if val, err := jp.GetFloat(value, "value"); err == nil {
				out.TChol = val
			}
			if val, err := jp.GetString(value, "units"); err == nil {
				out.CholUnit = val
			}
		case "hdl":
			if val, err := jp.GetFloat(value, "value"); err == nil {
				out.Hdl = val
			}
		case "ldl":
			if val, err := jp.GetFloat(value, "value"); err == nil {
				out.Ldl = val
			}
		case "tg":
			if val, err := jp.GetFloat(value, "value"); err == nil {
				out.Tg = val
			}
		}
	}, "params", "components", "biological-samples")

	// Family History
	jp.ArrayEach(data, func(value []byte, dataType jp.ValueType, offset int, err error) {
		name := ""
		if stringValue, err = jp.GetString(value, "name"); err == nil {
			name = stringValue
		}

		switch name {
		case "cardiovascular-disease":
			out.FamilyCvd = true
		case "kidney-disease":
			out.FamilyCkd = true
		}
	}, "params", "components", "family_history")

	// Medications
	out.Medications = make(map[string]bool)
	jp.ArrayEach(data, func(value []byte, dataType jp.ValueType, offset int, err error) {
		category := ""
		if stringValue, err = jp.GetString(value, "category"); err == nil {
			category = stringValue
		}

		out.Medications[category] = true

		// switch category {
		// case "anti-hypertensive":
		// 	out.Antihypertensives = true
		// case "oral-hypoglycaemic":
		// 	out.OralHypoglycaemic = true
		// case "insulin":
		// 	out.Insulin = true
		// case "lipid-lowering":
		// 	out.LipidLowering = true
		// case "anti-platelet":
		// 	out.Antiplatelet = true
		// case "anti-coagulant":
		// 	out.AntiCoagulant = true
		// case "bronchodilator":
		// 	out.Bronchodilator = true
		// }
	}, "params", "components", "medications")

	// Medical history
	out.Conditions = make(map[string]bool)
	out.ConditionNames = make(map[string]bool)
	jp.ArrayEach(data, func(value []byte, dataType jp.ValueType, offset int, err error) {
		category := ""
		if stringValue, err = jp.GetString(value, "category"); err == nil {
			category = stringValue
		}

		switch category {
		case "condition":
			name := ""
			if stringValue, err = jp.GetString(value, "name"); err == nil {
				name = stringValue
			}
			conditionBool := false
			if boolValue, err := jp.GetBoolean(value, "is_active"); err == nil {
				conditionBool = boolValue
			}
			out.Conditions[strings.ToUpper(name)] = conditionBool
			out.ConditionNames[name] = conditionBool
			switch name {
			// case "asthma":
			// 	out.Asthma = conditionBool
			// case "tuberculosis":
			// 	out.Tuberculosis = conditionBool
			case "diabetes":
				out.Diabetes = conditionBool
				// case "hypertension":
				// 	out.Hypertension = conditionBool
				// case "ckd":
				// 	out.Ckd = conditionBool
				// case "cvd":
				// 	out.Cvd = conditionBool
				// case "pvd":
				// 	out.Pvd = conditionBool
				// case "pregnant":
				// 	out.Pregnant = conditionBool
				// case "arrhythmia":
				// 	out.Arrhythmia = conditionBool
			}
		}
	}, "params", "components", "medical_history")

	// // Calculations
	// // Diabetes
	// if value, err := assessments.GetDiabetes(out.Bsl, out.BslUnit, out.BslType, out.Diabetes); err == nil {
	// 	out.DiabetesAssessment = value
	// 	out.Diabetes = value.Status
	// } else {
	// 	out.DiabetesAssessment = assessments.DiabetesAssessment{}
	// }

	// // BMI
	// if value, err := assessments.GetBMI(out.Weight, out.Height); err == nil {
	// 	out.BMIAssessment = value
	// } else {
	// 	out.BMIAssessment = assessments.BMIAssessment{}
	// }

	// // WHR
	// if value, err := assessments.GetWHR(out.Waist, out.Hip, out.Gender); err == nil {
	// 	out.WHRAssessment = value
	// } else {
	// 	out.WHRAssessment = assessments.WHRAssessment{}
	// }

	// // BP
	// if value, err := assessments.GetBP(out.Sbp, out.Dbp, out.Diabetes); err == nil {
	// 	out.BPAssessment = value
	// } else {
	// 	out.BPAssessment = assessments.BPAssessment{}
	// }

	// // Smoking
	// if value, err := assessments.GetSmoking(out.CurrentSmoker, out.ExSmoker, out.QuitWithinYear); err == nil {
	// 	out.SmokingAssessment = value
	// } else {
	// 	out.SmokingAssessment = assessments.SmokingAssessment{}
	// }

	// // Exercise
	// if value, err := assessments.GetExercise(out.PhysicalActivity); err == nil {
	// 	out.ExerciseAssessment = value
	// } else {
	// 	out.ExerciseAssessment = assessments.ExerciseAssessment{}
	// }

	// // Diet
	// if value, err := assessments.GetDiet(out.Fruits, out.Vegetables); err == nil {
	// 	out.DietAssessment = value
	// } else {
	// 	out.DietAssessment = assessments.DietAssessment{}
	// }

	// // HighRisk
	// if value, err := assessments.GetHighRisks(out.Sbp, out.Dbp, out.Age, out.Conditions); err == nil {
	// 	out.HighRisksAssessment = value
	// } else {
	// 	out.HighRisksAssessment = assessments.HighRisksAssessment{}
	// }
	// fmt.Printf("%+v\n", out.Region)

	if len(mandatory) > 0 {
		return out, errors.Errorf("missing mandatory attributes: %v", JoinStringsSep(", ", mandatory...))
	}

	return out, err
}
