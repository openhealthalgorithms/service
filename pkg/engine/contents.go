package engine

// GuideContents object
type GuideContents struct {
	Meta *MetaContents `json:"meta"`
	Body *BodyContents `json:"body"`
}

// MetaContents object
type MetaContents struct {
	GuidelineName   *string `json:"guideline_name"`
	Publisher       *string `json:"publisher"`
	PublicationDate *string `json:"publication_date"`
	ContentType     *string `json:"content_type"`
	PublishedBy     *string `json:"published_by"`
	VersionNumber   *string `json:"version_number"`
}

// BodyContents object
type BodyContents struct {
	Contents *Contents `json:"contents"`
	Gradings *Gradings `json:"gradings"`
}

// Contents map
type Contents map[string]Content

// Content object
type Content struct {
	Eval    *string `json:"eval"`
	Grading *int    `json:"grading"`
	TFL     *string `json:"tfl"`
	Message *string `json:"message"`
	Refer   *string `json:"refer"`
}

// Gradings object
type Gradings struct {
	BodyComposition *GradingScales `json:"body-composition"`
	Lifestyle       *GradingScales `json:"lifestyle"`
	Diet            *GradingScales `json:"diet"`
	Cholesterol     *GradingScales `json:"cholesterol"`
}

// GradingScales slice
type GradingScales []GradingScale

// GradingScale object
type GradingScale struct {
	Grading *RangeInt `json:"grading"`
	Message *string   `json:"message"`
}
