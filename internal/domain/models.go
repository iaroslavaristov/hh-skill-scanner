package domain

type Vacancy struct {
	ID          string
	Title       string
	Description string
	URL         string
}

type Skill struct {
	Name  string
	Count int
}

type Config struct {
	GeminiAPIKey string
	SkillsFile   string
}
