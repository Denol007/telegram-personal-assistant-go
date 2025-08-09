#!/bin/bash

# Скрипт для быстрого деплоймента Telegram бота в Google Cloud Functions

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

echo -e "${BLUE}🚀 Начинаем деплоймент Telegram бота...${NC}"

# Проверки файлов
if [[ ! -f "function.go" || ! -f "go.mod" ]]; then
    echo -e "${RED}❌ Ошибка: файлы function.go и/или go.mod не найдены!${NC}"
    echo -e "${YELLOW}Убедитесь, что вы запускаете скрипт из корневой директории проекта.${NC}"
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

# НОВОЕ: Проверяем ID проекта Google Cloud
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
# ИСПРАВЛЕНО: Добавлена переменная GCP_PROJECT_ID
gcloud functions deploy "$FUNCTION_NAME" \
    --gen2 \
    --runtime="$RUNTIME" \
    --region="$REGION" \
    --source=. \
    --entry-point="$ENTRY_POINT" \
    --trigger-http \
    --allow-unauthenticated \
    --set-env-vars GCP_PROJECT_ID=$GCP_PROJECT_ID,TELEGRAM_BOT_TOKEN=$TELEGRAM_BOT_TOKEN \
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
