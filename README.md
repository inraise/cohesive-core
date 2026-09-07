# Cohesive Core

Бэкенд-сервис на Go для приложения по учёту домохозяйств: авторизация, дома, участники, приглашения, задачи. Модульный монолит с чётким разделением инфраструктуры и бизнес-логики.

---

## Содержание

- [Архитектура](#-архитектура)
- [Стек технологий](#-стек-технологий)
- [Быстрый старт](#-быстрый-старт)
- [Переменные окружения](#-переменные-окружения)
- [Команды Makefile](#-команды-makefile)
- [API-документация](#-api-документация)
- [Схема базы данных](#-схема-базы-данных)
- [Логирование](#-логирование)
- [Roadmap](#-roadmap)

---

## Архитектура

Clean Architecture внутри модульного монолита: общая инфраструктура — в `internal/core`, бизнес-логика изолирована по фичам в `internal/features`, каждая со слоями `repository → service → transport`.

```text
cohesive-core/
├── cmd/cohesive/               # main.go: сборка зависимостей и запуск сервера
├── docs/                       # Сгенерированная swag-документация
├── internal/
│   ├── core/
│   │   ├── config/              # Общая конфигурация (тайм-зона и т.п.)
│   │   ├── domain/               # Общие доменные типы (User, RefreshToken, Household, Task)
│   │   ├── errors/               # Базовые доменные ошибки → HTTP-статусы
│   │   ├── jwt/                   # Генерация и валидация JWT access-токенов
│   │   ├── logger/                # zap + ротация файлов логов
│   │   ├── repository/postgres/pool/  # Пул соединений PostgreSQL (pgx)
│   │   ├── repository/redis/pool/     # Клиент Redis (rate limiting)
│   │   └── transport/http/        # Сервер, роутер, middleware, request/response, swagger
│   │
│   └── features/
│       ├── auth/         # Регистрация, логин, refresh, logout
│       ├── users/         # Профиль текущего пользователя (/users/me)
│       ├── households/    # Дома, участники, роли, приглашения
│       └── tasks/          # Задачи внутри дома
│
├── migrations/                  # SQL-миграции (golang-migrate)
├── docker-compose.yaml          # cohesive, postgres, redis, migrate, port-forwarder
├── Makefile
└── .env.example
```

**Ключевые решения:**

- **Версионирование API** через `APIVersionRouter` (`/api/v1`, префикс срезается до хендлера).
- **Единая цепочка middleware:** `CORS → RequestID → Logger → Trace → Panic`, сквозной `request_id` в логах.
- **Фичи не знают друг о друге на уровне Go** — только через свой интерфейс сервиса и `core_domain`-типы. Физически одна БД — точечный SQL к чужой таблице допустим (например, `tasks` проверяет членство в доме прямым запросом к `household_members`).
- **Явные ошибки домена** (`core/errors`) мапятся в HTTP-статусы централизованно (`response`-пакет): 400/401/403/404/409/429.
- **Stateless access + отзываемый refresh.** Access — обычный JWT, 15 минут, не хранится и не отзывается раньше `exp`. Refresh — opaque-строка, SHA-256 хеш в БД, отзыв и ротация — реальные.
- **Точечная авторизация через middleware** (`Authenticate`) — вешается на конкретные роуты, а не глобально; `user_id` берётся из контекста, не из тела запроса.
- **Rate limiting через Redis** (`INCR`+`EXPIRE`) на `/auth/login` и `/auth/refresh`. Недоступность Redis не блокирует запросы.
- **Partial update через `Nullable[T]`** — отличает «поле не прислали» от «поле = null» (кастомный `UnmarshalJSON`).
- **Оптимистичная блокировка** через `version` у `User`/`Household`/`Task` — конкурентная запись → `409`.
- **Роль — свойство членства**, не сущности: `household_members.role`, не колонка в `households`. Единственность `owner` держится атомарной операцией `TransferOwnership`, не схемой.
- **Многошаговые операции без Go-транзакций** — `Pool` не даёт `Begin`/`Commit`, атомарность там, где нужна (создание дома + owner-membership, приём инвайта, передача владения) — через data-modifying CTE одним SQL-запросом.

---

## Стек технологий

| Категория           | Технология                                                                   |
| ------------------- | ---------------------------------------------------------------------------- |
| Язык                | Go 1.26                                                                      |
| HTTP                | `net/http` (`http.ServeMux`, Go 1.22+ `{wildcard}`-паттерны), без фреймворка |
| База данных         | PostgreSQL 17 + [`pgx/v5`](https://github.com/jackc/pgx)                     |
| Кэш / rate limiting | Redis 7 + [`go-redis/v9`](https://github.com/redis/go-redis)                 |
| Миграции            | [`golang-migrate`](https://github.com/golang-migrate/migrate)                |
| Логирование         | [`zap`](https://github.com/uber-go/zap)                                      |
| Конфигурация        | [`envconfig`](https://github.com/kelseyhightower/envconfig)                  |
| Валидация           | [`go-playground/validator`](https://github.com/go-playground/validator)      |
| Хеширование паролей | `bcrypt`                                                                     |
| Токены              | JWT (`golang-jwt/jwt/v5`) для access, opaque + SHA-256 для refresh           |
| API-документация    | [`swaggo/swag`](https://github.com/swaggo/swag) + `http-swagger`             |
| Контейнеризация     | Docker / Docker Compose                                                      |

---

## Быстрый старт

### Предварительные требования

- Go 1.26+
- Docker и Docker Compose

### Шаги

```bash
cp .env.example .env        # 1. заполнить переменные окружения (см. таблицу ниже)
make env-up                 # 2. поднять PostgreSQL и Redis
make migrate-up             # 3. применить миграции
make cohesive-run           # 4. запустить приложение локально
```

Сервис поднимется на `HTTP_ADDR` (по умолчанию `http://localhost:5050`).

Альтернатива — всё в Docker: `make cohesive-deploy` / `make cohesive-undeploy`.

---

## Переменные окружения

| Переменная              | Обязательна | По умолчанию | Описание                               |
| ----------------------- | ----------- | ------------ | -------------------------------------- |
| `HTTP_ADDR`             | +           | —            | Адрес HTTP-сервера, например `:5050`   |
| `HTTP_SHUTDOWN_TIMEOUT` |             | `30s`        | Таймаут graceful shutdown              |
| `ALLOWED_ORIGINS`       | +           | —            | Origin'ы для CORS через запятую        |
| `POSTGRES_HOST`         | +           | —            | Хост PostgreSQL                        |
| `POSTGRES_PORT`         |             | `5432`       | Порт PostgreSQL                        |
| `POSTGRES_USER`         | +           | —            | Пользователь БД                        |
| `POSTGRES_PASSWORD`     | +           | —            | Пароль БД                              |
| `POSTGRES_DB`           | +           | —            | Имя базы данных                        |
| `POSTGRES_TIMEOUT`      | +           | —            | Таймаут соединения с БД                |
| `REDIS_ADDR`            | +           | —            | Адрес Redis, например `localhost:6379` |
| `REDIS_PASSWORD`        |             | `""`         | Пароль Redis                           |
| `REDIS_DB`              |             | `0`          | Номер БД Redis                         |
| `JWT_SECRET`            | +           | —            | Секрет подписи access-токенов (HMAC)   |
| `JWT_ACCESS_TTL`        |             | `15m`        | Время жизни access-токена              |
| `JWT_REFRESH_TTL`       |             | `720h`       | Время жизни refresh-токена             |
| `LOGGER_LEVEL`          |             | `DEBUG`      | Уровень логирования                    |
| `LOGGER_FOLDER`         | +           | —            | Папка для файлов логов                 |
| `TIME_ZONE`             |             | `UTC`        | Тайм-зона приложения                   |

> `make cohesive-run` подставляет `LOGGER_FOLDER`/`POSTGRES_HOST` автоматически.

---

## Команды Makefile

| Команда                                 | Что делает                                           |
| --------------------------------------- | ---------------------------------------------------- |
| `make env-up` / `make env-down`         | Поднять / остановить PostgreSQL и Redis              |
| `make env-port-forward` / `-close`      | Проброс порта `5432` наружу через `socat`            |
| `make env-cleanup`                      | Полностью снести окружение и данные БД               |
| `make migrate-create seq=<name>`        | Создать новую пару миграций `up`/`down`              |
| `make migrate-up` / `make migrate-down` | Применить / откатить миграции                        |
| `make cohesive-run`                     | Запустить приложение локально (`go run`)             |
| `make cohesive-deploy` / `-undeploy`    | Собрать/запустить или остановить приложение в Docker |
| `make logs-cleanup`                     | Очистить локальные логи                              |
| `make ps`                               | Статус контейнеров Compose                           |
| `refresh-swag`                          | Обновление конфигурации Swagger                      |

---

## API-документация

Полное описание всех эндпоинтов (auth, users, households, tasks — запросы, ответы, коды ошибок) — в Swagger UI, генерируется из `@swag`-аннотаций над хендлерами (`swaggo/swag`):

**http://localhost:5050/swagger/index.html**

JSON-спека отдельно: `http://localhost:5050/swagger/doc.json`.

> Раздача Swagger подключается вызовом `httpServer.RegisterSwagger()` в `main.go` — если ещё не добавлен, эндпоинты выше не заработают. Сама спека собирается командой `swag init` (перегенерировать после правки `@swag`-аннотаций).

---

## Схема базы данных

| Таблица             | Миграция                | Назначение                                                                     |
| ------------------- | ----------------------- | ------------------------------------------------------------------------------ |
| `users`             | `000001_init_schema`    | Аккаунты пользователей                                                         |
| `refresh_tokens`    | `000002_refresh_tokens` | Хеши refresh-токенов, отзыв и ротация                                          |
| `households`        | `000003_households`     | Дома (`id`, `name`, `version`)                                                 |
| `household_members` | `000003_households`     | Членство: `household_id` + `user_id` + `role`, `UNIQUE(household_id, user_id)` |
| `household_invites` | `000003_households`     | Инвайт-коды: срок жизни, лимит использований, отзыв                            |
| `tasks`             | `000004_tasks`          | Задачи внутри дома: `title`, `status`, `assigned_to`                           |

Все таблицы, кроме `users`, ссылаются на `households`/`users` с `ON DELETE CASCADE` (кроме `tasks.assigned_to` — `ON DELETE SET NULL`, задача не удаляется вместе с исполнителем). Точные колонки и ограничения — в файлах `migrations/*.up.sql` или в Swagger-моделях ответов.

---

## Логирование

`zap`, структурированный вывод, отдельный файл на запуск в `LOGGER_FOLDER` (имя — таймстемп старта), уровень — `LOGGER_LEVEL`. Каждый запрос получает сквозной `request_id`.

---

## Roadmap

Проект в активной разработке.

**Готово:** JWT auth (login/refresh/logout с ротацией), `/users/me` CRUD, `/households` (дома, участники, роли, передача владения, инвайты), `/households/{id}/tasks` CRUD, rate limiting через Redis.
