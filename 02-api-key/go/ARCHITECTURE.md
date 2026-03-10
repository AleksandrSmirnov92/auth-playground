# Архитектура: 02 — API Key

## Содержание

1. [Обзор проекта](#обзор-проекта)
2. [API Key: как это работает](#api-key-как-это-работает)
3. [Точка входа: main.go](#точка-входа-maingo)
4. [Clean Architecture](#clean-architecture)
5. [Domain Layer (Доменный слой)](#domain-layer-доменный-слой)
6. [Repository Layer (Слой данных)](#repository-layer-слой-данных)
7. [Use Case Layer (Бизнес-логика)](#use-case-layer-бизнес-логика)
8. [Delivery Layer (HTTP и Middleware)](#delivery-layer-http-и-middleware)
9. [Полный Flow запроса](#полный-flow-запроса)
10. [Зависимости между слоями](#зависимости-между-слоями)
11. [Преимущества архитектуры](#преимущества-архитектуры)
12. [Резюме](#резюме)

---

## Обзор проекта

Этот мини-проект реализует авторизацию по **API Key**: клиент передаёт секретный ключ в заголовке (например `X-API-Key`). Ключ выдаётся при регистрации или через отдельный эндпоинт; сервер по ключу находит пользователя и разрешает доступ. Типично используется для сервис-сервисного взаимодействия и доступа к API без передачи пароля при каждом запросе.

**Статус:** планируемая реализация (плейсхолдеры).

### Планируемая структура файлов

```
02-api-key/go/
├── cmd/server/main.go
├── internal/
│   ├── domain/
│   │   ├── user.go              # User + поле ApiKey (или отдельная сущность)
│   │   └── repository.go        # UserRepository, GetByAPIKey
│   ├── repository/memory/
│   │   └── user_repository.go
│   ├── usecase/
│   │   └── auth_usecase.go      # Register (генерация ключа), GetUserByAPIKey
│   └── delivery/
│       ├── auth_handler.go
│       └── middleware/
│           └── auth.go          # Проверка X-API-Key
├── go.mod
└── go.sum
```

---

## API Key: как это работает

Клиент при каждом запросе к защищённому API передаёт ключ в заголовке:

```
X-API-Key: <secret_key>
```

Ключ генерируется сервером при регистрации (или по отдельному запросу) — криптографически случайная строка (например `crypto/rand` или UUID). Сервер хранит связку key → user; middleware читает заголовок, ищет пользователя по ключу и кладёт user_id в context. Отдельного «логина» нет — авторизация при каждом запросе по ключу.

**Плюсы:** простота для машинного доступа, легко отозвать ключ (удалить/перегенерировать). **Минусы:** ключ в каждом запросе — только по HTTPS; компрометация ключа даёт полный доступ.

---

## Точка входа: main.go

Планируется: инициализация Repository, UseCase, Handler, **API Key middleware** (принимает use case или repository для поиска по ключу). Публичные роуты: `/health`, `POST /api/v1/auth/register` (возврат user + api_key в ответе). Защищённые роуты: `GET /api/v1/auth/me`, `DELETE /api/v1/auth/me` — обёрнуты в middleware, проверяющий `X-API-Key`. Сервер и Graceful Shutdown — по тому же шаблону, что и в 01-basic-auth.

---

## Clean Architecture

Те же слои: Domain, Use Case, Repository, Delivery. В Delivery добавляется middleware, который по заголовку `X-API-Key` вызывает use case (или repository) для получения пользователя по ключу и передачи user_id в context. Бизнес-логика не знает про HTTP и заголовки.

---

## Domain Layer (Доменный слой)

**User:** поля ID, Email, ApiKey (или хеш ключа), CreatedAt. ApiKey хранится в одном экземпляре (или только хеш для сравнения).

**UserRepository (планируется):** Create, GetByID, GetByEmail, **GetByAPIKey(apiKey string) (*User, error)**, Delete. GetByAPIKey — основа для middleware.

---

## Repository Layer (Слой данных)

In-memory: map по ID; для GetByAPIKey — перебор по пользователям или второй map `apiKeyToUserID`. С мьютексом для конкурентного доступа. В production — индекс/поля в БД по api_key.

---

## Use Case Layer (Бизнес-логика)

**Register(email, password):** проверка существования по email, генерация API key (например 32 байта `crypto/rand`, hex/base64), хеш пароля (bcrypt), сохранение user с api_key. Возврат user + plain api_key клиенту один раз (дальше ключ только в заголовках).

**GetUserByAPIKey(apiKey string):** поиск пользователя по ключу, возврат user или error. Вызывается из middleware.

**GetUserByID, DeleteUserById** — по аналогии с 01-basic-auth.

---

## Delivery Layer (HTTP и Middleware)

**auth_handler.go:** RegisterHandler (JSON email/password, в ответе user + api_key); MeHandler и DeleteUserHandler берут user_id из context (установленный middleware).

**middleware/auth.go (планируется):** чтение заголовка `X-API-Key`, вызов GetUserByAPIKey, при отсутствии/неверном ключе — 401; при успехе — context.WithValue(UserIDKey, user.ID) и вызов next.ServeHTTP.

---

## Полный Flow запроса

Клиент регистрируется → получает api_key в теле ответа. Дальше при каждом запросе к защищённому эндпоинту отправляет заголовок `X-API-Key`. Middleware извлекает ключ → GetUserByAPIKey → user_id в context → handler возвращает данные пользователя или 204 при удалении.

---

## Зависимости между слоями

Как в 01-basic-auth: HTTP → Middleware → Handler → UseCase → Domain; Repository реализует интерфейсы Domain. Use Case не зависит от способа передачи ключа (заголовок).

---

## Преимущества архитектуры

Тестируемость (мок репозитория), замена хранилища в одном месте, чёткое разделение: middleware только извлекает ключ и проверяет его через use case.

---

## Резюме

**Планируется:** авторизация по заголовку X-API-Key, генерация ключа при регистрации, middleware по ключу, те же принципы Clean Architecture. После реализации сюда подставляется актуальный код и диаграммы.

**Следующий проект (Go):** [03-session-cookie](../../03-session-cookie/go/ARCHITECTURE.md) — сессионная авторизация через Cookie.
