package usecase

import (
	"context"
	"hh-parser/internal/domain"
	"sync"
	"time"
)

type Cache interface {
	Get(id string) ([]string, bool)
	Set(id string, skills []string)
}

type Parser struct {
	provider JobProvider
	analyzer Analyzer
	cache    Cache
}

func NewParser(p JobProvider, a Analyzer, c Cache) *Parser {
	return &Parser{provider: p, analyzer: a, cache: c}
}

func (p *Parser) Run(ctx context.Context, query string, limit int) (map[string]int, error) {
	vacancies, err := p.provider.Search(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int)
	batchSize := 5
	
	for i := 0; i < len(vacancies); i += batchSize {
		end := i + batchSize
		if end > len(vacancies) {
			end = len(vacancies)
		}
		
		currentBatch := vacancies[i:end]
		var textsToProcess []string
		var batchIds []string

		for _, v := range currentBatch {
			if skills, found := p.cache.Get(v.ID); found {
				p.updateStats(stats, skills)
			} else {
				textsToProcess = append(textsToProcess, v.Title+" "+v.Description)
				batchIds = append(batchIds, v.ID)
			}
		}

		if len(textsToProcess) > 0 {
			fmt.Printf("Отправляю пакет из %d вакансий в Gemini...\n", len(textsToProcess))

			batchResults, err := p.analyzer.ExtractSkillsBatch(ctx, textsToProcess)
			if err != nil {
				fmt.Printf("Ошибка батча: %v\n", err)
				continue
			}

			for idx, skills := range batchResults {
				p.cache.Set(batchIds[idx], skills)
				p.updateStats(stats, skills)
			}

			time.Sleep(3 * time.Second)
		}
	}

	return stats, nil
}

func (p *Parser) updateStats(stats map[string]int, skills []string) {
	for _, s := range skills {
		stats[strings.ToLower(strings.TrimSpace(s))]++
	}
}
