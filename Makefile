BINARY_NAME=scanner
MAIN_PATH=./cmd/hh-cli/main.go

.PHONY: build run clean test help

build:
	@echo "Сборка бинарного файла..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Готово! Запускай командой: ./$(BINARY_NAME)"

run:
	go run $(MAIN_PATH) -query="$(q)" -limit=$(l)

clean:
	@echo "🧹 Очистка..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe
	@echo "✨ Чисто."

fmt:
	go fmt ./...

help:
	@echo "Доступные команды:"
	@echo "  make build   - Скомпилировать проект в файл $(BINARY_NAME)"
	@echo "  make run q=X l=Y - Запустить проект с запросом X и лимитом Y"
	@echo "  make fmt     - Отформатировать весь код"
	@echo "  make clean   - Удалить бинарные файлы"
