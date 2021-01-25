package config

// AgeRange object
type AgeRange struct {
	From int    `json:"from"`
	To   int    `json:"to"`
	Key  string `json:"key"`
}

// CholRange object
type CholRange struct {
	From float64 `json:"from"`
	To   float64 `json:"to"`
}

// SBPRange object
type SBPRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// RiskRange object
type RiskRange struct {
	From  int    `json:"from"`
	To    int    `json:"to"`
	Value string `json:"value"`
}

// Configs object
type Configs struct {
	Ages         []AgeRange  `json:"ages"`
	Cholesterols []CholRange `json:"cholesterols"`
	BMIs         []CholRange `json:"bmis"`
	Systolics    []SBPRange  `json:"systolics"`
	RiskValues   []RiskRange `json:"riskValues"`
}

// Meta object
type Meta struct {
	Version       string  `json:"version"`
	WHOVersion    string  `json:"whoversion"`
	Configuration Configs `json:"config"`
}

// AgeRanges object
type AgeRanges [][]int

// Smoking object
type Smoking map[string]AgeRanges

// Gender object
type Gender struct {
	NonSmoker Smoking `json:"nonsmoker"`
	Smoker    Smoking `json:"smoker"`
}

// Diabetes object
type Diabetes struct {
	Male   Gender `json:"male"`
	Female Gender `json:"female"`
}

// Cholesterol object
type Cholesterol struct {
	Diabetic    Diabetes `json:"diabetic"`
	NonDiabetic Diabetes `json:"nondiabetic"`
}

// Regions object
type Regions struct {
	Chol    Cholesterol `json:"cholesterol"`
	NonChol Cholesterol `json:"noncholesterol"`
}

// WHOColorChartV1 object
type WHOColorChartV1 struct {
	Meta        Meta               `json:"meta"`
	ColorCharts map[string]Regions `json:"charts"`
}

// NonLab object
type NonLab struct {
	Male   Gender `json:"male"`
	Female Gender `json:"female"`
}

// RegionsV2 object
type RegionsV2 struct {
	LabBased    Cholesterol `json:"lab_based"`
	NonLabBased NonLab      `json:"non_lab_based"`
}

// WHOColorChartV2 object
type WHOColorChartV2 struct {
	Meta        Meta                 `json:"meta"`
	ColorCharts map[string]RegionsV2 `json:"charts"`
}
