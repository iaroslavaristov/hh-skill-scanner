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
	var mu sync.Mutex
	var wg sync.WaitGroup

	jobs := make(chan domain.Vacancy)

	workerCount := 3 
	
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range jobs {
				var skills []string
				var found bool

				if skills, found = p.cache.Get(v.ID); !found {
					content := v.Title + " " + v.Description
					skills, _ = p.analyzer.ExtractSkills(ctx, content)
					
					p.cache.Set(v.ID, skills)
					
					time.Sleep(2 * time.Second) 
				}

				mu.Lock()
				for _, s := range skills {
					stats[s]++
				}
				mu.Unlock()
			}
		}()
	}

	for _, v := range vacancies {
		jobs <- v
	}
	close(jobs)
	wg.Wait()

	return stats, nil
}
