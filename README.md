# Telegram Personal Assistant Bot

A personal assistant Telegram bot written in Go using Google Cloud Functions and Firestore for note storage.

## ✨ Features

- 📝 **Note saving** - send any text to the bot and it will save it as a note
- 📋 **View notes** - `/list` command shows all your notes
- 🔒 **Personal data** - each user only sees their own notes
- ☁️ **Cloud-native** - runs on Google Cloud Functions + Firestore
- 🚀 **Production-ready** - with logging and error handlingPersonal Assistant Bot

Персональный помощник на базе Telegram бота, написанный на Go с использованием Google Cloud Functions и Firestore для хранения заметок.

## ✨ Возможности

- 📝 **Сохранение заметок** - отправляйте боту любой текст и он сохранит его как заметку
- � **Просмотр заметок** - команда `/list` покажет все ваши заметки
- 🔒 **Персональные данные** - каждый пользователь видит только свои заметки
- ☁️ **Cloud-native** - работает на Google Cloud Functions + Firestore
- �🚀 **Готов к продакшену** - с логированием и обработкой ошибок

## 🚀 Quick Start

### 1. Environment Variables Setup

```bash
# Copy the example file
cp .env.example .env.local

# Edit .env.local and add your tokens
# TELEGRAM_BOT_TOKEN - get from: https://t.me/BotFather
# GCP_PROJECT_ID - your Google Cloud project ID

# Load variables
source .env.local
```

### 2. Deployment

#### Automatic deployment (recommended)
```bash
./deploy-v2.sh
```

#### Manual deployment
```bash
gcloud functions deploy telegram-webhook-handler \
  --gen2 \
  --runtime=go124 \
  --region=europe-central2 \
  --source=. \
  --entry-point=TelegramWebhookHandler \
  --trigger-http \
  --allow-unauthenticated \
  --set-env-vars TELEGRAM_BOT_TOKEN="your_token",GCP_PROJECT_ID="your_project"
```

## 📁 Project Structure

```
telegram-assistant/
├── function.go              # Cloud Function entry point
├── go.mod                   # Go module and dependencies
├── go.sum                   # Dependency lock file
├── deploy-v2.sh            # Automatic deployment script
├── .env.example            # Environment variables example
├── README.md               # Project documentation
└── internal/
    ├── bot/
    │   ├── handler.go      # Main message handling logic
    │   └── commands.go     # Bot command handlers
    ├── config/
    │   └── config.go       # Configuration loading
    ├── note/
    │   └── note.go         # Note model
    ├── store/
    │   └── firestore.go    # Firestore operations
    └── telegram/
        └── telegram.go     # Telegram API client
```

## 🤖 Bot Commands

- **Any text** - saved as a note
- **/list** - show all your notes

## 🔧 Technologies

- **Backend**: Go 1.20+
- **Cloud Platform**: Google Cloud Functions (Gen 2)
- **Database**: Google Firestore
- **API**: Telegram Bot API
- **Infrastructure**: Google Cloud CLI

## 🛠 Development

### Install Dependencies

```bash
go mod download
```

### Local Testing

```bash
# Install Functions Framework
go install github.com/GoogleCloudPlatform/functions-framework-go/funcframework@latest

# Set environment variables
export TELEGRAM_BOT_TOKEN="your_token"
export GCP_PROJECT_ID="your_project"

# Run locally
functions-framework --target TelegramWebhookHandler
```

### Code Verification

```bash
# Check compilation
go build -o /dev/null .

# Format code
go fmt ./...

# Run linter (requires golangci-lint)
golangci-lint run
```

### View Logs

```bash
# Cloud Function logs
gcloud functions logs read telegram-webhook-handler --region=europe-central2

# Real-time logs
gcloud functions logs tail telegram-webhook-handler --region=europe-central2
```

### Function Management

```bash
# Function information
gcloud functions describe telegram-webhook-handler --region=europe-central2

# Delete function
gcloud functions delete telegram-webhook-handler --region=europe-central2
```

## 📋 Requirements

- **Go**: 1.20+
- **Google Cloud CLI**: latest version
- **Google Cloud Project**: with enabled APIs
  - Cloud Functions API
  - Cloud Firestore API
- **Telegram Bot Token**: from [@BotFather](https://t.me/BotFather)

## 🔒 Environment Variables

```bash
# Required
TELEGRAM_BOT_TOKEN=your_bot_token_here
GCP_PROJECT_ID=your-gcp-project-id

# Optional (for deployment customization)
FUNCTION_NAME=telegram-webhook-handler
GOOGLE_CLOUD_REGION=europe-central2
GO_RUNTIME=go124
ENTRY_POINT=TelegramWebhookHandler
```

## 🔗 Useful Links

- [Telegram Bot API](https://core.telegram.org/bots/api)
- [Google Cloud Functions](https://cloud.google.com/functions/docs)
- [Google Cloud Firestore](https://cloud.google.com/firestore/docs)
- [Functions Framework for Go](https://github.com/GoogleCloudPlatform/functions-framework-go)
- [Go Cloud Development Kit](https://gocloud.dev/)

## 📝 License

MIT License - see LICENSE file for details.

## 🤝 Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
