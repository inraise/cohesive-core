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
│   │   ├── domain/             # Общие доменные типы (User, RefreshToken, Household и др.)
│   │   ├── errors/             # Базовые доменные ошибки (NotFound, Conflict, Unauthorized, Forbidden)
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
│       ├── users/                # Профиль текущего пользователя (/users/me)
│       │   ├── repository/postgres/  # SQL-запросы (users)
│       │   ├── service/              # GetMe / PatchMe (partial update + оптимистичная блокировка) / DeleteMe
│       │   └── transport/http/       # HTTP-хендлеры и DTO, все роуты под Authenticate middleware
│       │
│       └── households/           # Дома, участники и приглашения
│           ├── repository/postgres/  # SQL-запросы (households, household_members, household_invites)
│           ├── service/              # Бизнес-правила и проверки прав (owner/admin/member)
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
- **Фичи не знают друг о друге.** Каждая фича работает только через собственный интерфейс сервиса и общие `core_domain`-типы — добавление новой фичи не требует правок в существующих.
- **Явные ошибки домена.** `core/errors` определяет базовый набор ошибок (`ErrNotFound`, `ErrInvalidArgument`, `ErrConflict`, `ErrUnauthorized`, `ErrForbidden`), которые оборачиваются на каждом слое и мапятся в HTTP-статусы (404/400/409/401/403) в `response`-пакете.
- **Stateless access + отзываемый refresh.** Access-токен — обычный подписанный JWT (`core/jwt`), сервер его не хранит и не может отозвать раньше `exp` (15 минут). Refresh-токен — непрозрачная случайная строка, её SHA-256 хеш живёт в таблице `refresh_tokens`; именно это позволяет по-настоящему отзывать сессии на `/auth/logout` и делать ротацию на `/auth/refresh`.
- **Точечная авторизация через middleware.** `Authenticate` (`core/transport/http/middleware/auth.go`) вешается на конкретные роуты через `Route.Middleware`, а не глобально на сервер — так публичные `/auth/*`-эндпоинты остаются без токена, а всё под `/users/me` и `/households/*` его требуют. Хендлер достаёт `user_id` из контекста (`UserIDFromContext`), а не из тела/query запроса — иначе можно было бы подставить чужой id, имея свой валидный токен.
- **Partial update через `Nullable[T]`.** `PATCH /users/me` отличает «поле не прислали» от «поле прислали как `null`» с помощью generic-обёртки `Nullable[T]` с кастомным `UnmarshalJSON` (`core/transport/http/types`).
- **Оптимистичная блокировка через `Version`.** И у `User`, и у `Household` есть поле `version`. Обновления идут через `WHERE id = $1 AND version = $2` с одновременным `version = version + 1`; если конкурентный запрос успел изменить строку между чтением и записью — `UPDATE` не находит строк, и это мапится в `409 Conflict`.
- **Роль — свойство членства, а не сущности.** `household_members.role` хранит роль конкретного юзера в конкретном доме (`owner`/`admin`/`member`), а не колонка в `households` — один и тот же человек может быть `owner` в одном доме и `member` в другом. Единственность `owner` в доме поддерживается не схемой, а атомарной операцией `TransferOwnership` (см. ниже).
- **Многошаговые операции без транзакций на уровне Go.** `Pool` (`core/repository/postgres/pool`) не предоставляет `Begin`/`Commit` — везде, где нужна атомарность нескольких `INSERT`/`UPDATE` (создание дома + owner-membership, приём инвайта, передача владения), используются **data-modifying CTE**: несколько операций в одном SQL-запросе, который Postgres выполняет как единое целое.

-----

## Стек технологий

|Категория          |Технология                                                                |
|-------------------|--------------------------------------------------------------------------|
|Язык               |Go 1.26                                                                   |
|HTTP               |`net/http` (`http.ServeMux`, паттерны Go 1.22+ с `{wildcard}` в пути), без веб-фреймворка|
|База данных        |PostgreSQL 17 + [`pgx/v5`](https://github.com/jackc/pgx) (connection pool)|
|Миграции           |[`golang-migrate`](https://github.com/golang-migrate/migrate)             |
|Логирование        |[`zap`](https://github.com/uber-go/zap)                                   |
|Конфигурация       |[`envconfig`](https://github.com/kelseyhightower/envconfig)               |
|Валидация          |[`go-playground/validator`](https://github.com/go-playground/validator)   |
|Хеширование паролей|`bcrypt`                                                                  |
|Токены             |JWT ([`golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt)) для access, opaque-строка + SHA-256 для refresh|
|Инвайт-коды        |Случайные base32-строки (`crypto/rand`)                                   |
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

Базовый префикс всех эндпоинтов фич: **`/api/v1`**. Эндпоинты, помеченные 🔒, требуют заголовок `Authorization: Bearer <access_token>`.

### Auth

|Метод и путь                        |Описание                                                        |
|-------------------------------------|-----------------------------------------------------------------|
|`POST /auth/register`                |Регистрация нового пользователя                                  |
|`POST /auth/login`                   |Логин по email/паролю → access + refresh токены                  |
|`POST /auth/refresh`                 |Обмен refresh-токена на новую пару (с ротацией)                  |
|`POST /auth/logout`                  |Отзыв refresh-токена                                              |

#### `POST /api/v1/auth/register`

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
|`first_name`|string|+          |1–100 символов                              |
|`last_name` |string|—          |1–100 символов                              |
|`age`       |int   |—          |0–130                                       |

**Response `201 Created`** — профиль пользователя (без `password_hash`).

#### `POST /api/v1/auth/login`

**Request body**: `{"email": "...", "password": "..."}`

**Response `200 OK`**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "9f1c2b7a3e5d...",
  "expires_at": "2026-08-08T12:15:00Z"
}
```

Неверный email или пароль → `400 Bad Request` — намеренно одна и та же ошибка для обоих случаев.

#### `POST /api/v1/auth/refresh`

**Request body**: `{"refresh_token": "..."}` → тот же формат ответа, что у `/login`. Использованный refresh-токен сразу отзывается (ротация). Невалидный/истёкший/уже использованный токен → `401 Unauthorized`.

#### `POST /api/v1/auth/logout`

**Request body**: `{"refresh_token": "..."}` → `204 No Content`. Идемпотентен: если токена уже нет — тоже `204`.

-----

### Users 🔒

|Метод и путь        |Описание                              |
|---------------------|----------------------------------------|
|`GET /users/me`      |Профиль текущего пользователя           |
|`PATCH /users/me`    |Частичное обновление профиля            |
|`DELETE /users/me`   |Безвозвратное удаление аккаунта         |

#### `GET /api/v1/users/me`

**Response `200 OK`** — профиль (без `password_hash`).

#### `PATCH /api/v1/users/me`

Меняются только присланные поля (см. `Nullable[T]` в архитектурных решениях). Пароль передаётся как обычный текст — сервер сам его хеширует.

**Request body** (любое подмножество)

```json
{
  "email": "new-email@example.com",
  "password": "newsupersecurepassword",
  "first_name": "Jane",
  "last_name": null,
  "age": 29
}
```

|Поле        |Тип   |Валидация при указании                                          |
|------------|------|------------------------------------------------------------------|
|`email`     |string|5–100 символов                                                     |
|`password`  |string|10–100 символов; хешируется bcrypt, в ответе не возвращается       |
|`first_name`|string|1–100 символов, нельзя явно выставить `null`                       |
|`last_name` |string|1–100 символов, можно явно выставить `null`                        |
|`age`       |int   |0–130, можно явно выставить `null`                                  |

**Response `200 OK`** — обновлённый профиль. Конфликт версии → `409 Conflict`.

#### `DELETE /api/v1/users/me`

Жёсткое удаление, без возможности восстановления. Каскадом удаляются `refresh_tokens` и членство во всех домах (`household_members`). **Response `204 No Content`**.

-----

### Households 🔒

|Метод и путь                                    |Кто может               |Описание                                    |
|--------------------------------------------------|-------------------------|----------------------------------------------|
|`POST /households`                                |любой авторизованный     |Создать дом (создатель становится `owner`)     |
|`GET /households`                                 |любой авторизованный     |Список своих домов с ролью в каждом            |
|`GET /households/{id}`                            |участник дома            |Карточка дома                                  |
|`PATCH /households/{id}`                          |`owner`, `admin`         |Переименовать дом                              |
|`DELETE /households/{id}`                         |`owner`                  |Удалить дом целиком                            |
|`GET /households/{id}/members`                    |участник дома            |Список участников с ролями                     |
|`DELETE /households/{id}/members/{user_id}`       |см. ниже                 |Убрать участника / выйти самому                |
|`PATCH /households/{id}/members/{user_id}`        |`owner`                  |Сменить роль участника или передать владение    |
|`POST /households/{id}/invites`                   |`owner`, `admin`         |Создать инвайт-код                             |
|`GET /households/{id}/invites`                    |`owner`, `admin`         |Список инвайтов                                |
|`DELETE /households/{id}/invites/{invite_id}`     |`owner`, `admin`         |Отозвать инвайт                                |
|`POST /households/invites/{code}/accept`          |любой авторизованный     |Принять приглашение по коду → стать `member`   |

Если дома с указанным `id` не существует **или** юзер в нём не состоит — везде отдаётся один и тот же `404 Not Found`, а не `403` — иначе по разнице кодов можно было бы перебором узнавать существование чужих домов.

#### `POST /api/v1/households`

**Request body**: `{"name": "Моя квартира"}` (1–100 символов)

**Response `201 Created`**

```json
{
  "id": "e5c1f2b0-...-uuid",
  "version": 1,
  "name": "Моя квартира",
  "role": "owner",
  "created_at": "2026-08-20T10:00:00Z",
  "updated_at": "2026-08-20T10:00:00Z"
}
```

#### `GET /api/v1/households`

**Response `200 OK`** — массив объектов в том же формате, что у `POST`, каждый со своей `role`. Пустой список → `[]`, не `null`.

#### `GET /api/v1/households/{id}`

Тот же формат объекта, что выше. Доступно только участнику дома.

#### `PATCH /api/v1/households/{id}`

**Request body**: `{"name": "Новое название"}`. Конфликт версии → `409`. Роль ниже `admin` → `403`.

#### `DELETE /api/v1/households/{id}`

Только `owner`. Каскадом удаляются `household_members` и `household_invites`. **Response `204`**.

#### `GET /api/v1/households/{id}/members`

**Response `200 OK`**

```json
[
  {
    "user_id": "...",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "owner",
    "joined_at": "2026-08-20T10:00:00Z"
  }
]
```

#### `DELETE /api/v1/households/{id}/members/{user_id}`

Покрывает и «выгнать другого», и «выйти самому» (когда `user_id` == себе):

- Выход из дома разрешён всем, **кроме** единственного `owner` — сначала нужно передать владение (`PATCH .../members/{id}` с `role: owner`) или удалить дом целиком (`409 Conflict`, если попытаться).
- `admin` может убрать только `member`, не другого `admin` и не `owner` (`403`).
- `owner` может убрать любого.

**Response `204`**.

#### `PATCH /api/v1/households/{id}/members/{user_id}`

Только `owner`. **Request body**: `{"role": "admin" | "member" | "owner"}`.

- `role: admin` / `role: member` — обычная смена роли.
- `role: owner` — **передача владения**: атомарно (одним SQL-запросом) caller становится `admin`, а target — новым `owner`. Нельзя применить к себе (`400`).

**Response `204`**.

#### `POST /api/v1/households/{id}/invites`

`owner`/`admin`. **Request body**: `{"max_uses": 5}` (необязательно, лимит использований). Срок жизни инвайта фиксирован — 7 дней.

**Response `201 Created`**

```json
{
  "id": "...",
  "code": "K3F9XQPA",
  "expires_at": "2026-08-27T10:00:00Z",
  "max_uses": 5,
  "use_count": 0,
  "revoked_at": null,
  "created_at": "2026-08-20T10:00:00Z"
}
```

#### `GET /api/v1/households/{id}/invites`

`owner`/`admin`. **Response `200 OK`** — массив инвайтов (включая уже отозванные/исчерпанные — по `revoked_at`/`use_count` видно их состояние).

#### `DELETE /api/v1/households/{id}/invites/{invite_id}`

`owner`/`admin`. **Response `204`**.

#### `POST /api/v1/households/invites/{code}/accept`

Единственная ручка в фиче, не требующая членства в доме — только валидный access-токен. Атомарно проверяет, что код не отозван, не истёк, не исчерпан по `max_uses`, и что юзер ещё не участник — и добавляет его как `member`.

**Response `200 OK`** — карточка дома, `role: "member"`. Невалидный/истёкший/исчерпанный/уже принятый код → `400 Bad Request`.

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

Индекс `idx_refresh_tokens_user_id` — по `user_id`.

Таблицы `households`, `household_members`, `household_invites` (миграция `000003_households`):

**`households`**

|Колонка      |Тип           |Ограничения                                |
|-------------|--------------|--------------------------------------------|
|`id`         |`UUID`        |`PRIMARY KEY`, `DEFAULT uuid_generate_v4()`|
|`version`    |`INT`         |`NOT NULL DEFAULT 1`                       |
|`name`       |`VARCHAR(100)`|`NOT NULL`, длина 1–100                    |
|`created_at` |`TIMESTAMPTZ` |`NOT NULL DEFAULT now()`                   |
|`updated_at` |`TIMESTAMPTZ` |`NOT NULL DEFAULT now()`, `CHECK(created_at <= updated_at)`|

**`household_members`** — владелец дома не отдельная колонка, а участник с `role = 'owner'`; единственность `owner` в доме поддерживается атомарностью `TransferOwnership`, а не схемой.

|Колонка        |Тип           |Ограничения                                          |
|---------------|--------------|-------------------------------------------------------|
|`id`           |`UUID`        |`PRIMARY KEY`, `DEFAULT uuid_generate_v4()`            |
|`household_id` |`UUID`        |`NOT NULL`, `REFERENCES households(id) ON DELETE CASCADE`|
|`user_id`      |`UUID`        |`NOT NULL`, `REFERENCES users(id) ON DELETE CASCADE`   |
|`role`         |`VARCHAR(20)` |`NOT NULL DEFAULT 'member'`, `CHECK IN ('owner','admin','member')`|
|`joined_at`    |`TIMESTAMPTZ` |`NOT NULL DEFAULT now()`                               |
|                |              |`UNIQUE(household_id, user_id)`                        |

Индексы: `idx_household_members_user_id`, `idx_household_members_household_id`.

**`household_invites`**

|Колонка       |Тип           |Ограничения                                              |
|--------------|--------------|------------------------------------------------------------|
|`id`          |`UUID`        |`PRIMARY KEY`, `DEFAULT uuid_generate_v4()`                |
|`household_id`|`UUID`        |`NOT NULL`, `REFERENCES households(id) ON DELETE CASCADE`  |
|`code`        |`VARCHAR(32)` |`NOT NULL UNIQUE` (глобально, не per-household)             |
|`created_by`  |`UUID`        |`NOT NULL`, `REFERENCES users(id) ON DELETE CASCADE`        |
|`expires_at`  |`TIMESTAMPTZ` |`NOT NULL`                                                   |
|`max_uses`    |`INT`         |nullable — `NULL` = безлимитный                              |
|`use_count`   |`INT`         |`NOT NULL DEFAULT 0`                                         |
|`revoked_at`  |`TIMESTAMPTZ` |nullable                                                     |
|`created_at`  |`TIMESTAMPTZ` |`NOT NULL DEFAULT now()`                                     |
|              |              |`CHECK(max_uses IS NULL OR use_count <= max_uses)` — защита от гонки при одновременном использовании последнего инвайта|

Индекс `idx_household_invites_household_id`.

-----

## Логирование

Логи пишутся через `zap` в структурированном виде и складываются в `LOGGER_FOLDER` отдельным файлом на каждый запуск сервиса (имя файла — таймстемп старта). Уровень регулируется переменной `LOGGER_LEVEL`. Каждый HTTP-запрос получает `request_id`, который прокидывается через middleware и попадает в логи для сквозной трассировки.

-----

## Roadmap

Проект в активной разработке.

**Авторизация, пользователи, дома — базовый CRUD**

- [x] JWT access-токены + login/refresh/logout
- [x] Refresh-токены с ротацией и server-side отзывом
- [x] `Authenticate` middleware на защищённых роутах
- [x] `/users/me` — GET/PATCH/DELETE
- [x] `/households` — создание, список, карточка, переименование, удаление
- [x] Участники домов — список, удаление/выход, смена роли, передача владения
- [x] Инвайты — создание, список, отзыв, приём по коду