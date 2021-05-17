package assessments

import (
	"github.com/openhealthalgorithms/service/tools"
)

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
	Contents    *Contents    `json:"contents"`
	MessagePool *MessagePool `json:"message-pool"`
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
	Version *string `json:"version"`
}

// MessagePool object
type MessagePool []MessageRules

// MessageRules object
type MessageRules struct {
	Assessment *string            `json:"assessment"`
	Conditions []MessageCondition `json:"conditions"`
	Message    *string            `json:"message"`
}

// MessageCondition object
type MessageCondition []string

// Process codes for MessagePool check
func (mp *MessagePool) Process(codes []string, assessment string) string {
	message := ""

	for _, m := range *mp {
		if *m.Assessment == assessment {
			match := true
			for _, mc := range m.Conditions {
				match = match && tools.SliceContainsAnyString(mc, codes)
			}
			if match {
				message = *m.Message
				break
			}
		}
	}

	return message
}
