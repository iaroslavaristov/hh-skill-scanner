package usecase 

import "context"

type Parser struct {
	provider JobProvider
	analyzer Analyzer
}

func NewParser(p JobProvider, a Analyzer) *Parser {
	return &Parser{provider: p, analyzer: a}
}

func (p *Parser) Run(ctx context.Context, query string, limit int) (map[string]int, error) {
	vacancies, err := p.provider.Search(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int)
	mu := &sync.Mutex{}
	
	jobs := make(chan domain.Vacancy, len(vacancies))
	results := make(chan []string, len(vacancies))


	workerCount := 5 
	var wg sync.WaitGroup

	for w := 1; w <= workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range jobs {
				content := v.Title + " " + v.Description
				skills, _ := p.analyzer.ExtractSkills(ctx, content)
				results <- skills
			}
		}()
	}


	for _, v := range vacancies {
		jobs <- v
	}
	close(jobs)


	go func() {
		wg.Wait()
		close(results)
	}()

	for skills := range results {
		mu.Lock()
		for _, s := range skills {
			stats[s]++
		}
		mu.Unlock()
	}

	return stats, nil
}
