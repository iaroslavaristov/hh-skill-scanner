package hh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"hh-parser/internal/domain"
)

type Client struct {
	httpClient *http.Client
}

type hhSearchResponse struct {
	Items []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"items"`
}

type hhDetailResponse struct {
	Description string `json:"description"`
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}
}

func (c *Client) SearchVacancies(query string, limit int) ([]domain.Vacancy, error) {
	u, err := url.Parse("https://api.hh.ru/vacancies")
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("text", query)
	params.Add("per_page", fmt.Sprintf("%d", limit))
	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HH search error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HH search status: %d", resp.StatusCode)
	}

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
	uri := fmt.Sprintf("https://api.hh.ru/vacancies/%s", vacancyID)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HH detail error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HH detail status: %d", resp.StatusCode)
	}

	var detailResp hhDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&detailResp); err != nil {
		return "", err
	}

	return detailResp.Description, nil
}
