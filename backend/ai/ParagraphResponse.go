package ai

type AnalyzedParagraph struct {
	Paragraph string              `json:"paragraph"`
	Sentences []*AnalyzedSentence `json:"sentences"`
}
