# Quiz API

REST API для викторин и тестов с использованием Go, PostgreSQL и Docker.

## 🚀 Быстрый старт

### Запуск приложения

1. **Клонируйте репозиторий:**
```bash
git clone <repository-url>
cd quiz-api
```

2. **Настройте .env:**
```
# Application
APP_NAME=quiz-api
APP_ENV=local
APP_SHUTDOWN_TIMEOUT=15s
APP_DB_URL=postgres://quiz_user:quiz_password@db:5432/quiz_db?sslmode=disable

# HTTP Server
HTTP_PORT=8080
HTTP_HOST=0.0.0.0
HTTP_READ_TIMEOUT=10s
HTTP_WRITE_TIMEOUT=10s

# Database
POSTGRES_NAME=quiz_db
POSTGRES_USER=quiz_user
POSTGRES_PASSWORD=quiz_password
```

3. **Запустите приложение:**
```bash
docker-compose up --build
```
