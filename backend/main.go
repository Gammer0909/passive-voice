package main

import (
	"context"
	"log"
	"net/http"
	"passive-voice-detect/ai"

	"google.golang.org/genai"
)

func main() {
	contentConfig := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseJsonSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"paragraph": map[string]any{
					"type":        "string",
					"description": "The original paragraph text.",
				},
				"sentences": map[string]any{
					"type": "array",
					"items": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"content": map[string]any{
								"type":        "string",
								"description": "The original sentence content.",
							},
							"isPassiveVoice": map[string]any{
								"type":        "boolean",
								"description": "Whether the sentence uses passive voice.",
							},
							"potentialFix": map[string]any{
								"type":        "string",
								"description": "A suggested revision to convert to active voice. Null if sentence is already active.",
							},
						},
						"required": []string{"content", "isPassiveVoice"},
					},
				},
			},
			"required": []string{"paragraph", "sentences"},
		},
	}

	model, err := ai.NewAI(context.Background(), contentConfig)
	if err != nil {
		log.Fatal("Could not create instance of AI:", err)
		return
	}

	http.HandleFunc("/api/upload", model.HandleFileUpload)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
