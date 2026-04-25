package ai

import (
	"os"
	"regexp"
	"unicode"

	"github.com/fumiama/go-docx"
)

type WordDocument struct {
	Sentences  []string
	Analyzed   []*AnalyzedSentence
	Paragraphs []string
}

func NewDocument(readFile *os.File) (*WordDocument, error) {
	fileinfo, err := readFile.Stat()
	if err != nil {
		panic(err)
	}
	size := fileinfo.Size()
	doc, err := docx.Parse(readFile, size)
	if err != nil {
		panic(err)
	}

	wordDoc := &WordDocument{
		Sentences:  make([]string, 0),
		Analyzed:   make([]*AnalyzedSentence, 0),
		Paragraphs: make([]string, 0),
	}

	for _, it := range doc.Document.Body.Items {
		switch it := it.(type) {
		case *docx.Paragraph:
			{
				wordDoc.Paragraphs = append(wordDoc.Paragraphs, it.String())
			}

		}
	}

	for _, paragraph := range wordDoc.Paragraphs {
		sentences := ParagraphToSentences(paragraph)
		for _, sentence := range sentences {
			// w nesting for loopsx
			wordDoc.Sentences = append(wordDoc.Sentences, sentence)
		}
	}

	return wordDoc, nil
}

func ParagraphToSentences(text string) []string {
	re := regexp.MustCompile(`[.!?]+\s+`)
	indices := re.FindAllStringIndex(text, -1)

	var result []string
	last := 0

	for _, idx := range indices {
		end := idx[1]

		// Ensure we do not go out of bounds
		if end < len(text) {
			r := rune(text[end])

			// Only split if next char matches your original intent
			if unicode.IsUpper(r) || r == '"' || r == '\'' {
				result = append(result, text[last:end-1]) // exclude trailing space
				last = end
			}
		}
	}

	// Add remaining text
	if last < len(text) {
		result = append(result, text[last:])
	}

	return result
}
