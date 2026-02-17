package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

type Client struct {
	genaiClient *genai.Client
	model *genai.GenerativeModel
}

func NewClient(ctx context.Context, apiKey string) (*Client, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}


	model := client.GenerativeModel("gemini-3.0-fast")
	
	model.SetTemperature(0.1)
	model.ResponseMIMEType = "application/json"
	return &Client{
		genaiClient: client,
		model:       model,
	}, nil
}

func (c *Client) ExtractSkillsBatch(ctx context.Context, texts []string) ([][]string, error) {
	var promptBuilder strings.Builder
	promptBuilder.WriteString("Extract IT skills as a JSON array of string arrays. No talk, just JSON. Format: [[\"skill1\"], [\"skill2\"]]\n")

	for i, t := range texts {
		if len(t) > 3000 {
			t = t[:3000]
		}
		promptBuilder.WriteString(fmt.Sprintf("Job %d: %s\n", i+1, t))
	}
	resp, err := c.model.GenerateContent(ctx, genai.Text(promptBuilder.String()))
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) == 0 {
		return nil, nil
	}

	rawJSON := fmt.Sprint(resp.Candidates[0].Content.Parts[0])
	rawJSON = strings.TrimPrefix(rawJSON, "```json")
	rawJSON = strings.TrimSuffix(rawJSON, "```")
	rawJSON = strings.TrimSpace(rawJSON)

	var result [][]string
	if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON от Gemini: %w", err)
	}

	return result, nil
}
