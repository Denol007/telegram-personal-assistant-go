# Telegram Personal Assistant Bot

[![CI/CD Pipeline](https://github.com/Denol007/telegram-personal-assistant-go/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/Denol007/telegram-personal-assistant-go/actions/workflows/ci-cd.yml)
[![Code Quality](https://github.com/Denol007/telegram-personal-assistant-go/actions/workflows/quality.yml/badge.svg)](https://github.com/Denol007/telegram-personal-assistant-go/actions/workflows/quality.yml)
[![Security](https://github.com/Denol007/telegram-personal-assistant-go/actions/workflows/security.yml/badge.svg)](https://github.com/Denol007/telegram-personal-assistant-go/actions/workflows/security.yml)

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

#### CI/CD Deployment (recommended)
The project includes automated CI/CD pipeline using GitHub Actions:

1. **Automatic deployment**: Push to `main` branch triggers automatic deployment to Google Cloud Functions
2. **Quality checks**: Every PR runs tests, linting, and security scans
3. **Required secrets** in GitHub repository settings:
   - `GCP_SA_KEY` - Google Cloud Service Account JSON key
   - `GCP_PROJECT_ID` - Your Google Cloud project ID
   - `TELEGRAM_BOT_TOKEN` - Your Telegram bot token

To set up CI/CD:
1. Create a Google Cloud Service Account with Cloud Functions permissions
2. Download the JSON key file
3. Add the secrets to your GitHub repository settings
4. Push to main branch - deployment will happen automatically!

#### Manual deployment
```bash
./deploy-v2.sh
```

#### Manual deployment (detailed)
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

### Code Quality

```bash
# Check compilation
go build -o /dev/null .

# Format code
go fmt ./...

# Run linter (requires golangci-lint)
golangci-lint run

# Run tests
go test ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### CI/CD Pipeline

The project uses GitHub Actions for automated CI/CD:

- **CI Pipeline**: Runs on every push and PR
  - Code formatting checks
  - Linting with golangci-lint
  - Build verification
  - Test execution
  - Security scanning

- **CD Pipeline**: Runs on main branch pushes
  - Deploys to Google Cloud Functions
  - Updates function environment variables
  - Outputs deployment URL

- **Quality Pipeline**: 
  - Code coverage reporting
  - Dependency vulnerability scanning
  - Automated dependency updates (weekly)

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

### Runtime Environment
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

### CI/CD Secrets (GitHub Repository Settings)
```bash
# Required for automated deployment
GCP_SA_KEY={"type":"service_account",...}  # Google Cloud Service Account JSON
GCP_PROJECT_ID=your-gcp-project-id         # Same as runtime
TELEGRAM_BOT_TOKEN=your_bot_token_here      # Same as runtime
```

### Setting up CI/CD Secrets

1. **Google Cloud Service Account**:
   ```bash
   # Create service account
   gcloud iam service-accounts create github-actions \
     --description="GitHub Actions deployment" \
     --display-name="GitHub Actions"

   # Grant necessary permissions
   gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
     --member="serviceAccount:github-actions@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
     --role="roles/cloudfunctions.developer"

   gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
     --member="serviceAccount:github-actions@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
     --role="roles/iam.serviceAccountUser"

   # Create and download key
   gcloud iam service-accounts keys create key.json \
     --iam-account=github-actions@YOUR_PROJECT_ID.iam.gserviceaccount.com
   ```

2. **Add to GitHub Secrets**:
   - Go to your repository → Settings → Secrets and variables → Actions
   - Add `GCP_SA_KEY` with the content of `key.json`
   - Add `GCP_PROJECT_ID` with your project ID
   - Add `TELEGRAM_BOT_TOKEN` with your bot token

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
