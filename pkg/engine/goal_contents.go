package engine

import (
	"github.com/openhealthalgorithms/service/pkg/tools"
)

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
	TFL     *string `json:"tfl"`
	Message *string `json:"message"`
}

// GenerateGoalsGuideline function
func (g *GoalGuideContents) GenerateGoalsGuideline(codes ...string) []GoalGuidelinesContent {
	ggc := make([]GoalGuidelinesContent, 0)

	for k, v := range *g.Body.GoalGuidelinesContents {
		_, found := tools.SliceContainsString(codes, k)
		if found {
			ggc = append(ggc, v)
		}
	}

	return ggc
}
