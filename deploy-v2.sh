#!/bin/bash

# Новый скрипт деплоя (основан на последней успешной команде)
# Использует переменные окружения из .env.local при наличии

set -e

# Цвета
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

# Подхват переменных из .env.local, если файл существует и не был загружен
if [[ -f .env.local ]]; then
  echo -e "${YELLOW}📦 Загружаем .env.local...${NC}"
  # shellcheck disable=SC1091
  source .env.local
fi

# Конфигурация по умолчанию (можно переопределить переменными окружения)
FUNCTION_NAME=${FUNCTION_NAME:-telegram-webhook-handler}
REGION=${GOOGLE_CLOUD_REGION:-europe-central2}
RUNTIME=${GO_RUNTIME:-go124}
ENTRY_POINT=${ENTRY_POINT:-TelegramWebhookHandler}
SOURCE_DIR=${SOURCE_DIR:-.}
SOURCE_ROOT=${SOURCE_ROOT:-.}

# Обязательные переменные
: "${GCP_PROJECT_ID:?GCP_PROJECT_ID не задан}" 
: "${TELEGRAM_BOT_TOKEN:?TELEGRAM_BOT_TOKEN не задан}"

# Проверки
if [[ ! -f "function.go" ]]; then
  echo -e "${RED}❌ Не найден файл function.go${NC}"
  exit 1
fi
if [[ ! -f go.mod ]]; then
  echo -e "${RED}❌ Не найден go.mod в корне проекта${NC}"
  exit 1
fi

# Обновим зависимости
echo -e "${YELLOW}📦 Обновляем зависимости...${NC}"
go mod tidy

# Деплой
echo -e "${BLUE}🚀 Деплой функции ${FUNCTION_NAME} в регион ${REGION}...${NC}"
gcloud functions deploy "${FUNCTION_NAME}" \
  --gen2 \
  --runtime="${RUNTIME}" \
  --region="${REGION}" \
  --source="${SOURCE_ROOT}" \
  --entry-point="${ENTRY_POINT}" \
  --trigger-http \
  --allow-unauthenticated \
  --set-env-vars "GCP_PROJECT_ID=${GCP_PROJECT_ID},TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}"

# Вывести URL
echo -e "${GREEN}✅ Деплой завершён.${NC}"
FUNC_URL=$(gcloud functions describe "${FUNCTION_NAME}" --region="${REGION}" --format="value(serviceConfig.uri)" 2>/dev/null || true)
if [[ -n "${FUNC_URL}" ]]; then
  echo -e "${BLUE}🌐 URL: ${FUNC_URL}${NC}"
fi
