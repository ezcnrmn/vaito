# Vaito (Avito Clone)

Проект площадки объявлений, построенный на микросервисной архитектуре.

## Архитектура

```mermaid
graph TB
  Client[Client] -->|HTTP/REST| Gateway[Gateway service]

  subgraph "Microservices"
    Gateway -->|gRPC| User[User service]
    Gateway -->|gRPC| Listing[Listing service]
    Listing -->|gRPC| User

    User ~~~ Notification[Notification service]
  end

  Listing -->|Publish AMQP| RMQ[[RabbitMQ Broker]]
  RMQ -->|Consume AMQP| Notification

  subgraph DB ["PostgreSQL (Physical Instance)"]
    DB_User[(User Database)]
    DB_Listing[(Listing Database)]
  end

  User -->|SQL| DB_User
  Listing -->|SQL| DB_Listing

  DB ~~~ RMQ
```

[Схемы базы данных](#схемы-базы-данных)

## Стек технологий

- Language: Go 1.25+
- Communication: gRPC (Protobuf)
- Database: PostgreSQL
- DevOps: Docker Compose
- Docs: Swagger (OpenAPI)
- Broker: RabbitMQ

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

3. Поменять название сервиса `rabbitmq` на `localhost` в url

```
БЫЛО:
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
                                ^^^^^^^^
СТАЛО:
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
                                ^^^^^^^^^
```

4. Запуск сервисов

  - запуск базы данных в docker:

```bash
make db/up
```

  - запуск RabbitMQ в docker:

```bash
make rabbitmq/up
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
  - запуск сервиса notification:

```bash
make run/notification
```

5. После запуска api будет доступен на порту `:4000` (указывается в `.env.local`)

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


## Схемы базы данных

### Users

```mermaid
erDiagram
  User ||--o{ Role: "назначена"
  User ||--o{ Token: "владеет"
  Role ||--o{ RolePermission: "включает"
  Permission ||--o{ RolePermission: "представлено"

  User {
    int id PK
    string name
    string email UK
    binary password_hash
    int role_id FK
    timestamp created_at
    int version
  }

  Role {
    int id PK
    string name UK
  }

  Permission {
    int id PK
    string code UK
  }

  RolePermission {
    int role_id PK, FK
    int permission_id PK, FK
  }

  Token {
    binary hash PK
    int user_id FK
    string scope
    timestamp expires_at
  }
```

### Listings

```mermaid
erDiagram
  Listing ||--o{ ListingStatus: "имеет статус"
  Listing ||--o{ Category: "принадлежит"
  User ||..o{ Listing : "создаёт"

  Category {
    int id PK
    string name
  }

  ListingStatus {
    int id PK
    string name UK
  }

  Listing {
    int id PK
    string title
    string description
    int category_id FK
    int user_id "ID из внешней БД"
    timestamp created_at
    timestamp published_at
    int version
    int status_id FK
    int price
  }
```
