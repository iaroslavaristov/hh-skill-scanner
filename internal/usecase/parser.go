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
	stats := make(map[string]int)

	vacancies, err := p.provider.Search(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	for _, v := range vacancies {
		content := v.Title + " " + v.Description

		skills, err := p.analyzer.ExtractSkills(ctx, content)
		if err != nil {
			continue
		}

		for _, s := range skills {
			stats[s]++
		}
	}

	return stats, nil
}
