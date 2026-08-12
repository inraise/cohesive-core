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
│   │   ├── domain/             # Общие доменные типы (User, RefreshToken и др.)
│   │   ├── errors/             # Базовые доменные ошибки (NotFound, Conflict, Unauthorized, ...)
│   │   ├── jwt/                 # Генерация и валидация JWT access-токенов
│   │   ├── logger/              # Обёртка над zap + ротация файлов логов
│   │   ├── repository/postgres/ # Пул соединений с PostgreSQL (pgx)
│   │   └── transport/http/      # HTTP-сервер: роутер, middleware (в т.ч. Auth), request/response
│   │
│   └── features/
│       ├── auth/                # Регистрация, логин, refresh и logout
│       │   ├── repository/postgres/  # SQL-запросы (users, refresh_tokens)
│       │   ├── service/              # Бизнес-логика (хеширование пароля, JWT, ротация refresh-токенов)
│       │   └── transport/http/       # HTTP-хендлеры и DTO
│       │
│       └── users/                # Профиль текущего пользователя (/users/me)
│           ├── repository/postgres/  # SQL-запросы (users)
│           ├── service/              # GetMe / PatchMe (partial update + оптимистичная блокировка) / DeleteMe
│           └── transport/http/       # HTTP-хендлеры и DTO, все роуты под Authenticate middleware
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
- **Явные ошибки домена.** `core/errors` определяет базовый набор ошибок (`ErrNotFound`, `ErrInvalidArgument`, `ErrConflict`, `ErrUnauthorized`), которые оборачиваются на каждом слое и мапятся в HTTP-статусы в `response`-пакете.
- **Stateless access + отзываемый refresh.** Access-токен — обычный подписанный JWT (`core/jwt`), сервер его не хранит и не может отозвать раньше `exp` (15 минут). Refresh-токен — непрозрачная случайная строка, её SHA-256 хеш живёт в таблице `refresh_tokens`; именно это позволяет по-настоящему отзывать сессии на `/auth/logout` и делать ротацию на `/auth/refresh`.
- **Точечная авторизация через middleware.** `Authenticate` (`core/transport/http/middleware/auth.go`) вешается на конкретные роуты через `Route.Middleware`, а не глобально на сервер — так публичные `/auth/*`-эндпоинты остаются без токена, а все `/users/me` его требуют. Хендлер достаёт `user_id` из контекста (`UserIDFromContext`), а не из тела/query запроса — иначе можно было бы подставить чужой id, имея свой валидный токен.
- **Partial update через `Nullable[T]`.** `PATCH /users/me` отличает «поле не прислали» от «поле прислали как `null`» с помощью generic-обёртки `Nullable[T]` с кастомным `UnmarshalJSON` (`core/transport/http/types`) — JSON-декодер вызывает `UnmarshalJSON` только для ключей, которые реально есть в теле запроса. Домен (`core_domain.UserPatch.ApplyPatch`) применяет только `Set == true` поля и валидирует результат целиком, не зная деталей HTTP-слоя.
- **Оптимистичная блокировка через `Version`.** `PatchMe` в репозитории обновляет строку через `WHERE id = $1 AND version = $2` и одновременно увеличивает `version`; если конкурентный запрос успел изменить профиль между чтением и записью — `UPDATE` не находит строку, и это мапится в `409 Conflict`.

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
|Токены             |JWT ([`golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt)) для access, opaque-строка + SHA-256 для refresh|
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
|`JWT_SECRET`           |+          |—           |Секрет для подписи access-токенов (HMAC)               |
|`JWT_ACCESS_TTL`       |           |`15m`       |Время жизни access-токена                              |
|`JWT_REFRESH_TTL`      |           |`720h`      |Время жизни refresh-токена (30 дней)                   |
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

-----

### `POST /api/v1/auth/login`

Логин по email/паролю. Выдаёт пару токенов: короткоживущий access (JWT, `JWT_ACCESS_TTL`) и долгоживущий refresh (opaque-строка, `JWT_REFRESH_TTL`), хеш которого сохраняется в `refresh_tokens`.

**Request body**

```json
{
  "email": "user@example.com",
  "password": "supersecurepassword"
}
```

**Response `200 OK`**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "9f1c2b7a3e5d...",
  "expires_at": "2026-08-08T12:15:00Z"
}
```

`expires_at` относится к `access_token`. Неверный email или пароль дают `400 Bad Request` — намеренно одну и ту же ошибку для обоих случаев, чтобы нельзя было перебором выяснить, какие email зарегистрированы.

**Пример запроса**

```bash
curl -X POST http://localhost:5050/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "supersecurepassword"
  }'
```

-----

### `POST /api/v1/auth/refresh`

Меняет валидный refresh-токен на новую пару access + refresh. Использованный refresh-токен сразу отзывается (ротация) — повторно предъявить его нельзя, даже если он был перехвачен.

**Request body**

```json
{
  "refresh_token": "9f1c2b7a3e5d..."
}
```

**Response `200 OK`** — такой же формат, как у `/auth/login`.

Невалидный, истёкший, уже отозванный или уже использованный refresh-токен — `401 Unauthorized`.

**Пример запроса**

```bash
curl -X POST http://localhost:5050/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "9f1c2b7a3e5d..."}'
```

-----

### `POST /api/v1/auth/logout`

Отзывает refresh-токен: им больше нельзя получить новый access-токен. Access-токен, уже выданный клиенту, продолжит работать до истечения своего TTL — это осознанный компромисс stateless access-токенов, полноценного server-side kill switch для них нет.

**Request body**

```json
{
  "refresh_token": "9f1c2b7a3e5d..."
}
```

**Response `204 No Content`**

Если токен уже не существует/уже отозван — тоже `204`, logout идемпотентен.

**Пример запроса**

```bash
curl -X POST http://localhost:5050/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "9f1c2b7a3e5d..."}'
```

-----

### `GET /api/v1/users/me`

Профиль текущего пользователя. Требует заголовок `Authorization: Bearer <access_token>` — обрабатывается `Authenticate` middleware до хендлера.

**Response `200 OK`**

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

`password_hash` в ответе никогда не присутствует — DTO собирается вручную, без него.

Нет/просрочен/невалиден токен — `401 Unauthorized`, до хендлера дело не доходит.

**Пример запроса**

```bash
curl http://localhost:5050/api/v1/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

-----

### `PATCH /api/v1/users/me`

Частичное обновление профиля. Меняются только присланные поля — отсутствие ключа в JSON и присланное значение `null` различаются (см. `Nullable[T]` в архитектурных решениях выше). Пароль передаётся как обычный текст (сервер сам его хеширует), а не как хеш.

**Request body** (любое подмножество полей)

```json
{
  "email": "new-email@example.com",
  "password": "newsupersecurepassword",
  "first_name": "Jane",
  "last_name": null,
  "age": 29
}
```

|Поле        |Тип   |Валидация при указании                      |
|------------|------|---------------------------------------------|
|`email`     |string|5–100 символов                                |
|`password`  |string|10–100 символов; хешируется bcrypt перед сохранением, `password_hash` клиенту не возвращается|
|`first_name`|string|1–100 символов, нельзя явно выставить `null`  |
|`last_name` |string|1–100 символов, можно явно выставить `null`   |
|`age`       |int   |0–130, можно явно выставить `null`             |

**Response `200 OK`** — обновлённый профиль, формат как у `GET /users/me`.

Конфликт версии (профиль успели изменить между чтением и записью, например из другой сессии) — `409 Conflict`. Некорректные значения — `400 Bad Request`.

**Пример запроса**

```bash
curl -X PATCH http://localhost:5050/api/v1/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "Content-Type: application/json" \
  -d '{"first_name": "Jane"}'
```

-----

### `DELETE /api/v1/users/me`

Безвозвратное удаление аккаунта (жёсткое, без soft-delete). Каскадом удаляются и все `refresh_tokens` пользователя (`ON DELETE CASCADE`).

**Response `204 No Content`**

**Пример запроса**

```bash
curl -X DELETE http://localhost:5050/api/v1/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
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

Таблица `refresh_tokens` (миграция `000002_refresh_tokens`):

|Колонка      |Тип           |Ограничения                                        |
|-------------|--------------|----------------------------------------------------|
|`id`         |`UUID`        |`PRIMARY KEY`, `DEFAULT uuid_generate_v4()`          |
|`user_id`    |`UUID`        |`NOT NULL`, `REFERENCES users(id) ON DELETE CASCADE` |
|`token_hash` |`VARCHAR(64)` |`UNIQUE NOT NULL` (SHA-256 hex от refresh-токена)    |
|`expires_at` |`TIMESTAMPTZ` |`NOT NULL`                                           |
|`revoked_at` |`TIMESTAMPTZ` |nullable — `NULL`, пока токен активен                |
|`created_at` |`TIMESTAMPTZ` |`NOT NULL DEFAULT now()`                             |

Индекс `idx_refresh_tokens_user_id` — по `user_id`, на будущее для операций вида «отозвать все сессии пользователя».

-----

## Логирование

Логи пишутся через `zap` в структурированном виде и складываются в `LOGGER_FOLDER` отдельным файлом на каждый запуск сервиса (имя файла — таймстемп старта). Уровень регулируется переменной `LOGGER_LEVEL`. Каждый HTTP-запрос получает `request_id`, который прокидывается через middleware и попадает в логи для сквозной трассировки.

-----

## Roadmap

Проект в активной разработке. В ближайших планах:

**Модуль пользователей**

- [ ] Роли и права доступа

**Модуль семьи**

- [ ] Создание и редактирование семьи
- [ ] Приглашение участников по ссылке/коду
- [ ] Разграничение прав внутри группы

**Качество и CI/CD**

- [ ] Unit-тесты для слоёв service и repository (моки)
- [ ] Интеграционные тесты через `testcontainers-go`
- [ ] `golangci-lint` и CI/CD pipeline на GitHub Actions