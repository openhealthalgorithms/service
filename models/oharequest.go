package models

// OHARequest object
type OHARequest struct {
	Config *ORConfig `json:"config" validate:"required"`
	Params *ORParams `json:"params" validate:"required"`
}

type (
	// ORConfig object
	ORConfig struct {
		Algorithm        *string `json:"algorithm" validate:"required"`
		RiskModel        *string `json:"risk_model" validate:"required"`
		Debug            *bool   `json:"debug" validate:""`
		CarePlan         *bool   `json:"careplan" validate:""`
		RiskModelVersion *string `json:"risk_model_version" validate:""`
		LabBased         *bool   `json:"lab_based" validate:""`
	}

	// ORParams object
	ORParams struct {
		Demographics *ORDemographics    `json:"demographics" validate:"required"`
		Components   *ORParamComponents `json:"components" validate:"required"`
	}

	// ORAge object
	ORAge struct {
		Value *int64  `json:"value" validate:"required"`
		Unit  *string `json:"unit" validate:"required"`
	}

	// ORDemographics object
	ORDemographics struct {
		Gender            *string `json:"gender" validate:"required"`
		Age               *ORAge  `json:"age" validate:"required"`
		BirthCountry      *string `json:"birth_country" validate:""`
		BirthCountryCode  *string `json:"birth_country_code" validate:"required"`
		LivingCountry     *string `json:"living_country" validate:""`
		LivingCountryCode *string `json:"living_country_code" validate:""`
		Race              *string `json:"race" validate:""`
		Ethnicity         *string `json:"ethnicity" validate:""`
	}

	// ORParamComponents object
	ORParamComponents struct {
		Lifestyle         []ORLifestyle        `json:"lifestyle" validate:"required"`
		BodyMeasurements  []ORBodyMeasurement  `json:"body-measurements" validate:"required"`
		BiologicalSamples []ORBiologicalSample `json:"biological-samples" validate:"required"`
		MedicalHistory    []ORMedicalHistory   `json:"medical_history" validate:""`
		Medications       []ORMedication       `json:"medications" validate:""`
		FamilyHistory     []ORFamilyHistory    `json:"family_history" validate:""`
	}

	// ORLifestyle object
	ORLifestyle struct {
		Name           *string      `json:"name" validate:"required"`
		Category       *string      `json:"category" validate:"required"`
		Value          *interface{} `json:"value" validate:"required"`
		Units          *string      `json:"units" validate:""`
		Frequency      *string      `json:"frequency" validate:""`
		Intensity      *string      `json:"intensity" validate:""`
		QuitWithinYear *bool        `json:"quit_within_year" validate:""`
	}

	// ORBodyMeasurement object
	ORBodyMeasurement struct {
		EffectiveDate *string      `json:"effectiveDate" validate:"required"`
		Name          *string      `json:"name" validate:"required"`
		Category      *string      `json:"category" validate:"required"`
		Value         *interface{} `json:"value" validate:"required"`
		Units         *string      `json:"units" validate:"required"`
		Arm           *string      `json:"arm" validate:""`
	}

	// ORBiologicalSample object
	ORBiologicalSample struct {
		EffectiveDate *string      `json:"effectiveDate" validate:"required"`
		Name          *string      `json:"name" validate:"required"`
		Category      *string      `json:"category" validate:"required"`
		Value         *interface{} `json:"value" validate:"required"`
		Units         *string      `json:"units" validate:"required"`
		Type          *string      `json:"type" validate:""`
	}

	// ORMedicalHistory object
	ORMedicalHistory struct {
		Name        *string `json:"name" validate:""`
		Category    *string `json:"category" validate:"required"`
		IsActive    *bool   `json:"is_active" validate:""`
		Type        *string `json:"type" validate:""`
		Allergen    *string `json:"allergen" validate:""`
		Criticality *string `json:"criticality" validate:""`
		Reaction    *string `json:"reaction" validate:""`
	}

	// ORMedication object
	ORMedication struct {
		Generic   *string `json:"generic" validate:""`
		Category  *string `json:"category" validate:"required"`
		Class     *string `json:"class" validate:"required"`
		Status    *string `json:"status" validate:""`
		Dose      *string `json:"dose" validate:""`
		Frequency *string `json:"frequency" validate:""`
	}

	// ORFamilyHistory object
	ORFamilyHistory struct {
		Name     *string `json:"name" validate:"required"`
		Relative *string `json:"relative" validate:"required"`
	}
)
