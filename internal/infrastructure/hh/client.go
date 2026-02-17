package hh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"hh-parser/internal/domain"
)

type Client struct {
	httpClient *http.Client
}

type hhSearchResponse struct {
	Items []struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"items"`
}

type hhDetailResponse struct {
	Description string `json:"description"`
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *Client) SearchVacancies(query string, limit int) ([]domain.Vacancy, error) {
	url := fmt.Sprintf("https://api.hh.ru/vacancies?text=%s&per_page=%d", query, limit)
	
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "HH-Skill-Scanner/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResp hhSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	var vacancies []domain.Vacancy
	for _, item := range searchResp.Items {
		vacancies = append(vacancies, domain.Vacancy{
			ID:    item.ID,
			Title: item.Name,
			URL:   item.URL,
		})
	}

	return vacancies, nil
}

func (c *Client) GetFullDescription(vacancyID string) (string, error) {
	url := fmt.Sprintf("https://api.hh.ru/vacancies/%s", vacancyID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "HH-Skill-Scanner/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var detailResp hhDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&detailResp); err != nil {
		return "", err
	}

	return detailResp.Description, nil
}
