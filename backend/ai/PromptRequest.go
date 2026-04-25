package ai

// Text Should be a paragraph
// ID should start from 0
type PromptRequest struct {
	ID   int
	Text string
}

func NewPrompt(id int, prompt string) PromptRequest {
	return PromptRequest{
		ID:   id,
		Text: prompt,
	}
}
