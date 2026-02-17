package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Client struct {
	genaiClient *genai.Client
	model       *genai.GenerativeModel
}

func NewClient(ctx context.Context, apiKey string) (*Client, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	model := client.GenerativeModel("gemini-3-flash-preview")
	model.SetTemperature(0.1)
	model.ResponseMIMEType = "application/json"

	return &Client{
		genaiClient: client,
		model:       model,
	}, nil
}

func (c *Client) ExtractSkills(texts []string) ([]string, error) {
	ctx := context.Background()
	var promptBuilder strings.Builder
	promptBuilder.WriteString("Extract all technical IT skills as a flat JSON array of strings. No prose.\n")

	for i, t := range texts {
		// РАДИКАЛЬНАЯ ОЧИСТКА
		cleanText := fixUtf8(t)

		if len(cleanText) > 2000 {
			cleanText = cleanText[:2000]
		}
		promptBuilder.WriteString(fmt.Sprintf("Job %d: %s\n", i+1, cleanText))
	}

	finalPrompt := fixUtf8(promptBuilder.String())

	resp, err := c.model.GenerateContent(ctx, genai.Text(finalPrompt))
	if err != nil {
		return nil, fmt.Errorf("gemini generation error: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response")
	}

	rawJSON := fmt.Sprint(resp.Candidates[0].Content.Parts[0])
	rawJSON = strings.TrimSpace(rawJSON)
	rawJSON = strings.TrimPrefix(rawJSON, "```json")
	rawJSON = strings.TrimPrefix(rawJSON, "```")
	rawJSON = strings.TrimSuffix(rawJSON, "```")
	rawJSON = strings.TrimSpace(rawJSON)

	var result []string
	if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
		return nil, fmt.Errorf("json error: %w | raw: %s", err, rawJSON)
	}

	return result, nil
}

func fixUtf8(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	v := make([]rune, 0, len(s))
	for i, r := range s {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(s[i:])
			if size == 1 {
				continue
			}
		}
		v = append(v, r)
	}
	return string(v)
}
