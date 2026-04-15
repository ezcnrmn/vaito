# Vaito (Avito Clone)

Проект площадки объявлений, построенный на микросервисной архитектуре.

## Архитектура

```mermaid
graph LR
  Client[Client] -->|HTTP/REST| Gateway[Gateway service]

  subgraph "Microservices"
    Gateway -->|gRPC| User[User service]
    Gateway -->|gRPC| Listing[Listing service]
    Listing -->|gRPC| User
  end

  subgraph "PostgreSQL (Physical Instance)"
    DB_User[(User Database)]
    DB_Listing[(Listing Database)]
  end

  User -->|SQL| DB_User
  Listing -->|SQL| DB_Listing
```

## Стек технологий

- Language: Go 1.25+
- Communication: gRPC (Protobuf)
- Database: PostgreSQL
- DevOps: Docker Compose
- Docs: Swagger (OpenAPI)

## Запуск

### Все сервисы в docker:

1. Создать в корне проекта `.env` файл. Скопировать в него содержимое `.env.example`.

2. Выполнить

```bash
docker compose up --build
```

3. После запуска api будет доступен на порту `:4000` (указывается в `.env`)

> Swagger доступен по `http://localhost:4000/swagger/index.html`

### Также возможен локальный запуск каждого сервиса:

1. Создать в корне проекта `.env.local` файл. Скопировать в него содержимое `.env.example`.

2. Поменять название сервиса базы данных `postgres_db` на `localhost` в dsn

```
БЫЛО:
USER_DB_DSN=postgres://vaito_user:pa55word@postgres_db:5432/users?sslmode=disable
                                           ^^^^^^^^^^^
СТАЛО:
USER_DB_DSN=postgres://vaito_user:pa55word@localhost:5432/users?sslmode=disable
                                           ^^^^^^^^^
```

3. Запуск сервисов

  - запуск базы данных в docker:

```bash
make db/up
```

  - запуск сервиса gateway:

```bash
make run/gateway
```

  - запуск сервиса user:

```bash
make run/user
```

  - запуск сервиса listing:

```bash
make run/listing
```

4. После запуска api будет доступен на порту `:4000` (указывается в `.env.local`)

> Swagger доступен по `http://localhost:4000/swagger/index.html`

## Тестовые данные

Для добавления двух пользователей (admin, user) и тестовых объявлений можно выполнить команду:

```bash
make db/seed
```

Войти под созданными пользователями можно с помощью email и пароля:

|          | Admin            | User            |
| -------- | ---------------- | --------------- |
| email    | `admin@test.com` | `user@test.com` |
| password | `12345678`       | `12345678`      |
