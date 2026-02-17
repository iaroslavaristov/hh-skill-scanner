package domain

type Vacancy struct {
	ID string
	Title string
	Description string
	Skills []string 
}

type Result struct {
	TechName string
	Count int
}
