package usecase

import (
	"hh-parser/internal/domain"
	"testing"
)

// --- MOCKS ---

type MockHH struct{}

func (m *MockHH) SearchVacancies(query string, limit int) ([]domain.Vacancy, error) {
	return []domain.Vacancy{
		{ID: "vac_1", Title: "Golang Developer"},
		{ID: "vac_2", Title: "Go Engineer"},
	}, nil
}

func (m *MockHH) GetFullDescription(id string) (string, error) {
	return "We need Go and Docker experience", nil
}

type MockAI struct{}

func (m *MockAI) ExtractSkills(descriptions []string) ([]string, error) {
	return []string{"Go, Docker", "Go, Kubernetes"}, nil
}

type MockCache struct{}

func (m *MockCache) Get(id string) ([]string, bool) {
	return nil, false
}

func (m *MockCache) Set(id string, skills []string) error {
	return nil
}

func (m *MockCache) Save() error {
	return nil
}


func TestParser_Analyze(t *testing.T) {
	hh := &MockHH{}
	ai := &MockAI{}
	cache := &MockCache{}

	parser := NewParser(hh, ai, cache)

	limit := 2
	results, err := parser.Analyze("golang", limit)

	if err != nil {
		t.Fatalf("Analyze вернул ошибку: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("Ожидались навыки, но список пуст")
	}

	var goCount int
	for _, s := range results {
		if s.Name == "Go" {
			goCount = s.Count
		}
	}

	if goCount != 2 {
		t.Errorf("Ожидалось 2 упоминания Go, получено %d", goCount)
	}
}
