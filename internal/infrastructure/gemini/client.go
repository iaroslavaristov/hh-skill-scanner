package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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
	if len(texts) == 0 {
		return []string{}, nil
	}

	ctx := context.Background()
	var promptBuilder strings.Builder

	promptBuilder.WriteString("Extract all technical terms, frameworks, languages, and software tools from these job descriptions. Return them as a flat JSON array of strings. Example: [\"1C\", \"SQL\", \"SKD\", \"ERP\"]. Focus on professional keywords. Return ONLY the JSON array, no other text.\n")

	for i, t := range texts {
		if len(t) > 2000 {
			t = t[:2000]
		}
		promptBuilder.WriteString(fmt.Sprintf("Text %d: %s\n", i+1, t))
	}

	resp, err := c.model.GenerateContent(ctx, genai.Text(promptBuilder.String()))
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("gemini returned no data")
	}

	rawJSON := fmt.Sprint(resp.Candidates[0].Content.Parts[0])
	
	rawJSON = strings.TrimSpace(rawJSON)
	rawJSON = strings.TrimPrefix(rawJSON, "```json")
	rawJSON = strings.TrimPrefix(rawJSON, "```")
	rawJSON = strings.TrimSuffix(rawJSON, "```")
	rawJSON = strings.TrimSpace(rawJSON)

	var result []string
	if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
		fmt.Printf("\n[DEBUG] Ошибка парсинга. Ответ ИИ: %s\n", rawJSON)
		return nil, err
	}

	return result, nil
}
