package usecase

import (
	"sort"
	"strings"

	"github.com/schollz/progressbar/v3"
	"hh-parser/internal/domain"
)

type Parser struct {
	hh     HHClient
	ai     GeminiClient
	cache  Cache
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

	bar := progressbar.Default(int64(len(vacancies)), "Анализ вакансий")
	skillMap := make(map[string]int)

	var toProcess []string
	var toProcessIDs []string

	for _, v := range vacancies {
		if cachedSkills, found := p.cache.Get(v.ID); found {
			p.addSkillsToMap(skillMap, cachedSkills)
			bar.Add(1)
			continue
		}

		desc, err := p.hh.GetFullDescription(v.ID)
		if err != nil {
			bar.Add(1)
			continue
		}

		toProcess = append(toProcess, desc)
		toProcessIDs = append(toProcessIDs, v.ID)

		if len(toProcess) >= 20 {
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
		bar.Add(len(descs))
		return
	}

	for i, skillsStr := range extracted {
		skills := strings.Split(skillsStr, ",")
		var cleaned []string
		for _, s := range skills {
			s = strings.TrimSpace(s)
			if s != "" {
				cleaned = append(cleaned, s)
			}
		}
		p.cache.Set(ids[i], cleaned)
		p.addSkillsToMap(skillMap, cleaned)
		bar.Add(1)
	}
}

func (p *Parser) addSkillsToMap(m map[string]int, skills []string) {
	for _, s := range skills {
		m[strings.Title(strings.ToLower(s))]++
	}
}

func (p *Parser) sortSkills(m map[string]int) []domain.Skill {
	var result []domain.Skill
	for name, count := range m {
		result = append(result, domain.Skill{Name: name, Count: count})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	return result
}
