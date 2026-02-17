package main 

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"hh-parser/internal/infrastructure/gemini"
	"hh-parser/internal/infrastructure/hh"
	"hh-parser/internal/usecase"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set in .env file")
	}

	query := flag.String("query", "Golang developer", "Вакансия для поиска")
	limit := flag.Int("limit", 5, "Количество вакансий для анализа")
	flag.Parse()
	ctx := context.Background()

	// 3. Инициализация адаптеровол
	hhClient := hh.NewClient()
	geminiClient, err := gemini.NewClient(ctx, apiKey)
	if err != nil {
		log.Fatalf("Gemini init error: %v", err)
	}

	parser := usecase.NewParser(hhClient, geminiClient)

	fmt.Printf("Анализ вакансий: %s...\n", *query)

	stats, err := parser.Run(ctx, *query, *limit)
	if err != nil {
		log.Fatalf("Execution error: %v", err)
	}

	printResults(stats)
}

func printResults(stats map[string]int) {
	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range stats {
		sorted = append(sorted, kv{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	fmt.Println("\n📊 Популярность технологий:")
	for _, entry := range sorted {
		fmt.Printf("%-15s: %d\n", entry.Key, entry.Value)
	}
}
}
