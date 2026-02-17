package usecase

import (
	"context"
	"fmt"
	"hh-parser/internal/domain"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

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

	total := len(vacancies)
	if total == 0 {
		return nil, fmt.Errorf("вакансий по запросу '%s' не найдено", query)
	}

	stats := make(map[string]int)
	batchSize := 20

	bar := progressbar.NewOptions(total,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan][1/3][reset] Анализ навыков..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]track[reset]",
			SaucerHead:    "[green]track[reset]",
			SaucerPadding: " ",
			BarStart:      "|",
			BarEnd:        "|",
		}),
	)

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		currentBatch := vacancies[i:end]
		var textsToProcess []string
		var batchIds []string

		for _, v := range currentBatch {
			if skills, found := p.cache.Get(v.ID); found {
				p.updateStats(stats, skills)
				bar.Add(1)
			} else {
				textsToProcess = append(textsToProcess, v.Title+" "+v.Description)
				batchIds = append(batchIds, v.ID)
			}
		}

		if len(textsToProcess) > 0 {
			batchResults, err := p.analyzer.ExtractSkillsBatch(ctx, textsToProcess)
			if err != nil {
				bar.Add(len(textsToProcess))
				continue
			}

			for idx, skills := range batchResults {
				if idx < len(batchIds) {
					p.cache.Set(batchIds[idx], skills)
					p.updateStats(stats, skills)
				}
				bar.Add(1)
			}
			
			time.Sleep(500 * time.Millisecond)
		}
	}

	fmt.Println()
	return stats, nil
}

func (p *Parser) updateStats(stats map[string]int, skills []string) {
	for _, s := range skills {
		name := strings.ToLower(strings.TrimSpace(s))
		if name != "" && name != "none" {
			stats[name]++
		}
	}
}
