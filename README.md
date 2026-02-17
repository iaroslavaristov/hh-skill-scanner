# HH Skill Scanner (Gemini 2.5 Powered)

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![AI](https://img.shields.io/badge/AI-Gemini%203.0%20Flash-orange?style=flat)](https://ai.google.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

Мощная CLI-утилита на **Go**, которая анализирует вакансии с HeadHunter и с помощью нейросети **Gemini 2.5 Flash Lite** извлекает из них актуальный список требуемых IT-технологий и навыков.



## Особенности

* **AI-Анализ**: Использует возможности **Gemini 2.5 Flash Lite** для глубокого понимания текста вакансий.
* **Batch Processing**: Группировка вакансий (до 20 штук) в одном запросе к AI, что обеспечивает высокую скорость и экономию лимитов API.
* **Умное кэширование**: Результаты сохраняются в локальный `skills_cache.json`. Повторный анализ тех же вакансий происходит мгновенно без обращения к AI.
* **Clean Architecture**: Код разделен на слои (Domain, Usecase, Infrastructure), что делает его легко тестируемым и расширяемым.
* **Интерактивный UI**: Встроенный прогресс-бар позволяет наглядно отслеживать статус обработки в реальном времени.

## Стек технологий

* **Язык**: Go (Golang)
* **AI SDK**: Google Generative AI SDK
* **API**: HeadHunter API
* **UI**: schollz/progressbar
* **Инструменты**: Makefile, JSON Storage

## Структура проекта

```text
├── cmd/
│   └── hh-cli/        # Точка входа (CLI интерфейс)
├── internal/
│   ├── domain/        # Бизнес-сущности и модели данных
│   ├── usecase/       # Оркестрация логики (анализ, батчинг, статистика)
│   └── infrastructure/# Реализация адаптеров (HH Client, Gemini Client, File Cache)
├── Makefile           # Автоматизация сборки и запуска
├── .env               # Конфигурация секретов (API ключи)
└── skills_cache.json  # Локальная база данных навыков
```


## Быстрый старт

### 1. Настройка окружения
Создайте файл `.env` в корне проекта и добавьте ваш API ключ от Google AI Studio (Gemini):
```env
GEMINI_API_KEY=your_api_key_here
```
### 2. Сборка проекта
```bash
make build
```
### 3. Запуск анализа
```bash
# -query - запрос, -limit - лимит вакансий
./scanner -query="Golang developer" -limit=10
```

## Пример работы
При запуске программы вы увидите интерактивный прогресс-бар:
[1/3] Анализ навыков... |██████████████░| 84%

После завершения сканер выведет итоговую таблицу:
```Plaintext
Топ навыков по запросу "Golang developer":
1. Go (Golang)      - 48
2. PostgreSQL       - 32
3. Docker           - 29
4. Kubernetes       - 21
5. gRPC             - 18
```

## Лицензия
Данный проект распространяется под лицензией MIT. Подробности в файле [LICENSE](./LICENSE).
