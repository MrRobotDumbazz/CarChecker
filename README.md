# Car Status Backend

Go backend сервис для системы определения состояния автомобиля по фотографии. Использует чистый `net/http` без внешних фреймворков.

## Архитектура проекта

```
car-status-detection/
├── frontend/           # React/Vue фронтенд (отдельная команда)
├── backend/           # Go API сервер (этот репозиторий)
└── ml/               # Python ML сервис (отдельная команда)
```

## Структура backend/

```
backend/
├── cmd/server/         # Точка входа в приложение
├── internal/
│   ├── config/        # Конфигурация
│   ├── database/      # Подключение к БД и миграции
│   ├── handlers/      # HTTP обработчики
│   ├── models/        # Модели данных
│   ├── services/      # Бизнес-логика
│   ├── middleware/    # HTTP middleware
│   └── server/        # HTTP сервер и роутинг
├── pkg/utils/         # Общие утилиты
├── uploads/           # Загруженные изображения
└── shared/uploads/    # Shared storage с ML сервисом
```

## API Documentation

### 📚 Интерактивная документация

Полная документация API доступна в нескольких форматах:

- **🌐 [Swagger UI](http://localhost:8081/api/docs/swagger)** - интерактивная документация с возможностью тестирования
- **📖 [ReDoc](http://localhost:8081/api/docs/redoc)** - красивая документация в современном стиле
- **📄 [OpenAPI Spec](http://localhost:8081/api/docs/openapi.yaml)** - машиночитаемая спецификация
- **🏠 [Documentation Index](http://localhost:8081/api/docs)** - главная страница документации

### 🚀 Быстрый старт с API

```bash
# 1. Проверить состояние сервиса
curl http://localhost:8081/api/v1/health

# 2. Загрузить изображение автомобиля
curl -X POST http://localhost:8081/api/v1/images/upload \
  -F "image=@my_car.jpg"

# 3. Запустить анализ (получить image_id из шага 2)
curl -X POST http://localhost:8081/api/v1/predict/{image_id}

# 4. Получить результат анализа
curl http://localhost:8081/api/v1/predictions/{prediction_id}
```

## API Endpoints

### Health Check
- `GET /api/v1/health` - Общий health check всех компонентов
- `GET /api/v1/health/ready` - Readiness probe для Kubernetes
- `GET /api/v1/health/live` - Liveness probe для Kubernetes

### Изображения
- `POST /api/v1/images/upload` - Загрузка изображения автомобиля
- `GET /api/v1/images/{id}` - Получение метаданных изображения
- `DELETE /api/v1/images/{id}` - Удаление изображения

### Анализ автомобилей
- `POST /api/v1/predict/{image_id}` - Запуск анализа состояния автомобиля
- `GET /api/v1/predictions/{id}` - Получение результата анализа
- `GET /api/v1/predictions/stats` - Статистика анализов за 24 часа

### Документация
- `GET /api/docs` - Главная страница документации
- `GET /api/docs/swagger` - Swagger UI
- `GET /api/docs/redoc` - ReDoc документация

## Запуск

### 1. Подготовка окружения

```bash
# Клонируйте репозиторий
git clone <repository-url>
cd CarChecker

# Скопируйте конфигурацию
cp .env.example .env

# Отредактируйте .env под ваше окружение
vim .env
```

### 2. Запуск базы данных

```bash
# Только PostgreSQL
docker-compose up postgres

# С RabbitMQ (если нужны очереди)
docker-compose --profile queue up postgres rabbitmq

# С Kafka (альтернатива RabbitMQ)
docker-compose --profile kafka up postgres zookeeper kafka

# Со всеми сервисами
docker-compose --profile queue --profile cache up
```

### 3. Запуск сервера

```bash
# Установка зависимостей
go mod tidy

# Создание папки для загрузок
mkdir -p uploads shared/uploads

# Запуск сервера
go run cmd/server/main.go
```

Сервер будет доступен на `http://localhost:8080`

## Конфигурация

Основная конфигурация через переменные окружения в `.env`:

```bash
# Сервер
SERVER_HOST=localhost
SERVER_PORT=8080

# База данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=car_status_db

# ML сервис
ML_SERVICE_URL=http://localhost:8000
ML_SERVICE_TIMEOUT=30s

# Файловое хранилище
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=10485760  # 10MB
ALLOWED_TYPES=image/jpeg,image/jpg,image/png

# Очереди (опционально)
QUEUE_ENABLED=false
QUEUE_TYPE=db  # "rabbitmq" | "kafka" | "db"
```

## Интеграция с ML сервисом

### 🔗 Для ML команды (Python)

Необходимо реализовать HTTP API на порту 8000 со следующим endpoint:

#### `POST http://localhost:8000/api/predict`

**Request:**
```json
{
  "image_path": "/absolute/path/to/image.jpg",
  "model_version": "v1.0.0" // опционально
}
```

**Response (success):**
```json
{
  "success": true,
  "cleanliness": {
    "status": "clean",     // "clean" | "dirty"
    "confidence": 0.8765
  },
  "integrity": {
    "status": "intact",    // "intact" | "damaged"
    "confidence": 0.9234
  },
  "processing_time_ms": 2450,
  "model_version": "v1.2.0"
}
```

**Response (error):**
```json
{
  "success": false,
  "error": "Model inference failed"
}
```

📋 **Подробная документация для ML команды:** [Swagger UI - ML Service](http://localhost:8081/api/docs/swagger#tag/ML-Service)

### Shared Storage
- Go backend сохраняет файлы в `./uploads/` или `./shared/uploads/`
- ML сервис получает абсолютный путь к файлу в поле `image_path`
- Обе службы должны иметь доступ к одной файловой системе

### Режимы работы

#### Синхронный режим (по умолчанию)
- Прямые HTTP запросы к ML API
- Подходит для быстрых предсказаний (<5 сек)
- Пользователь сразу получает результат

#### Асинхронный режим (опционально)
- Через очереди: RabbitMQ, Kafka или DB queue
- Для длительных предсказаний (>10 сек)
- Пользователь получает статус "queued", затем polling результата

## Shared Storage

Изображения сохраняются в `shared/uploads/` директории, которая доступна обоим сервисам:
- Go backend сохраняет файлы
- ML сервис читает файлы по абсолютному пути

## Примеры использования

### Загрузка изображения

```bash
curl -X POST http://localhost:8080/api/v1/images/upload \
  -F "image=@car.jpg" \
  -H "Content-Type: multipart/form-data"
```

### Запуск предсказания

```bash
curl -X POST http://localhost:8080/api/v1/predict/{image_id} \
  -H "Content-Type: application/json"
```

### Получение результата

```bash
curl http://localhost:8080/api/v1/predictions/{prediction_id}
```

## Разработка

### Структура запроса/ответа

**Успешный ответ:**
```json
{
  "success": true,
  "data": { ... },
  "message": "Operation completed successfully",
  "timestamp": "2023-11-20T10:00:00Z"
}
```

**Ошибка:**
```json
{
  "error": "Bad Request",
  "message": "Invalid image format",
  "timestamp": "2023-11-20T10:00:00Z"
}
```

### Добавление новых endpoints

1. Создайте handler в `internal/handlers/`
2. Добавьте route в `internal/server/routes.go`
3. Зарегистрируйте в `internal/server/server.go`

### Тестирование

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Список всех endpoints
curl http://localhost:8080/
```

## Производственное развёртывание

### Docker

```bash
# Сборка образа
docker build -t car-status-backend .

# Запуск контейнера
docker run -p 8080:8080 --env-file .env car-status-backend
```

### Kubernetes

```bash
# Применить конфигурацию
kubectl apply -f k8s/
```

## Мониторинг

### Health Checks
- `/api/v1/health` - полная проверка всех компонентов
- `/api/v1/health/ready` - готовность к приёму запросов
- `/api/v1/health/live` - проверка работоспособности

### Логирование
- Структурированные логи всех HTTP запросов
- Логи ошибок с контекстом
- Метрики производительности

## Технические решения

### Почему net/http вместо фреймворка?
- Минимальные зависимости
- Полный контроль над поведением
- Высокая производительность
- Простота отладки

### Когда использовать очереди?
- ML обработка > 10 секунд
- Высокая нагрузка
- Нужна надёжность доставки
- Несколько ML воркеров

### База данных
- PostgreSQL для надёжности
- Индексы для производительности
- Миграции при старте сервера

## Troubleshooting

### Проблемы подключения к БД
```bash
# Проверка состояния контейнера
docker-compose ps postgres

# Логи PostgreSQL
docker-compose logs postgres
```

### Проблемы с ML сервисом
```bash
# Проверка доступности
curl http://localhost:8000/health

# Проверка конфигурации
echo $ML_SERVICE_URL
```

### Проблемы с файлами
```bash
# Проверка прав на папку
ls -la uploads/

# Создание папки
mkdir -p uploads shared/uploads
chmod 755 uploads shared/uploads
```