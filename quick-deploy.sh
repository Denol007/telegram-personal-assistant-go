#!/bin/bash

# Скрипт для быстрого деплоймента во время разработки
# Более простая версия без интерактивных вопросов

set -e

# Цвета
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Конфигурация
FUNCTION_NAME="telegram-webhook-handler"
REGION="europe-central2"
RUNTIME="go124"
ENTRY_POINT="TelegramWebhookHandler"

echo -e "${BLUE}🔄 Быстрый деплоймент...${NC}"

# Проверяем файлы
if [[ ! -f "function.go" ]]; then
    echo "❌ function.go не найден!"
    exit 1
fi

# Проверяем токен
if [[ -z "$TELEGRAM_BOT_TOKEN" ]]; then
    echo "❌ Установите TELEGRAM_BOT_TOKEN:"
    echo "export TELEGRAM_BOT_TOKEN='ваш_токен'"
    exit 1
fi

# Проверяем ID проекта
if [[ -z "$GCP_PROJECT_ID" ]]; then
    echo "❌ Установите GCP_PROJECT_ID:"
    echo "export GCP_PROJECT_ID='ваш_проект'"
    exit 1
fi

echo -e "${YELLOW}📦 Обновляем зависимости...${NC}"
go mod tidy

echo -e "${YELLOW}🚀 Деплоим...${NC}"

# Быстрый деплоймент
gcloud functions deploy "$FUNCTION_NAME" \
    --gen2 \
    --runtime="$RUNTIME" \
    --region="$REGION" \
    --source=. \
    --entry-point="$ENTRY_POINT" \
    --trigger-http \
    --allow-unauthenticated \
    --set-env-vars TELEGRAM_BOT_TOKEN="$TELEGRAM_BOT_TOKEN",GCP_PROJECT_ID="$GCP_PROJECT_ID" \
    --quiet

echo -e "${GREEN}✅ Готово!${NC}"

# Показываем URL
FUNCTION_URL=$(gcloud functions describe "$FUNCTION_NAME" --region="$REGION" --format="value(serviceConfig.uri)" 2>/dev/null || echo "")
if [[ -n "$FUNCTION_URL" ]]; then
    echo -e "${BLUE}🌐 URL: ${FUNCTION_URL}${NC}"
fi
