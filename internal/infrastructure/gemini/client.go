package gemini

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
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

func (c *Client) ExtractSkills(ctx context.Context, text string) ([]string, error) {
	prompt := fmt.Sprintf(`
		Проанализируй текст вакансии ниже. 
		Извлеки из него список IT-технологий, языков программирования и инструментов.
		Верни ответ строго в виде списка через запятую. 
		Если технологий не найдено, напиши "none".
		Текст: %s`, text)

	resp, err := c.model.GenerateContext(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) == 0 {
		return nil, nil
	}

	content := fmt.Sprint(resp.Candidates[0].Content.Parts[0])
	if strings.ToLower(strings.TrimSpace(content)) == "none" {
		return nil, nil
	}

	parts := strings.Split(content, ",")
	var skills []string
	for _, p := range parts {
		trimmed := strings.ToLower(strings.TrimSpace(p))
		if trimmed != "" {
			skills = append(skills, trimmed)
		}
	}

	return skills, nil
}
