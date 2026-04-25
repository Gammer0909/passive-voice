package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"google.golang.org/genai"
)

const MODEL_NAME = "gemma-4-31b-it"

type AI struct {
	Context       context.Context
	Client        *genai.Client
	ContentConfig *genai.GenerateContentConfig
}

func NewAI(ctx context.Context, config *genai.GenerateContentConfig) (*AI, error) {

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &AI{
		Context:       ctx,
		Client:        client,
		ContentConfig: config,
	}, nil
}

func (ai *AI) AnalyzeParagraph(id int, paragraph string) (*AnalyzedParagraph, error) {

	log.Printf("In AnalyzeParagraph():\n------\n%s\n-------------------------\n", paragraph)

	resp, err := ai.prompt(NewPrompt(id, paragraph))
	if err != nil {
		return nil, err
	}

	return resp, nil

}

func (ai *AI) prompt(req PromptRequest) (*AnalyzedParagraph, error) {

	result, err := ai.Client.Models.GenerateContent(
		ai.Context,
		MODEL_NAME,
		genai.Text(req.Text),
		ai.ContentConfig,
	)
	if err != nil {
		return nil, err
	}

	// so for some reason the LLM likes to put "```" at the end of it's JSON
	// (see end of file for a guess as to why)
	// I tried a lazy prompt to get it to stop, but I realized a more reliable approach was just cutting it out instead
	jsonString := strings.TrimSpace(result.Text())
	jsonString = strings.ReplaceAll(jsonString, "```json", "")
	jsonString = strings.ReplaceAll(jsonString, "```", "")

	var resp *AnalyzedParagraph
	err = json.Unmarshal([]byte(jsonString), &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil

}

// /api/upload/
func (model *AI) HandleFileUpload(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodPost {
		http.Error(w, "/api/upload is a POST.", http.StatusMethodNotAllowed)
		return
	}

	log.Println("User sent a request.")

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("An error occurred getting the file from the form: ", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	dst, err := os.Create(fmt.Sprintf("./uploads/%s", handler.Filename))
	if err != nil {
		log.Println("An error occurred while creating the file on disk: ", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Println("An error occurred copying the file: ", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	doc, err := NewDocument(dst)
	if err != nil {
		log.Println("go-docx could not open the file:", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	log.Printf("All paragraphs:\n-------------------------------\n%#v\n-------------------------", doc.Sentences)

	analyzedParagraphs := make([]*AnalyzedParagraph, 0)

	for i, paragraph := range doc.Paragraphs {
		if paragraph == "" {
			log.Printf("Skipped %s because it's empty", paragraph)
			continue
		}

		log.Printf("In for loop:\n------\n%s\n-------------------------\n", paragraph)
		analyzed, err := model.AnalyzeParagraph(i, paragraph)
		if err != nil {
			log.Printf("Error occurred analyzing paragraph number %b: %v", i, err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
		analyzedParagraphs = append(analyzedParagraphs, analyzed)
	}

	jsonData, err := json.Marshal(analyzedParagraphs)
	if err != nil {
		log.Println("An error occurred marshaling data to JSON:", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

/*

My guess as to why the LLM enjoys adding "```" is because it's probably been trained on either a lot of previous LLM output
that includes that, or its been trained on a lot of articles or MDs that include ``` tags so it assumes they should come next even
though I said not to.

Basically, google's probably yoinked some training data from somewhere it shouldn't or from its own LLM
(the latter makes a lot of sense cause im using the open-source model)

*/
