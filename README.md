# Telegram Personal Assistant Bot

Telegram бот на Go с использованием Google Cloud Functions.

## 🚀 Быстрый старт

### 1. Настройка переменных окружения

```bash
# Скопируйте пример файла
cp .env.example .env.local

# Отредактируйте .env.local и добавьте свой токен бота
# Получить токен: https://t.me/BotFather

# Загрузите переменные
source .env.local
```

### 2. Деплоймент

#### Полный деплоймент (первый раз)
```bash
./deploy.sh
```

#### Быстрый деплоймент (для разработки)
```bash
./quick-deploy.sh
```

#### Ручной деплоймент
```bash
gcloud functions deploy telegram-webhook-handler \
  --gen2 \
  --runtime=go124 \
  --region=europe-central2 \
  --source=. \
  --entry-point=TelegramWebhookHandler \
  --trigger-http \
  --allow-unauthenticated \
  --set-env-vars TELEGRAM_BOT_TOKEN="ваш_токен"
```

## 📁 Структура проекта

```
telegram-assistant/
├── function.go           # Основная логика бота
├── go.mod               # Go модуль
├── go.sum               # Зависимости
├── deploy.sh            # Полный скрипт деплоймента
├── quick-deploy.sh      # Быстрый деплоймент
├── .env.example         # Пример переменных окружения
└── README.md           # Этот файл
```

## 🔧 Функциональность

- ✅ Эхо-бот (повторяет сообщения)
- ✅ Обработка текстовых сообщений
- ✅ Автоматическая установка webhook
- ✅ Логирование ошибок

## 🛠 Разработка

### Локальное тестирование

```bash
# Запуск локально (требует functions-framework)
go mod download
export TELEGRAM_BOT_TOKEN="ваш_токен"
functions-framework --target TelegramWebhookHandler
```

### Просмотр логов

```bash
gcloud functions logs read telegram-webhook-handler --region=europe-central2
```

### Удаление функции

```bash
gcloud functions delete telegram-webhook-handler --region=europe-central2
```

## 📋 Требования

- Go 1.24+
- Google Cloud CLI
- Telegram Bot Token

## 🔗 Полезные ссылки

- [Telegram Bot API](https://core.telegram.org/bots/api)
- [Google Cloud Functions](https://cloud.google.com/functions)
- [Functions Framework for Go](https://github.com/GoogleCloudPlatform/functions-framework-go)
