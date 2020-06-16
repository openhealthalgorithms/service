package assessments

import "github.com/pkg/errors"

// Guidelines object
type Guidelines struct {
	Meta *MetaStructure `json:"meta"`
	Body *BodyStructure `json:"body"`
}

// MetaStructure object
type MetaStructure struct {
	GuidelineName   *string `json:"guideline_name"`
	Publisher       *string `json:"publisher"`
	PublicationDate *string `json:"publication_date"`
}

// BodyStructure object
type BodyStructure struct {
	Lifestyle       *LifestyleGuideline        `json:"lifestyle"`
	BodyComposition *BodyCompositionGuideline  `json:"body_composition"`
	BloodPressure   *BloodPressureGuidelines   `json:"blood_pressure"`
	Diabetes        *DiabetesGuidelines        `json:"diabetes"`
	Cholesterol     *CholesterolGuidelinesFull `json:"cholesterol"`
	CVD             *CVDGuidelinesFull         `json:"cvd"`
}

// Response object
type Response struct {
	Code   string `json:"code"`
	Value  string `json:"value"`
	Target string `json:"target"`
}

// GetResponse function
func GetResponse(assessmentName, code, value, target string) (Response, error) {
	if len(code) > 0 {
		resp := Response{
			Code:   code,
			Value:  value,
			Target: target,
		}
		return resp, nil
	}
	return Response{}, errors.New("no matching condition found for " + assessmentName)
}
