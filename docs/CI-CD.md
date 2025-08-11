# CI/CD Pipeline Documentation

## Overview
This repository uses GitHub Actions for automated CI/CD pipeline that:

1. **Tests and builds** code on every push and pull request
2. **Deploys automatically** to Google Cloud Functions on main branch pushes
3. **Monitors code quality** with linting and security scans
4. **Manages dependencies** with automated updates

## Workflows

### 1. Main CI/CD Pipeline (`.github/workflows/ci-cd.yml`)
- **Triggers**: Push to main/master, Pull Requests
- **Jobs**:
  - **Test**: Build verification, tests, go vet, formatting checks
  - **Deploy**: Automated deployment to Google Cloud Functions (main branch only)

### 2. Code Quality (`.github/workflows/quality.yml`)
- **Triggers**: Push to main/master, Pull Requests
- **Jobs**:
  - **Lint**: Code linting with golangci-lint
  - **Coverage**: Test coverage reporting and badge updates

### 3. Security & Dependencies (`.github/workflows/security.yml`)
- **Triggers**: Weekly schedule, manual dispatch, Pull Requests
- **Jobs**:
  - **Security Scan**: Vulnerability scanning with gosec and govulncheck
  - **Dependency Review**: Security review of dependencies in PRs
  - **Update Dependencies**: Weekly automated dependency updates

## Required Secrets

Set these in GitHub repository Settings → Secrets and variables → Actions:

| Secret | Description | How to obtain |
|--------|-------------|---------------|
| `GCP_SA_KEY` | Google Cloud Service Account JSON | Create service account with Cloud Functions permissions |
| `GCP_PROJECT_ID` | Google Cloud Project ID | Your GCP project identifier |
| `TELEGRAM_BOT_TOKEN` | Telegram Bot Token | Get from @BotFather on Telegram |

## Deployment Process

1. **Development**: Create feature branch, make changes
2. **Pull Request**: Automated testing and quality checks
3. **Code Review**: Manual review and approval
4. **Merge to Main**: Automatic deployment to production
5. **Monitoring**: Check deployment status and function logs

## Local Development

```bash
# Install dependencies
go mod download

# Run quality checks (same as CI)
go build -v -o /dev/null .
go test -v ./...
go vet ./...
go fmt ./...

# Run linter (requires golangci-lint installation)
golangci-lint run
```

## Troubleshooting

- **Deployment failures**: Check Google Cloud Functions logs and GitHub Actions logs
- **Permission errors**: Verify service account has correct IAM roles
- **Secret errors**: Ensure all required secrets are set in GitHub repository settings
- **Build failures**: Run local tests and linting to identify issues before pushing