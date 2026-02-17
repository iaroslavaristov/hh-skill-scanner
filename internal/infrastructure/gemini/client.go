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

func NewClient(cxt context.Context, apiKey string) (*Client, error) {
	client, err := genai.NewClient(ctx, option.With.APIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	model := client.GenerativeModel("gemini-3.0-fast")
	model.SetTemperature(0.2)

	return &Client{
		genaiClient: client,
		model: model,
	}, nil
}

func (c *Client) ExtractSkillsBatch(ctx context.Context, texts []string) ([][]string, error) {
	var promptBuilder strings.Builder
	promptBuilder.WriteString("Проанализируй список вакансий и извлеки IT-технологии для каждой. " +
		"Верни ответ СТРОГО в формате JSON: массив массивов строк, где каждый внутренний массив — это навыки одной вакансии. " +
		"Пример: [[\"go\", \"docker\"], [\"react\", \"js\"]].\n\n")

	for i, t := range texts {
		promptBuilder.WriteString(fmt.Sprintf("Вакансия %d: %s\n\n", i+1, t))
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
