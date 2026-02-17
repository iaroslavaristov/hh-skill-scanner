package usecase

import "hh-parser/internal/domain"

type HHClient interface {
	SearchVacancies(query string, limit int) ([]domain.Vacancy, error)
	GetFullDescription(vacancyID string) (string, error)
}

type GeminiClient interface {
	ExtractSkills(descriptions []string) ([]string, error)
}

type Cache interface {
	Get(vacancyID string) ([]string, bool)
	Set(vacancyID string, skills []string) error
	Save() error
}
