package ai

type AnalyzedSentence struct {
	Content        string  `json:"content"`
	IsPassiveVoice bool    `json:"isPassiveVoice"`
	PotentialFix   *string `json:"potentialFix"`
}
