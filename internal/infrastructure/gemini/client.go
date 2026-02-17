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

	model := client.GenerativeModel("gemini-1.5-flash")
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
	promptBuilder.WriteString("Extract IT skills for each job description. Return a JSON array of strings, where each element is a comma-separated list of skills for that job. Example: [\"Go, Docker\", \"Java, Spring\"]. No prose, just JSON.\n")

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

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from gemini")
	}

	rawJSON := fmt.Sprint(resp.Candidates[0].Content.Parts[0])
	rawJSON = strings.TrimPrefix(rawJSON, "```json")
	rawJSON = strings.TrimSuffix(rawJSON, "```")
	rawJSON = strings.TrimSpace(rawJSON)

	var result []string
	if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal gemini response: %w", err)
	}

	return result, nil
}
