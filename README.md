# Cohesive Core

Бэкенд-сервис на Go для приложения по учёту домохозяйств: авторизация пользователей, домохозяйства и связанные с ними сущности. Построен как модульный монолит с чётким разделением инфраструктуры и бизнес-логики.

-----

## Содержание

- [Архитектура](#-архитектура)
- [Стек технологий](#-стек-технологий)
- [Быстрый старт](#-быстрый-старт)
- [Переменные окружения](#-переменные-окружения)
- [Команды Makefile](#-команды-makefile)
- [API](#-api)
- [Схема базы данных](#-схема-базы-данных)
- [Логирование](#-логирование)
- [Roadmap](#-roadmap)

-----

## Архитектура

Проект следует принципам **Clean Architecture** внутри **модульного монолита**: общая инфраструктура вынесена в `internal/core`, а бизнес-логика изолирована по фичам в `internal/features`, каждая — со своими слоями `repository → service → transport`.

```text
cohesive-core/
├── cmd/
│   └── cohesive/
│       ├── main.go            # Точка входа: сборка зависимостей и запуск сервера
│       └── Dockerfile
│
├── internal/
│   ├── core/                  # Инфраструктурный слой, общий для всех фич
│   │   ├── config/            # Общая конфигурация приложения (тайм-зона и т.п.)
│   │   ├── domain/             # Общие доменные типы (User и др.)
│   │   ├── errors/             # Базовые доменные ошибки (NotFound, Conflict, ...)
│   │   ├── logger/              # Обёртка над zap + ротация файлов логов
│   │   ├── repository/postgres/ # Пул соединений с PostgreSQL (pgx)
│   │   └── transport/http/      # HTTP-сервер: роутер, middleware, request/response
│   │
│   └── features/
│       └── auth/                # Регистрация и авторизация пользователей
│           ├── repository/postgres/  # SQL-запросы
│           ├── service/              # Бизнес-логика (хеширование пароля, валидация)
│           └── transport/http/       # HTTP-хендлеры и DTO
│
├── migrations/                 # SQL-миграции (golang-migrate)
├── docker-compose.yaml         # cohesive, cohesive-postgres, migrate, port-forwarder
├── Makefile
└── .env.example
```

**Ключевые архитектурные решения:**

- **API-роутинг с версионированием.** `APIVersionRouter` регистрирует роуты фичи под префиксом `/api/v1`, который «срезается» перед тем, как запрос доходит до хендлера — фичи ничего не знают о версии API.
- **Единая цепочка middleware.** На сервер накручены `CORS → RequestID → Logger → Trace → Panic` — запросы логируются и трассируются сквозным `request_id`, а паника в хендлере не роняет процесс.
- **Фичи не знают друг о друге.** `auth` работает только через собственный интерфейс `AuthService` и общий `core_domain.User` — добавление новой фичи (например, `households`) не требует правок в существующих.
- **Явные ошибки домена.** `core/errors` определяет базовый набор ошибок (`ErrNotFound`, `ErrInvalidArgument`, `ErrConflict`), которые оборачиваются на каждом слое и мапятся в HTTP-статусы в `response`-пакете.

-----

## Стек технологий

|Категория          |Технология                                                                |
|-------------------|--------------------------------------------------------------------------|
|Язык               |Go 1.26                                                                   |
|HTTP               |`net/http` (`http.ServeMux`), без веб-фреймворка                          |
|База данных        |PostgreSQL 17 + [`pgx/v5`](https://github.com/jackc/pgx) (connection pool)|
|Миграции           |[`golang-migrate`](https://github.com/golang-migrate/migrate)             |
|Логирование        |[`zap`](https://github.com/uber-go/zap)                                   |
|Конфигурация       |[`envconfig`](https://github.com/kelseyhightower/envconfig)               |
|Валидация          |[`go-playground/validator`](https://github.com/go-playground/validator)   |
|Хеширование паролей|`bcrypt`                                                                  |
|Контейнеризация    |Docker / Docker Compose                                                   |

-----

## Быстрый старт

### Предварительные требования

- Go 1.26+
- Docker и Docker Compose

### 1. Переменные окружения

```bash
cp .env.example .env
```

Заполните `.env` своими значениями (см. [таблицу переменных](#-переменные-окружения) ниже).

### 2. База данных

Поднять только PostgreSQL в Docker:

```bash
make env-up
```

### 3. Миграции

```bash
make migrate-up
```

### 4. Запуск приложения

Локально (без контейнера, БД — из шага 2):

```bash
make cohesive-run
```

Сервис поднимется на адресе, указанном в `HTTP_ADDR` (по умолчанию `http://localhost:5050`).

### Альтернатива: всё в Docker

Если не хочется поднимать Go локально — можно собрать и запустить сам сервис в контейнере:

```bash
make cohesive-deploy   # сборка и запуск контейнера cohesive
make cohesive-undeploy # остановка
```

-----

## Переменные окружения

|Переменная             |Обязательна|По умолчанию|Описание                                               |
|-----------------------|-----------|------------|-------------------------------------------------------|
|`HTTP_ADDR`            |+          |—           |Адрес, на котором слушает HTTP-сервер, например `:5050`|
|`HTTP_SHUTDOWN_TIMEOUT`|           |`30s`       |Таймаут graceful shutdown                              |
|`ALLOWED_ORIGINS`      |+          |—           |Список origin’ов для CORS через запятую                |
|`POSTGRES_HOST`        |+          |—           |Хост PostgreSQL                                        |
|`POSTGRES_PORT`        |           |`5432`      |Порт PostgreSQL                                        |
|`POSTGRES_USER`        |+          |—           |Пользователь БД                                        |
|`POSTGRES_PASSWORD`    |+          |—           |Пароль БД                                              |
|`POSTGRES_DB`          |+          |—           |Имя базы данных                                        |
|`POSTGRES_TIMEOUT`     |+          |—           |Таймаут соединения с БД                                |
|`JWT_SECRET`           |           |—           |Секрет для подписи JWT (зарезервировано, см. Roadmap)  |
|`JWT_ACCESS_TTL`       |           |`15m`       |TTL для JWT ключа                                      |
|`LOGGER_LEVEL`         |           |`DEBUG`     |Уровень логирования                                    |
|`LOGGER_FOLDER`        |+          |—           |Папка для файлов логов                                 |
|`TIME_ZONE`            |           |`UTC`       |Тайм-зона приложения                                   |


> При запуске через `make cohesive-run` переменные `LOGGER_FOLDER` и `POSTGRES_HOST` подставляются автоматически (логи пишутся в `./out/logs`, БД — на `localhost`).

-----

## Команды Makefile

|Команда                                |Что делает                                                                    |
|---------------------------------------|------------------------------------------------------------------------------|
|`make env-up`                          |Поднять PostgreSQL для локальной разработки                                   |
|`make env-down`                        |Остановить PostgreSQL                                                         |
|`make env-port-forward`                |Прокинуть порт `5432` наружу через `socat` (доступ к БД из контейнера снаружи)|
|`make env-port-close`                  |Закрыть проброс порта                                                         |
|`make env-cleanup`                     |Полностью снести окружение и данные БД (с подтверждением)                     |
|`make migrate-create seq=<name>`       |Создать новую пару миграций `up`/`down`                                       |
|`make migrate-up` / `make migrate-down`|Применить / откатить миграции                                                 |
|`make cohesive-run`                    |Запустить приложение локально (`go run`)                                      |
|`make cohesive-deploy`                 |Собрать и запустить приложение в Docker                                       |
|`make cohesive-undeploy`               |Остановить контейнер приложения                                               |
|`make logs-cleanup`                    |Очистить локальные логи (с подтверждением)                                    |
|`make ps`                              |Статус контейнеров Compose                                                    |

-----

## API

Базовый префикс всех эндпоинтов фич: **`/api/v1`**.

### `POST /api/v1/auth/register`

Регистрация нового пользователя.

**Request body**

```json
{
  "email": "user@example.com",
  "password": "supersecurepassword",
  "first_name": "John",
  "last_name": "Doe",
  "age": 28
}
```

|Поле        |Тип   |Обязательно|Валидация                                   |
|------------|------|-----------|--------------------------------------------|
|`email`     |string|+          |5–100 символов                              |
|`password`  |string|+          |10–100 символов, хранится в виде bcrypt-хеша|
|`first_name`|string|+          |3–100 символов                              |
|`last_name` |string|—          |3–100 символов                              |
|`age`       |int   |—          |0–130                                       |

**Response `201 Created`**

```json
{
  "id": "e5c1f2b0-...-uuid",
  "version": 1,
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "age": 28,
  "created_at": "2026-08-04T05:00:00Z",
  "updated_at": "2026-08-04T05:00:00Z"
}
```

**Ответ с ошибкой**

```json
{
  "error": "validate user domain: invalid `Email` len: 3: invalid argument",
  "message": "failed to create user"
}
```

**Пример запроса**

```bash
curl -X POST http://localhost:5050/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "supersecurepassword",
    "first_name": "John",
    "last_name": "Doe",
    "age": 28
  }'
```

### `POST /api/v1/auth/login` 

Аутентификация пользователя и получение токенов.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```
**Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1Ni...",
  "expires_at": "2026-08-05T20:00:00Z"
}
```

-----

## Схема базы данных

Таблица `users` (миграция `000001_init_schema`):

|Колонка        |Тип           |Ограничения                                |
|---------------|--------------|-------------------------------------------|
|`id`           |`UUID`        |`PRIMARY KEY`, `DEFAULT uuid_generate_v4()`|
|`version`      |`INT`         |`NOT NULL DEFAULT 1`                       |
|`email`        |`VARCHAR(100)`|`UNIQUE NOT NULL`, длина 5–100             |
|`password_hash`|`VARCHAR(255)`|`NOT NULL`                                 |
|`first_name`   |`VARCHAR(100)`|`NOT NULL`, длина 1–100                    |
|`last_name`    |`VARCHAR(100)`|длина 1–100, nullable                      |
|`age`          |`INT`         |0–130, nullable                            |
|`created_at`   |`TIMESTAMPTZ` |                                           |
|`updated_at`   |`TIMESTAMPTZ` |`CHECK(created_at <= updated_at)`          |

-----

## Логирование

Логи пишутся через `zap` в структурированном виде и складываются в `LOGGER_FOLDER` отдельным файлом на каждый запуск сервиса (имя файла — таймстемп старта). Уровень регулируется переменной `LOGGER_LEVEL`. Каждый HTTP-запрос получает `request_id`, который прокидывается через middleware и попадает в логи для сквозной трассировки.

-----

## Roadmap

Проект в активной разработке. В ближайших планах:

**Модуль пользователей**

- [ ] `GET /api/v1/users/me` и обновление профиля
- [ ] Роли и права доступа

**Модуль домохозяйств**

- [ ] Создание и редактирование домохозяйств
- [ ] Приглашение участников по ссылке/коду
- [ ] Разграничение прав внутри группы

**Качество и CI/CD**

- [ ] Unit-тесты для слоёв service и repository (моки)
- [ ] Интеграционные тесты через `testcontainers-go`
- [ ] `golangci-lint` и CI/CD pipeline на GitHub Actions