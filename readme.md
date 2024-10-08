# WB Orders Service

WB Orders Service — это микросервис для обработки заказов, который использует Kafka для чтения сообщений о новых заказах и PostgreSQL для их хранения. Приложение также имеет кеширование данных для быстрого доступа к заказам.

## Стек технологий

- **Go** — основной язык разработки.
- **Kafka** — очередь сообщений для обработки заказов.
- **PostgreSQL** — реляционная база данных для хранения заказов.
- **Docker** — контейнеризация сервисов (Kafka, Zookeeper, PostgreSQL).
- **Gin** — фреймворк для создания HTTP-сервера.
- **Viper** — библиотека для управления конфигурациями.
- **Zap** — логирование.
- **GORM** — ORM для взаимодействия с PostgreSQL.
- **golangci-lint** — инструмент для проверки кода (линтер).

## Установка и запуск

### 1. Клонирование репозитория

```bash
git clone git@github.com:mmmaks777/WB.git
cd wb-orders-service
```

### 2. Запуск всех сервисов (Kafka, Zookeeper и PostgreSQL) через docker-compose

```bash
docker-compose up -d
```

### 3. Запуск приложения

```bash
go run main.go
```

### 4. Отправка тестового сообщения в Kafka

```bash
go run order_usecase.go
```

### 5. API (получение данные по id)

```bash
curl http://localhost:8081/orders/:id
```