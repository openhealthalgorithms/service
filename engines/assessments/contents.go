package assessments

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
}

// BodyContents object
type BodyContents struct {
    Contents *Contents `json:"contents"`
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
