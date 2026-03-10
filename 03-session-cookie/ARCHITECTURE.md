# Архитектура: 03 — Session + Cookie

## Содержание

1. [Обзор проекта](#обзор-проекта)
2. [Session + Cookie: как это работает](#session--cookie-как-это-работает)
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

Этот мини-проект реализует **серверную сессионную авторизацию**: после успешного логина сервер создаёт сессию (session_id → user_id), записывает session_id в **Cookie**, при последующих запросах клиент отправляет cookie автоматически, middleware по session_id находит пользователя и передаёт user_id в context. Logout — удаление сессии на сервере и инвалидация cookie.

**Статус:** планируемая реализация (плейсхолдеры).

### Планируемая структура файлов

```
03-session-cookie/
├── cmd/server/main.go
├── internal/
│   ├── domain/
│   │   ├── user.go
│   │   ├── session.go           # Session: ID, UserID, ExpiresAt (опционально)
│   │   └── repository.go       # UserRepository, SessionRepository
│   ├── repository/memory/
│   │   ├── user_repository.go
│   │   └── session_repository.go
│   ├── usecase/
│   │   └── auth_usecase.go      # Register, Login (создание сессии), Logout, GetUserByID
│   └── delivery/
│       ├── auth_handler.go      # Register, Login (Set-Cookie), Logout, Me, Delete
│       └── middleware/
│           └── auth.go           # Чтение cookie session_id, проверка сессии
├── go.mod
└── go.sum
```

---

## Session + Cookie: как это работает

1. **Login:** клиент отправляет email/password; сервер проверяет, создаёт сессию (уникальный session_id), сохраняет в хранилище (memory/redis), в ответе ставит `Set-Cookie: session_id=...; HttpOnly; Path=/; SameSite=Strict`.
2. **Последующие запросы:** браузер автоматически отправляет cookie; middleware читает session_id, ищет сессию в хранилище, по user_id кладёт в context и вызывает handler.
3. **Logout:** клиент вызывает POST /logout; сервер удаляет сессию из хранилища и отдаёт Set-Cookie с пустым/истёкшим значением.

**Плюсы:** не передаём пароль после логина, удобно для браузеров. **Минусы:** состояние на сервере (или в Redis), нужна политика истечения сессий.

---

## Точка входа: main.go

Планируется: инициализация UserRepository, SessionRepository, AuthUsecase, AuthHandler, **Session middleware** (принимает session store / use case). Публичные роуты: `/health`, `POST /register`, `POST /login`, `POST /logout`. Защищённые: `GET /me`, `DELETE /me` — обёрнуты в session middleware. Сервер и Graceful Shutdown — по шаблону 01-basic-auth.

---

## Clean Architecture

Domain: User, Session, интерфейсы UserRepository и SessionRepository. Use Case: Register, Login (создание сессии), Logout (удаление сессии), GetUserByID. Delivery: handlers + middleware, читающий cookie и проверяющий сессию через use case или session repository.

---

## Domain Layer (Доменный слой)

**User:** как в 01-basic-auth (ID, Email, Password hash, CreatedAt).

**Session (планируется):** ID (session_id), UserID, CreatedAt; опционально ExpiresAt для TTL. Хранилище: map[session_id]Session или в Redis.

**UserRepository:** Create, GetByID, GetByEmail, Delete.

**SessionRepository (планируется):** Create(session), GetByID(sessionID), Delete(sessionID). GetByID возвращает сессию → по ней получаем UserID.

---

## Repository Layer (Слой данных)

**user_repository:** как в 01-basic-auth, in-memory map с RWMutex.

**session_repository:** in-memory map session_id → Session, мьютекс. Create генерирует уникальный session_id (crypto/rand), сохраняет; GetByID — поиск по ключу; Delete — удаление. В production — Redis с TTL.

---

## Use Case Layer (Бизнес-логика)

**Register:** как в 01-basic-auth (bcrypt, проверка email).

**Login(email, password):** проверка пароля через GetByEmail + bcrypt; при успехе — создание сессии (SessionRepository.Create), возврат user + session_id (или только session_id для установки в cookie в handler).

**Logout(sessionID):** SessionRepository.Delete(sessionID). Handler затем сбрасывает cookie.

**GetUserByID:** как в 01-basic-auth. Middleware по session_id получает user_id из сессии и кладёт в context.

---

## Delivery Layer (HTTP и Middleware)

**auth_handler.go:** RegisterHandler; LoginHandler — после успешного Login устанавливает Set-Cookie с session_id (HttpOnly, Path=/, SameSite); LogoutHandler — удаляет сессию по cookie и отдаёт Set-Cookie с пустым/истёкшим; MeHandler, DeleteUserHandler — user_id из context.

**middleware/auth.go:** чтение cookie `session_id`, вызов SessionRepository.GetByID или use case; при отсутствии/невалидной сессии — 401; при успехе — context.WithValue(UserIDKey, userID), next.ServeHTTP.

---

## Полный Flow запроса

**Логин:** POST /login {email, password} → UseCase Login → создание сессии → ответ с Set-Cookie(session_id). **Запрос к /me:** браузер отправляет Cookie(session_id) → middleware извлекает session_id → поиск сессии → user_id в context → MeHandler возвращает user. **Logout:** POST /logout с cookie → удаление сессии → Set-Cookie сброс.

---

## Зависимости между слоями

Как в 01/02: HTTP → Middleware → Handler → UseCase → Domain; Repository реализует интерфейсы. Use Case не знает о cookie — только о session_id и операциях с сессиями.

---

## Преимущества архитектуры

Тестируемость (мок репозиториев), возможность заменить in-memory сессии на Redis без изменения use case, единая точка проверки сессии в middleware.

---

## Резюме

**Планируется:** серверные сессии, cookie с session_id, Login (Set-Cookie), Logout (удаление сессии и сброс cookie), middleware по cookie. После реализации — подставить актуальный код и диаграммы.

**Следующий проект:** [04-jwt-bearer](../04-jwt-bearer/) — JWT Bearer Token (stateless).
