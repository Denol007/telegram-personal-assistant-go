#!/bin/bash

# Скрипт для полного деплоймента Telegram бота в Google Cloud Functions

set -e  # Остановить выполнение при любой ошибке

# Цвета для красивого вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Конфигурация
FUNCTION_NAME="telegram-webhook-handler"
REGION="europe-central2"
RUNTIME="go124" 
ENTRY_POINT="TelegramWebhookHandler"
SOURCE_DIR="cmd/functions"

echo -e "${BLUE}🚀 Начинаем деплоймент Telegram бота...${NC}"

# Проверки файлов
if [[ ! -f "$SOURCE_DIR/function.go" ]]; then
    echo -e "${RED}❌ Ошибка: файл $SOURCE_DIR/function.go не найден!${NC}"
    echo -e "${YELLOW}Убедитесь, что исходный код находится в папке $SOURCE_DIR/.${NC}"
    exit 1
fi

if [[ ! -f "go.mod" ]]; then
    echo -e "${RED}❌ Ошибка: файл go.mod не найден в корне проекта!${NC}"
    exit 1
fi

echo -e "${YELLOW}📦 Проверяем зависимости...${NC}"
go mod tidy

# --- Проверка переменных окружения ---

# Проверяем токен бота
if [[ -z "$TELEGRAM_BOT_TOKEN" ]]; then
    echo -e "${YELLOW}⚠️  Переменная TELEGRAM_BOT_TOKEN не установлена.${NC}"
    echo -e "${BLUE}Введите токен вашего Telegram бота:${NC}"
    read -s TOKEN
    if [[ -z "$TOKEN" ]]; then
        echo -e "${RED}❌ Токен не может быть пустым!${NC}"
        exit 1
    fi
    TELEGRAM_BOT_TOKEN="$TOKEN"
fi

# Проверяем ID проекта Google Cloud
if [[ -z "$GCP_PROJECT_ID" ]]; then
    echo -e "${YELLOW}⚠️  Переменная GCP_PROJECT_ID не установлена.${NC}"
    echo -e "${BLUE}Введите ID вашего Google Cloud проекта:${NC}"
    read -r PROJECT_ID
    if [[ -z "$PROJECT_ID" ]]; then
        echo -e "${RED}❌ ID проекта не может быть пустым!${NC}"
        exit 1
    fi
    GCP_PROJECT_ID="$PROJECT_ID"
fi


echo -e "${YELLOW}🔨 Деплоим функцию в Google Cloud...${NC}"

# --- Выполняем деплоймент ---
gcloud functions deploy "$FUNCTION_NAME" \
    --gen2 \
    --runtime="$RUNTIME" \
    --region="$REGION" \
    --source=. \
    --entry-point="$ENTRY_POINT" \
    --trigger-http \
    --allow-unauthenticated \
    --set-env-vars GCP_PROJECT_ID=$GCP_PROJECT_ID,TELEGRAM_BOT_TOKEN=$TELEGRAM_BOT_TOKEN \
    --set-build-env-vars GOOGLE_FUNCTION_SOURCE=cmd/functions,GOOGLE_BUILDABLE=./cmd/functions \
    --quiet

# --- Обработка результата ---
if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ Деплоймент завершен успешно!${NC}"
    
    FUNCTION_URL=$(gcloud functions describe "$FUNCTION_NAME" --region="$REGION" --format="value(serviceConfig.uri)")
    
    echo -e "${GREEN}🌐 URL функции: ${FUNCTION_URL}${NC}"
    echo -e "${BLUE}📝 Установка webhook в Telegram...${NC}"
    
    WEBHOOK_RESPONSE=$(curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/setWebhook" -d "url=$FUNCTION_URL")
    
    if echo "$WEBHOOK_RESPONSE" | grep -q '"ok":true'; then
        echo -e "${GREEN}✅ Webhook установлен успешно!${NC}"
    else
        echo -e "${RED}❌ Ошибка при установке webhook:${NC}"
        echo "$WEBHOOK_RESPONSE"
    fi
    
else
    echo -e "${RED}❌ Деплоймент не удался!${NC}"
    exit 1
fi

echo -e "${GREEN}🎉 Готово! Ваш бот готов к работе!${NC}"