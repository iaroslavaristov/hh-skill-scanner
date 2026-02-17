package usecase

import (
	"context"
	"hh-parser/internal/domain"
)

type JobProvider interface {
	Search(ctx context.Context, query string, limit int) ([]domain.Vacancy, error)
}

type Analyzer interface {
	ExtractSkills(ctx context.Context, text string) ([]string, error)
}
