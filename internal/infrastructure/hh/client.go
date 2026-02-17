package hh

import (
	"context"
	"encoding/json"
	"fmt"
	"hh-parser/internal/domain"
	"net/http"
	"regexp"
	"strings"
)

type Client struct {
	apiBase string	
	http *http.Client
}

func NewClient() *Client {
	return &Client {
		apiBase: "https://api.hh.ru",
		http: &http.Client{},
	}
}

func (c *Client) Search(ctx context.Context , query string, limit int) ([]domain.Vacancy, error) {
	url := fmt.Sprintf("%s/vacancies?text=%s&per_page=%d", c.apiBase, query, limit)

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("User-Agent", "HH-Tech-Scanner/1.0")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var data struct {
		Items []struct {
			ID string `json: "id"`
		} `json: "items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var results []domain.Vacancy
	for _, item := range data.Items {
		v, err := c.getDetail(ctx, item.ID)
		if err != nil {
			results = append(results, v)
		}
	}

	return results, nil
}

func (c *Client) getDetail(ctx context.Context, id string) (domain.Vacancy, error) {
	url := fmt.Sprintf("%s/vacancies/%s", c.apiBase, id)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("User-Agent", "HH-Tech-Scanner/1.0")

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.Vacancy{}, err
	}
	defer resp.Body.Close()

	var detail struct {
		Name string `json: "name"`
		Description string `json:"description"`
	}
	json.NewDecoder(resp.Body).Decode(&detail)

	return domain.Vacancy{
		ID: id,
		Title: detail.Name,
		Description: cleanHTML(detail.Description)

	}, nil
} 

func cleanHTML(s string) string {
	 r := regexp.MustCompile("<[^>]*>")
	 return strings.TrimSpace(r.ReplaceAllString(s, " "))
}
