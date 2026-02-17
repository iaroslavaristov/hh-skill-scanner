package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"hh-parser/internal/infrastructure/gemini"
	"hh-parser/internal/infrastructure/hh"
	"hh-parser/internal/usecase"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("Ошибка: GEMINI_API_KEY не найден в .env")
	}

	queryFlag := flag.String("query", "", "Название вакансии (например, 'Frontend developer')")
	limitFlag := flag.Int("limit", 10, "Количество вакансий для анализа")
	flag.Parse()

	searchQuery := *queryFlag

	if searchQuery == "" {
		fmt.Print("Введите название позиции для анализа: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		searchQuery = strings.TrimSpace(input)
	}

	if searchQuery == "" {
		log.Fatal("Ошибка: поисковый запрос не может быть пустым")
	}

	ctx := context.Background()

	hhClient := hh.NewClient()
	geminiClient, err := gemini.NewClient(ctx, apiKey)
	if err != nil {
		log.Fatalf("Ошибка Gemini: %v", err)
	}
	fileCache := cache.NewFileCache("skills_cache.json")

	parser := usecase.NewParser(hhClient, geminiClient, fileCache)

	fmt.Printf("\nАнализирование топ-%d вакансий по запросу: '%s'...\n", *limitFlag, searchQuery)

	stats, err := parser.Run(ctx, searchQuery, *limitFlag)
	if err != nil {
		log.Fatalf("Ошибка при выполнении: %v", err)
	}

	printFinalStats(stats)
}

func printFinalStats(stats map[string]int) {
	if len(stats) == 0 {
		fmt.Println("Технологии не найдены.")
		return
	}

	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range stats {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	fmt.Println("\nАнализ закончен. Технологии:")
	fmt.Println("--------------------------------------------------")
	for i, entry := range ss {
		if i >= 20 { break }
		fmt.Printf("%2d) %-15s — встретилось %d раз(а)\n", i+1, entry.Key, entry.Value)
	}
}
