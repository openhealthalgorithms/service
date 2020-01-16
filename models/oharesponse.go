package models

type (
    ORRAssessments struct {
        BloodPressure   *ORRAssessment      `json:"blood_pressure"`
        BodyComposition *ORRBodyComposition `json:"body_composition"`
        Cholesterol     *ORRCholesterol     `json:"cholesterol"`
        CVD             *ORRAssessment      `json:"cvd"`
        Diabetes        *ORRAssessment      `json:"diabetes"`
        Lifestyle       *ORRLifestyle       `json:"lifestyle"`
    }

    ORRBodyComposition struct {
        Components *ORRBodyCompositionComponents `json:"components"`
        Message    *string                       `json:"message"`
    }

    ORRBodyCompositionComponents struct {
        BMI       *ORRAssessment `json:"bmi"`
        BodyFat   *ORRAssessment `json:"body_fat"`
        WaistCirc *ORRAssessment `json:"waist_circ"`
        WHR       *ORRAssessment `json:"whr"`
    }

    ORRCholesterol struct {
        Components *ORRCholesterolComponents `json:"components"`
        Message    *string                   `json:"message"`
    }

    ORRCholesterolComponents struct {
        HDL   *ORRAssessment `json:"hdl"`
        LDL   *ORRAssessment `json:"ldl"`
        TG    *ORRAssessment `json:"tg"`
        TChol *ORRAssessment `json:"total_cholesterol"`
    }

    ORRLifestyle struct {
        Components *ORRLifestyleComponents `json:"components"`
        Message    *string                 `json:"message"`
    }

    ORRDiet struct {
        Components *ORRDietComponents `json:"components"`
        Message    *string            `json:"message"`
    }

    ORRDietComponents struct {
        Fruit          *ORRAssessment `json:"fruit"`
        FruitVegetable *ORRAssessment `json:"fruit_vegetable"`
        Vegetable      *ORRAssessment `json:"vegetable"`
    }

    ORRLifestyleComponents struct {
        Alcohol          *ORRAssessment `json:"alcohol"`
        Diet             *ORRDiet       `json:"diet"`
        PhysicalActivity *ORRAssessment `json:"physical_activity"`
        Smoking          *ORRAssessment `json:"smoking"`
    }

    ORRAssessment struct {
        Code    *string `json:"code" validate:"required"`
        Eval    *string `json:"eval" validate:"required"`
        Message *string `json:"message" validate:"required"`
        Refer   *string `json:"refer" validate:"required"`
        Target  *string `json:"target" validate:"required"`
        TFL     *string `json:"tfl" validate:"required"`
        Value   *string `json:"value" validate:"required"`
    }

    ORRGoal struct {
        Code    *string `json:"code" validate:"required"`
        Eval    *string `json:"eval" validate:"required"`
        TFL     *string `json:"tfl" validate:"required"`
        Message *string `json:"message" validate:"required"`
    }

    ORRReferrals struct {
        Reasons []ORRReferralReason `json:"reasons" validate:""`
        Refer   *bool               `json:"refer" validate:""`
        Urgent  *bool               `json:"urgent" validate:""`
    }

    ORRReferralReason struct {
        Type   *string `json:"type" validate:"required"`
        Urgent *bool   `json:"urgent" validate:"required"`
    }
)

func NewORRReferrals() *ORRReferrals {
    return &ORRReferrals{
        Reasons: nil,
        Refer:   nil,
        Urgent:  nil,
    }
}

func NewORRAssessments() *ORRAssessments {
    return &ORRAssessments{
        BloodPressure: &ORRAssessment{},
        BodyComposition: &ORRBodyComposition{
            Components: &ORRBodyCompositionComponents{
                BMI:       nil,
                BodyFat:   nil,
                WaistCirc: nil,
                WHR:       nil,
            },
            Message: nil,
        },
        Cholesterol: &ORRCholesterol{
            Components: &ORRCholesterolComponents{
                HDL:   nil,
                LDL:   nil,
                TG:    nil,
                TChol: nil,
            },
            Message: nil,
        },
        CVD:      &ORRAssessment{},
        Diabetes: &ORRAssessment{},
        Lifestyle: &ORRLifestyle{
            Components: &ORRLifestyleComponents{
                Alcohol: nil,
                Diet: &ORRDiet{
                    Components: &ORRDietComponents{
                        Fruit:          nil,
                        FruitVegetable: nil,
                        Vegetable:      nil,
                    },
                    Message: nil,
                },
                PhysicalActivity: nil,
                Smoking:          nil,
            },
            Message: nil,
        },
    }
}
