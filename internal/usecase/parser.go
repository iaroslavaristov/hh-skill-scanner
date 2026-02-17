package usecase

import (
	"fmt"
	"sort"
	"strings"

	"github.com/schollz/progressbar/v3"
	"hh-parser/internal/domain"
)

type Parser struct {
	hh    HHClient
	ai    GeminiClient
	cache Cache
}

func NewParser(hh HHClient, ai GeminiClient, cache Cache) *Parser {
	return &Parser{
		hh:    hh,
		ai:    ai,
		cache: cache,
	}
}

func (p *Parser) Analyze(query string, limit int) ([]domain.Skill, error) {
	vacancies, err := p.hh.SearchVacancies(query, limit)
	if err != nil {
		return nil, err
	}

	if len(vacancies) == 0 {
		return nil, fmt.Errorf("no vacancies found for query: %s", query)
	}

	bar := progressbar.Default(int64(len(vacancies)), "Анализ вакансий")
	skillMap := make(map[string]int)

	var toProcess []string
	var toProcessIDs []string

	for _, v := range vacancies {
		if cachedSkills, found := p.cache.Get(v.ID); found && len(cachedSkills) > 0 {
			p.addSkillsToMap(skillMap, cachedSkills)
			bar.Add(1)
			continue
		}

		desc, err := p.hh.GetFullDescription(v.ID)
		if err != nil || desc == "" {
			bar.Add(1)
			continue
		}

		toProcess = append(toProcess, desc)
		toProcessIDs = append(toProcessIDs, v.ID)

		if len(toProcess) >= 5 {
			p.processBatch(toProcess, toProcessIDs, skillMap, bar)
			toProcess = nil
			toProcessIDs = nil
		}
	}

	if len(toProcess) > 0 {
		p.processBatch(toProcess, toProcessIDs, skillMap, bar)
	}

	p.cache.Save()

	return p.sortSkills(skillMap), nil
}

func (p *Parser) processBatch(descs []string, ids []string, skillMap map[string]int, bar *progressbar.ProgressBar) {
	extracted, err := p.ai.ExtractSkills(descs)
	if err != nil {
		fmt.Printf("\n[Gemini Error]: %v\n", err)
		bar.Add(len(descs))
		return
	}

	if len(extracted) == 0 {
		fmt.Printf("\n[Warning]: Gemini returned 0 skills for batch of %d\n", len(descs))
	}

	for _, s := range extracted {
		cleanSkill := strings.TrimSpace(s)
		if cleanSkill != "" {
			p.addSkillsToMap(skillMap, []string{cleanSkill})
		}
	}
	
	for _, id := range ids {
		p.cache.Set(id, extracted)
		bar.Add(1)
	}
}

func (p *Parser) addSkillsToMap(m map[string]int, skills []string) {
	for _, s := range skills {
		if len(s) < 2 {
			continue
		}
		name := strings.Title(strings.ToLower(s))
		m[name]++
	}
}

func (p *Parser) sortSkills(m map[string]int) []domain.Skill {
	var result []domain.Skill
	for name, count := range m {
		result = append(result, domain.Skill{Name: name, Count: count})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Name < result[j].Name
		}
		return result[i].Count > result[j].Count
	})

	if len(result) > 20 {
		result = result[:20]
	}
	return result
}
