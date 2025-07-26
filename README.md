# Online Subscriptions Service

REST API-сервис для управления онлайн-подписками. Поддерживает создание, чтение, обновление, удаление подписок, а также подсчёт общей суммы расходов по фильтрам.

## Стек технологий

- Golang  
- PostgreSQL  
- GORM  
- Gorilla Mux  
- Swagger (swaggo/swag)  
- Docker + Docker Compose

## Структура проекта

├── cmd/server # main.go — точка входа
├── internal
│ ├── handler # Обработчики API, разделены по CRUD
│ └── model # GORM-модель Subscription
│ └── reposirory # Подключение к хранилищу 
│ └── migrations # SQL-миграции
│ └── config # Чтение .env 
├── docs # Swagger-документация (авто)
├── .env # Переменные окружения
├── docker-compose.yml # Контейнеризация
└── README.md # Этот файл

## API-эндпоинты
* POST /subscriptions — создать подписку
* GET /subscriptions — получить все подписки
* GET /subscriptions/{user_id} — подписки по пользователю
* PUT /subscriptions/{id} — обновить подписку
* DELETE /subscriptions/{id} — удалить подписку
* GET /subscriptions/summary — сумма подписок по фильтрам

## Пример .env

DB_HOST=db
DB_PORT=5432
DB_USER=user
DB_PASSWORD=pass
DB_NAME=subscription_db
