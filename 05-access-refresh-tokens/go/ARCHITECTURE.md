# Архитектура: 05 — Access + Refresh Tokens

## Содержание

1. [Обзор проекта](#обзор-проекта)
2. [Access + Refresh: как это работает](#access--refresh-как-это-работает)
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

Этот мини-проект комбинирует **короткоживущий Access Token** (JWT) для доступа к API и **долгоживущий Refresh Token** (opaque, хранится на сервере) для получения новой пары токенов без повторного ввода пароля. Клиент использует access token в заголовке Bearer; при истечении access — отправляет refresh token на POST /refresh и получает новую пару.

**Статус:** планируемая реализация (плейсхолдеры).

### Планируемая структура файлов

```
05-access-refresh-tokens/go/
├── cmd/server/main.go
├── internal/
│   ├── domain/
│   │   ├── user.go
│   │   ├── refresh_token.go     # или храним в репозитории по token_id
│   │   └── repository.go        # UserRepository, RefreshTokenRepository
│   ├── repository/memory/
│   │   ├── user_repository.go
│   │   └── refresh_token_repository.go
│   ├── usecase/
│   │   └── auth_usecase.go      # Register, Login (пара токенов), Refresh, GetUserByID
│   └── delivery/
│       ├── auth_handler.go      # Register, Login, Refresh, Me, Delete
│       └── middleware/
│           └── auth.go          # Только Access JWT
├── go.mod
└── go.sum
```

---

## Access + Refresh: как это работает

1. **Login:** сервер проверяет email/password, выдаёт **access token** (JWT, короткий TTL, например 15 мин) и **refresh token** (случайная строка, долгий TTL, например 7 дней). Refresh token сохраняется на сервере (memory/redis/DB) в связке с user_id.
2. **Доступ к API:** клиент отправляет только access token в `Authorization: Bearer <access>`. Middleware проверяет JWT, извлекает user_id, передаёт в context.
3. **Обновление пары:** когда access истёк, клиент отправляет POST /refresh с телом { "refresh_token": "..." }. Сервер проверяет refresh token по хранилищу, находит user_id, выдаёт новую пару access + refresh (старый refresh можно инвалидировать или оставить одноразовым).

**Плюсы:** короткий срок жизни access ограничивает ущерб при утечке; refresh даёт удобство без частого логина. **Минусы:** нужна безопасная передача и хранение refresh token (HttpOnly cookie или защищённое хранилище на клиенте).

---

## Точка входа: main.go

Планируется: UserRepository, RefreshTokenRepository, AuthUsecase (JWT + refresh логика), Handler, **JWT middleware только для access token**. Публичные: `/health`, `POST /register`, `POST /login`, `POST /refresh`. Защищённые: `GET /me`, `DELETE /me` — обёрнуты в access-token middleware. Сервер и Graceful Shutdown — по шаблону.

---

## Clean Architecture

Domain: User, сущность/хранилище для refresh token (token_id → user_id, expires). Use Case: Register, Login (генерация access JWT + refresh, сохранение refresh), Refresh (проверка refresh, выдача новой пары), GetUserByID. Delivery: handlers + middleware по access JWT; refresh обрабатывается только в RefreshHandler.

---

## Domain Layer (Доменный слой)

**User:** как в 01-basic-auth.

**RefreshToken (или хранилище):** идентификатор токена (или сам хеш токена), user_id, expires_at. Репозиторий: Save(refreshToken), GetByToken(token) → user_id или error, Revoke(token) при необходимости.

**UserRepository:** Create, GetByID, GetByEmail, Delete.

**RefreshTokenRepository (планируется):** Create(tokenID, userID, expiresAt), GetUserByToken(token) (*User, error), Delete(tokenID) для одноразового использования или logout.

---

## Repository Layer (Слой данных)

**user_repository:** in-memory, как раньше.

**refresh_token_repository:** in-memory map token_id → {user_id, expires_at} с мьютексом. При Refresh можно удалять использованный refresh (rotation) или оставлять до истечения.

---

## Use Case Layer (Бизнес-логика)

**Register:** как в 01-basic-auth.

**Login(email, password):** проверка пароля; генерация access JWT (короткий TTL); генерация refresh token (crypto/rand), сохранение в RefreshTokenRepository; возврат пары { access_token, refresh_token, expires_in }.

**Refresh(refreshToken string):** поиск по refresh token в репозитории, проверка срока действия; генерация новой пары access + refresh; при rotation — удаление старого refresh; возврат новой пары.

**GetUserByID:** как в 01-basic-auth. Middleware по access JWT передаёт user_id в context.

---

## Delivery Layer (HTTP и Middleware)

**auth_handler.go:** RegisterHandler; LoginHandler — возврат JSON с access_token и refresh_token; RefreshHandler — приём refresh_token из body, вызов UseCase.Refresh, возврат новой пары; MeHandler, DeleteUserHandler — user_id из context.

**middleware/auth.go:** только проверка Access JWT (как в 04-jwt-bearer). Refresh token в middleware не используется — только в POST /refresh.

---

## Полный Flow запроса

**Логин:** POST /login → access + refresh в ответе. **Запрос к /me:** Authorization: Bearer <access> → middleware проверяет JWT → user_id в context → MeHandler. **Истечение access:** клиент вызывает POST /refresh { "refresh_token": "..." } → сервер проверяет refresh, выдаёт новую пару → клиент сохраняет и использует новый access.

---

## Зависимости между слоями

Как в 04: middleware только для access JWT; refresh-логика в use case и RefreshHandler. Domain не знает о JWT — только о хранении refresh token и пользователей.

---

## Преимущества архитектуры

Разделение ответственности: access для доступа, refresh для обновления; отзыв refresh при logout или компрометации; тестируемость use case и middleware по отдельности.

---

## Резюме

**Планируется:** короткоживущий access JWT, долгоживущий refresh token с хранением на сервере, endpoint POST /refresh, middleware только по access token. После реализации — подставить актуальный код и диаграммы.

**Следующий проект (Go):** [06-oauth2](../../06-oauth2/go/ARCHITECTURE.md) — OAuth 2.0 (Google/GitHub).
