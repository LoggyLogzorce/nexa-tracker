[English](README.md) | [Русский](README.ru.md)

# Nexa Task Tracker

API для отслеживания задач на Go, Gin, PostgreSQL и GORM.

## Возможности

### Реализовано
- **JWT аутентификация** — регистрация, вход, выход, обновление токена с ротацией и обнаружением повторного использования
- **OAuth 2.0 вход** — интеграция с Google и Yandex
- **Управление сессиями** — просмотр активных сессий, отзыв конкретных сессий
- **Управление пользователями** — получение/обновление профиля, загрузка аватара, смена пароля, удаление аккаунта с анонимизацией данных, поиск пользователей
- **Управление проектами** — полный CRUD, фильтрация (собственные/участие), поиск
- **Управление задачами** — полный CRUD с валидацией (исполнитель/статус/приоритет в рамках проекта), архивация/разархивация, глобальный поиск, задачи текущего пользователя (назначенные/созданные), история изменений через JSONB diffs
- **Свои статусы и приоритеты** — для каждого проекта, с поддержкой drag-reorder (`order_index`), значения по умолчанию (To Do / In Progress / Done; Low / Medium / High)
- **Участники проекта** — ролевой доступ (owner / member / read_only)
- **Комментарии к задачам** — CRUD с проверкой владения, обогащением данных пользователя
- **Файловые вложения** — загрузка, скачивание, удаление с отслеживанием метаданных
- **Шина событий** — синхронная pub/sub для межмодульного взаимодействия (например, создание проекта создаёт статусы и приоритеты по умолчанию, удаление пользователя каскадирует)
- **Ограничение запросов** — по IP, настраиваемое (auth endpoints: 5 запросов/мин)
- **Стандартизированные JSON ответы** — единый формат `{"success": bool, "data": ..., "error": ...}`
- **Версия и Changelog** — информация о сборке и заметки о релизах

### В разработке / Заглушки
- **2FA (TOTP)** — эндпоинты подключены, обработчики возвращают заглушки
- **Модуль уведомлений** — подготовлен, `Init()` закомментирован в `main.go`

### Запланировано
- **Чат** — файл-заглушка

## Технологии

| Уровень | Технология |
|---|---|
| Язык | Go 1.26.1 |
| HTTP фреймворк | Gin v1.12.0 |
| ORM | GORM v1.30.0 (с datatypes для JSONB) |
| База данных | PostgreSQL 16 |
| JWT | golang-jwt v5.2.1 |
| OAuth 2.0 | golang.org/x/oauth2 |
| Хеширование паролей | bcrypt (x/crypto v0.48.0) |
| UUID | google/uuid v1.6.0 |
| Rate Limiting | x/time v0.5.0 |
| Загрузка .env | godotenv v1.5.1 |
| Валидация | go-playground/validator v10.30.1 |
| CORS | gin-contrib/cors |
| Контейнеризация | Docker + Docker Compose |
| Фронтенд | Vite + React (отдельное SPA) |

## Архитектура

```
┌─────────────┐     ┌──────────────┐     ┌────────────────┐
│   Handler   │────▶│   Service    │────▶│  Repository    │
│ (HTTP only) │     │ (бизнес-     │     │ (GORM запросы) │
│             │     │  логика)     │     │                │
└─────────────┘     └──────┬───────┘     └────────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │  Event Bus   │
                    │ (pub/sub)    │
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
         status svc   priority svc   notify svc
```

- **Handler** — привязывает запросы, извлекает контекст, валидирует ввод, возвращает ответы
- **Service** — бизнес-логика, управление таймаутами (5-10s контекстные дедлайны)
- **Repository** — GORM запросы, транзакционные записи
- **Event Bus** — синхронная pub/sub для слабосвязанного межмодульного взаимодействия
- **Middleware стэк** — ограничитель запросов → лимит тела запроса → JWT auth → RBAC (уровень проекта)

## Структура проекта

```
├── cmd/app/main.go                 — Точка входа
├── docker/
│   ├── Dockerfile                  — Многоступенчатая сборка Go
│   ├── Dockerfile.frontend         — Сборка Node/Vite
│   └── docker-compose.yml          — App + PostgreSQL 16 + Frontend
├── internal/
│   ├── api/router.go               — Определение маршрутов
│   ├── config/
│   │   ├── config.go               — Конфигурация из переменных окружения
│   │   ├── google_conf.go          — Конфиг Google OAuth
│   │   └── yandex_conf.go          — Конфиг Yandex OAuth
│   ├── ctxkeys/ctxkeys.go          — Константы ключей контекста
│   ├── db/
│   │   ├── db.go                   — GORM подключение + авто-миграция
│   │   └── schema.sql              — SQL схема для справки
│   ├── middleware/
│   │   ├── auth.go                 — JWT аутентификация
│   │   ├── rbac.go                 — RBAC на уровне проекта (owner/member/read_only)
│   │   ├── ratelimit.go            — IP-based ограничитель запросов
│   │   └── requestbody.go          — Лимит тела запроса
│   ├── models/                     — GORM модели
│   │   ├── user.go
│   │   ├── user_provider.go        — Привязка OAuth провайдеров
│   │   ├── auth.go                 — Модель RefreshToken
│   │   ├── project.go
│   │   ├── participant.go
│   │   ├── status.go
│   │   ├── priority.go
│   │   ├── task.go                 — Task + UpdateHistory
│   │   ├── comment.go
│   │   └── attachment.go
│   ├── core/
│   │   ├── auth/
│   │   │   ├── handler.go          — Register, Login, Refresh, Logout, Sessions, ChangePassword, 2FA заглушки
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── model.go
│   │   │   ├── errors.go
│   │   │   └── oauth/
│   │   │       ├── handler.go      — Google/Yandex login & callback
│   │   │       ├── google_service.go
│   │   │       ├── yandex_service.go
│   │   │       └── ctxkeys.go
│   │   ├── user/                   — CRUD пользователей, аватар, поиск, удаление
│   │   ├── project/                — CRUD проектов, поиск
│   │   ├── task/                   — CRUD задач, архив, история, поиск
│   │   ├── status/                 — Статусы проектов с сортировкой
│   │   ├── priority/               — Приоритеты проектов
│   │   ├── participant/            — Участники и роли
│   │   ├── comment/                — Комментарии к задачам
│   │   └── attachment/             — Загрузка и скачивание файлов
│   ├── modules/
│   │   ├── notify/                 — Уведомления (заглушка, закомментировано)
│   │   └── chat/chat.go            — Чат (заглушка)
│   └── version/
│       ├── version.go              — Константа версии + встроенный changelog
│       └── changelog.md            — Заметки о релизах
├── pkg/
│   ├── cookie/                     — Хелперы установки/удаления кук
│   ├── events/                     — Шина событий + структуры событий
│   ├── hash/                       — bcrypt + SHA-256 хеширование токенов
│   ├── jwt/                        — Access/refresh JWT токены
│   ├── nullable/                   — Nullable типы для PATCH обновлений
│   ├── response/                   — Стандартизированные JSON ответы
│   └── validation/                 — Хелперы валидации запросов
├── frontend/                       — Vite + React SPA
├── uploads/                        — Директория загрузок (в .gitignore)
│   ├── avatars/
│   └── {task_id}/
├── .env.example                    — Шаблон переменных окружения
└── .gitignore
```

## Начало работы

### Требования

- Go 1.26+ (или Docker)
- PostgreSQL 16

### Локальная разработка

```bash
# Настроить окружение
cp .env.example .env
# Отредактируйте .env: данные БД, JWT секрет, OAuth ключи

# Запуск
go mod tidy
go run ./cmd/app
```

### Через Docker

```bash
docker compose -f docker/docker-compose.yml up --build
```

Сервер запускается на `http://localhost:8080`. Фронтенд (собранный) работает на порту `3000`.

## Конфигурация

Вся конфигурация через переменные окружения (см. `.env.example`):

| Переменная | По умолчанию | Описание |
|---|---|---|
| `ENV` | `development` | Режим окружения |
| `SERVER_HOST` | `0.0.0.0` | Адрес привязки |
| `SERVER_PORT` | `8080` | HTTP порт |
| `DB_HOST` | `127.0.0.1` | Хост PostgreSQL |
| `DB_PORT` | `5432` | Порт PostgreSQL |
| `DB_USER` | `postgres` | Пользователь БД |
| `DB_PASSWORD` | — | Пароль БД |
| `DB_NAME` | `nexa_tracker` | Имя БД |
| `DB_SSLMODE` | `disable` | SSL режим PostgreSQL |
| `JWT_SECRET` | — | HMAC ключ подписи |
| `JWT_ACCESS_EXPIRY` | `15m` | Время жизни access токена |
| `JWT_REFRESH_EXPIRY` | `168h` (7d) | Время жизни refresh токена |
| `UPLOAD_PATH` | `./uploads` | Директория загрузки файлов |
| `COOKIE_DOMAIN` | — | Домен куки |
| `COOKIE_SECURE` | `false` | Флаг Secure для куки |
| `COOKIE_SAMESITE` | `1` | SameSite (1=Default, 2=Lax, 3=Strict, 4=None) |
| `CORS_ORIGINS` | `http://localhost:5173` | CORS источники (через запятую) |
| `FRONTEND_URL` | `http://localhost:5173` | URL фронтенда для OAuth редиректов |
| `NOTIFY_MODULE` | `false` | Включить модуль уведомлений |
| `GOOGLE_CLIENT_ID` | — | Google OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | — | Google OAuth client secret |
| `GOOGLE_REDIRECT_URL` | — | Google OAuth redirect URL |
| `YANDEX_CLIENT_ID` | — | Yandex OAuth client ID |
| `YANDEX_CLIENT_SECRET` | — | Yandex OAuth client secret |
| `YANDEX_REDIRECT_URL` | — | Yandex OAuth redirect URL |

## API эндпоинты

### Health & Meta

| Метод | Путь | Auth |
|---|---|---|
| `GET` | `/health` | — |
| `GET` | `/api/v1/version` | — |
| `GET` | `/api/v1/changelog` | — |

### Auth (`/api/v1/auth`)

Ограничение запросов (5 запросов/мин). Для хранения токенов используются куки.

| Метод | Путь | Статус |
|---|---|---|
| `POST` | `/register` | ✅ |
| `POST` | `/login` | ✅ |
| `POST` | `/refresh` | ✅ |
| `POST` | `/logout` | ✅ |
| `GET` | `/google/login` | ✅ |
| `GET` | `/google/callback` | ✅ |
| `GET` | `/yandex/login` | ✅ |
| `GET` | `/yandex/callback` | ✅ |
| `POST` | `/2fa/setup` | 🚧 заглушка |
| `POST` | `/2fa/verify` | 🚧 заглушка |
| `POST` | `/2fa/enable` | 🚧 заглушка |
| `POST` | `/2fa/disable` | 🚧 заглушка |

### Users (`/api/v1/users`)

Требуется JWT auth.

| Метод | Путь | Статус |
|---|---|---|
| `GET` | `/me` | ✅ |
| `PUT` | `/me` | ✅ |
| `PUT` | `/me/avatar` | ✅ |
| `PUT` | `/me/change-password` | ✅ |
| `DELETE` | `/me` | ✅ |
| `GET` | `/me/sessions` | ✅ |
| `DELETE` | `/me/sessions/:id` | ✅ |
| `GET` | `/search?q=` | ✅ |

### Projects (`/api/v1/projects`)

Требуется JWT auth. Роли участников: `owner` (чтение+запись+удаление), `member` (чтение+запись), `read_only` (чтение).

| Метод | Путь | Доступ | Статус |
|---|---|---|---|
| `GET` | `/` | authenticated | ✅ |
| `GET` | `/owned` | authenticated | ✅ |
| `GET` | `/search?q=` | authenticated | ✅ |
| `POST` | `/` | authenticated | ✅ |
| `GET` | `/:id` | read_only+ | ✅ |
| `PUT` | `/:id` | owner | ✅ |
| `DELETE` | `/:id` | owner | ✅ |
| `GET` | `/:id/participants` | read_only+ | ✅ |
| `POST` | `/:id/participants` | owner | ✅ |
| `PUT` | `/:id/participants/:user_id` | owner | ✅ |
| `DELETE` | `/:id/participants/:user_id` | owner | ✅ |
| `GET` | `/:id/statuses` | read_only+ | ✅ |
| `POST` | `/:id/statuses` | member+ | ✅ |
| `PUT` | `/:id/statuses/:status_id` | member+ | ✅ |
| `DELETE` | `/:id/statuses/:status_id` | owner | ✅ |
| `GET` | `/:id/priorities` | read_only+ | ✅ |
| `POST` | `/:id/priorities` | member+ | ✅ |
| `PUT` | `/:id/priorities/:priority_id` | member+ | ✅ |
| `DELETE` | `/:id/priorities/:priority_id` | owner | ✅ |
| `GET` | `/:id/attachments` | read_only+ | ✅ |

### Tasks (`/api/v1/projects/:id/tasks`)

Требуется JWT auth + доступ к проекту. Поддерживает фильтр `?archived=true`. История отслеживается в формате JSONB diffs.

| Метод | Путь | Доступ | Статус |
|---|---|---|---|
| `GET` | `/?archived=` | read_only+ | ✅ |
| `POST` | `/` | member+ | ✅ |
| `GET` | `/:task_id?archived=` | read_only+ | ✅ |
| `PUT` | `/:task_id?archived=` | member+ | ✅ |
| `DELETE` | `/:task_id` | owner | ✅ |
| `GET` | `/:task_id/history` | read_only+ | ✅ |
| `GET` | `/:task_id/comments` | member+ | ✅ |
| `POST` | `/:task_id/comments` | member+ | ✅ |
| `PUT` | `/:task_id/comments/:comment_id` | member+ | ✅ |
| `DELETE` | `/:task_id/comments/:comment_id` | member+ | ✅ |
| `GET` | `/:task_id/attachments` | read_only+ | ✅ |
| `GET` | `/:task_id/attachments/:attachment_id` | read_only+ | ✅ |
| `POST` | `/:task_id/attachments` | member+ | ✅ |
| `DELETE` | `/:task_id/attachments/:attachment_id` | member+ | ✅ |

### Глобальный поиск задач (`/api/v1/tasks`)

Требуется JWT auth.

| Метод | Путь | Статус |
|---|---|---|
| `GET` | `/me?type=assigned\|reported` | ✅ |
| `GET` | `/search?q=` | ✅ |

### Notifications (`/api/v1/notifications`)

Закомментировано в `main.go`. По умолчанию недоступно.

## Схема БД

### Таблицы

| Таблица | Описание |
|---|---|
| `users` | Учётные записи пользователей (uuid PK, email, password_hash, name, role, 2fa secret) |
| `user_providers` | Привязка OAuth провайдеров (Google, Yandex) |
| `refresh_tokens` | Хранение JWT refresh токенов с отслеживанием отзыва и обнаружением повторного использования |
| `projects` | Проекты, принадлежащие пользователю |
| `project_participants` | Многие-ко-многим с ролями (owner/member/read_only) |
| `statuses` | Статусы задач для каждого проекта с `order_index` для drag-reorder |
| `priorities` | Приоритеты задач для каждого проекта |
| `tasks` | Задачи со ссылками на проект, статус, приоритет, исполнителя, автора |
| `update_history` | JSONB-based отслеживание изменений на уровне полей |
| `comments` | Комментарии в рамках задачи |
| `attachments` | Метаданные файлов (имя, путь, размер, mime тип) |

Полную схему см. в `internal/db/schema.sql`.

## Разработка

### Добавление нового модуля

1. Создать `internal/core/<module>/` с `model.go`, `handler.go`, `service.go`, `repository.go`, `errors.go`
2. Реализовать интерфейс `Repository` через GORM
3. Реализовать бизнес-логику в `Service`
4. Подключить HTTP обработчики в `Handler`
5. Зарегистрировать маршруты в `internal/api/router.go`
6. Инициализировать в `cmd/app/main.go`

### Стиль кода

- Трёхуровневая архитектура: Handler → Service → Repository
- Сервисы используют `context.WithTimeout` для всех операций с БД
- Используйте `pkg/nullable` типы для PATCH эндпоинтов, чтобы различать "не передано" и "явный null"
- Используйте `pkg/response` хелперы для единого форматирования JSON
- События для межмодульного взаимодействия проходят через синхронный EventBus

## Лицензия

MIT
