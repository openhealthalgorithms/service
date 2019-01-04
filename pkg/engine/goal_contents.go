package engine

// GoalGuideContents object
type GoalGuideContents struct {
	Meta *GoalGuidelinesMetaContents `json:"meta"`
	Body *GoalGuidelinesBodyContents `json:"body"`
}

// GoalGuidelinesMetaContents object
type GoalGuidelinesMetaContents struct {
	GoalGuidelineName *string `json:"goal_name"`
	Publisher         *string `json:"publisher"`
	PublicationDate   *string `json:"publication_date"`
	ContentType       *string `json:"content_type"`
	PublishedBy       *string `json:"published_by"`
	VersionNumber     *string `json:"version_number"`
}

// GoalGuidelinesBodyContents object
type GoalGuidelinesBodyContents struct {
	GoalGuidelinesContents *GoalGuidelinesContents `json:"contents"`
}

// GoalGuidelinesContents map
type GoalGuidelinesContents map[string]GoalGuidelinesContent

// GoalGuidelinesContent object
type GoalGuidelinesContent struct {
	Eval    *string `json:"eval"`
	Grading *int    `json:"grading"`
	TFL     *string `json:"tfl"`
	Message *string `json:"message"`
	Refer   *string `json:"refer"`
}
