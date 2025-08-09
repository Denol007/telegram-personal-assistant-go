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

# Проверяем, что мы в правильной директории
if [[ ! -f "function.go" ]]; then
    echo -e "${RED}❌ Ошибка: файл function.go не найден!${NC}"
    echo -e "${YELLOW}Убедитесь, что вы запускаете скрипт из директории проекта.${NC}"
    exit 1
fi

# Проверяем, что go.mod существует
if [[ ! -f "go.mod" ]]; then
    echo -e "${RED}❌ Ошибка: файл go.mod не найден!${NC}"
    exit 1
fi

echo -e "${YELLOW}📦 Проверяем зависимости...${NC}"
go mod tidy

# Проверяем переменную окружения TELEGRAM_BOT_TOKEN
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

echo -e "${YELLOW}🔨 Деплоим функцию в Google Cloud...${NC}"

# Выполняем деплоймент
gcloud functions deploy "$FUNCTION_NAME" \
    --gen2 \
    --runtime="$RUNTIME" \
    --region="$REGION" \
    --source=. \
    --entry-point="$ENTRY_POINT" \
    --trigger-http \
    --allow-unauthenticated \
    --set-env-vars TELEGRAM_BOT_TOKEN="$TELEGRAM_BOT_TOKEN" \
    --quiet

if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ Деплоймент завершен успешно!${NC}"
    
    # Получаем URL функции
    FUNCTION_URL=$(gcloud functions describe "$FUNCTION_NAME" --region="$REGION" --format="value(serviceConfig.uri)")
    
    echo -e "${GREEN}🌐 URL функции: ${FUNCTION_URL}${NC}"
    echo -e "${BLUE}📝 Не забудьте установить webhook в Telegram:${NC}"
    echo -e "${YELLOW}curl -X POST \"https://api.telegram.org/bot\$TELEGRAM_BOT_TOKEN/setWebhook\" -d \"url=\$FUNCTION_URL\"${NC}"
    
    # Предлагаем автоматически установить webhook
    echo -e "${BLUE}Хотите автоматически установить webhook? (y/n):${NC}"
    read -r SETUP_WEBHOOK
    
    if [[ "$SETUP_WEBHOOK" =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}🔗 Устанавливаем webhook...${NC}"
        WEBHOOK_RESPONSE=$(curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/setWebhook" -d "url=$FUNCTION_URL")
        
        if echo "$WEBHOOK_RESPONSE" | grep -q '"ok":true'; then
            echo -e "${GREEN}✅ Webhook установлен успешно!${NC}"
        else
            echo -e "${RED}❌ Ошибка при установке webhook:${NC}"
            echo "$WEBHOOK_RESPONSE"
        fi
    fi
    
else
    echo -e "${RED}❌ Деплоймент не удался!${NC}"
    exit 1
fi

echo -e "${GREEN}🎉 Готово! Ваш бот готов к работе!${NC}"
