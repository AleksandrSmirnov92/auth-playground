# Архитектура: 04 — JWT Bearer Token

## Содержание

1. [Обзор проекта](#обзор-проекта)
2. [JWT Bearer: как это работает](#jwt-bearer-как-это-работает)
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

Этот мини-проект реализует **stateless авторизацию через JWT** (JSON Web Token): после логина сервер выдаёт подписанный токен, клиент при каждом запросе передаёт его в заголовке `Authorization: Bearer <token>`. Сервер проверяет подпись и извлекает user_id из claims, без хранения сессий на сервере.

**Статус:** планируемая реализация (плейсхолдеры).

### Планируемая структура файлов

```
04-jwt-bearer/
├── cmd/server/main.go
├── internal/
│   ├── domain/
│   │   ├── user.go
│   │   └── repository.go
│   ├── repository/memory/
│   │   └── user_repository.go
│   ├── usecase/
│   │   └── auth_usecase.go      # Register, Login (возврат JWT), GetUserByID
│   └── delivery/
│       ├── auth_handler.go      # Register, Login (выдача JWT), Me, Delete
│       └── middleware/
│           └── auth.go          # Парсинг Bearer token, верификация JWT, user_id в context
├── pkg/ или internal/           # JWT: подпись и верификация (HS256/RS256)
│   └── jwt/
├── go.mod
└── go.sum
```

---

## JWT Bearer: как это работает

**JWT** — компактный формат: три части в Base64, разделённые точкой: `header.payload.signature`. В payload (claims) хранятся данные, например `sub` (user_id), `exp` (время истечения). Подпись (HMAC-SHA256 или RSA) гарантирует, что токен не подделан и не изменён.

**Flow:** клиент отправляет POST /login {email, password} → сервер проверяет пароль, генерирует JWT с claims (user_id, exp), возвращает токен в теле ответа (или в cookie). Клиент при запросах к защищённым эндпоинтам отправляет заголовок `Authorization: Bearer <token>`. Middleware извлекает токен, проверяет подпись и exp, достаёт user_id из claims и кладёт в context.

**Плюсы:** stateless, масштабируемость без общего хранилища сессий. **Минусы:** отзыв до истечения срока сложнее (нужен blacklist или короткий TTL); секрет/ключи должны храниться безопасно.

---

## Точка входа: main.go

Планируется: инициализация Repository, UseCase (с зависимостью от JWT-сервиса: секрет, TTL), Handler, **JWT Bearer middleware** (принимает JWT-сервис для верификации). Публичные: `/health`, `POST /register`, `POST /login`. Защищённые: `GET /me`, `DELETE /me` — обёрнуты в JWT middleware. Сервер и Graceful Shutdown — по шаблону 01-basic-auth.

---

## Clean Architecture

Domain: User, UserRepository. Use Case: Register, Login (генерация JWT через внутренний сервис/пакет), GetUserByID. Delivery: handlers + middleware, который парсит Bearer token, верифицирует JWT и кладёт user_id в context. JWT-логика (подпись/верификация) — в отдельном пакете или use case, не в HTTP-слое.

---

## Domain Layer (Доменный слой)

**User:** как в 01-basic-auth (ID, Email, Password hash, CreatedAt). JWT не хранит пароль — только идентификатор пользователя в claims.

**UserRepository:** Create, GetByID, GetByEmail, Delete. Сессии не хранятся; JWT сам несёт данные.

---

## Repository Layer (Слой данных)

In-memory UserRepository, как в 01-basic-auth. Отдельного хранилища для токенов нет (stateless).

---

## Use Case Layer (Бизнес-логика)

**Register:** как в 01-basic-auth (bcrypt, проверка email).

**Login(email, password):** проверка через GetByEmail + bcrypt; при успехе — генерация JWT (claims: sub=user.ID, exp=now+TTL), возврат токена (и опционально user) в вызывающий слой. Генерация JWT может быть в use case или в отдельном пакете, вызываемом из use case.

**GetUserByID:** как в 01-basic-auth. Middleware после верификации JWT получает user_id из claims и передаёт в context; handler вызывает GetUserByID(user_id).

---

## Delivery Layer (HTTP и Middleware)

**auth_handler.go:** RegisterHandler; LoginHandler — после успешного Login возвращает в JSON токен (и при необходимости user); MeHandler, DeleteUserHandler — user_id из context.

**middleware/auth.go:** чтение заголовка `Authorization`, извлечение строки после "Bearer "; вызов JWT Verify (проверка подписи и exp); при ошибке — 401; при успехе — извлечение user_id из claims (например claim "sub"), context.WithValue(UserIDKey, userID), next.ServeHTTP.

---

## Полный Flow запроса

**Логин:** POST /login {email, password} → UseCase Login → генерация JWT → ответ { "token": "eyJ..." }. **Запрос к /me:** клиент отправляет заголовок Authorization: Bearer eyJ... → middleware парсит токен → Verify подписи и exp → извлекает user_id из claims → context → MeHandler возвращает user. Сервер не хранит токен; вся информация в самом JWT.

---

## Зависимости между слоями

HTTP → Middleware (JWT verify) → Handler → UseCase → Domain; Repository реализует Domain. Use Case не знает о заголовке Bearer — только о генерации токена при логине; middleware не знает о паролях — только о верификации JWT.

---

## Преимущества архитектуры

Stateless сервер, горизонтальное масштабирование без общего session store; тестируемость use case без HTTP; JWT-логика изолирована в пакете/middleware.

---

## Резюме

**Планируется:** JWT при логине, заголовок Authorization: Bearer, middleware верификации JWT и извлечения user_id из claims, те же принципы Clean Architecture. После реализации — подставить актуальный код и диаграммы.

**Следующий проект:** [05-access-refresh-tokens](../05-access-refresh-tokens/) — Access + Refresh токены.
