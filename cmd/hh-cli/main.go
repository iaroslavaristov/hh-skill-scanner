package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"hh-parser/internal/infrastructure/cache"
	"hh-parser/internal/infrastructure/gemini"
	"hh-parser/internal/infrastructure/hh"
	"hh-parser/internal/usecase"
)

func main() {
	query := flag.String("query", "Golang developer", "Search query")
	limit := flag.Int("limit", 20, "Number of vacancies to analyze")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY must be set")
	}

	ctx := context.Background()

	hhClient := hh.NewClient()
	
	geminiClient, err := gemini.NewClient(ctx, apiKey)
	if err != nil {
		log.Fatalf("Failed to init Gemini: %v", err)
	}

	fileCache := cache.NewFileCache("skills_cache.json")

	parser := usecase.NewParser(hhClient, geminiClient, fileCache)

	fmt.Printf("Searching for %d vacancies: %s\n", *limit, *query)

	skills, err := parser.Analyze(*query, *limit)
	if err != nil {
		log.Fatalf("Analysis failed: %v", err)
	}

	fmt.Printf("\nTop Skills for %s:\n", *query)
	for i, skill := range skills {
		if i >= 15 {
			break
		}
		fmt.Printf("%d. %s - %d\n", i+1, skill.Name, skill.Count)
	}
}
