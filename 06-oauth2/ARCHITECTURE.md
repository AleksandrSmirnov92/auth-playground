# Архитектура: 06 — OAuth 2.0

## Содержание

1. [Обзор проекта](#обзор-проекта)
2. [OAuth 2.0: как это работает](#oauth-20-как-это-работает)
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

Этот мини-проект реализует **OAuth 2.0 Authorization Code flow** с внешними провайдерами (Google, GitHub): пользователь перенаправляется на страницу провайдера для входа, провайдер возвращает authorization code на наш callback URL, наш сервер обменивает code на access token провайдера и по нему получает профиль пользователя (email, id), после чего создаёт или находит пользователя в нашем приложении и выдаёт ему сессию/JWT для доступа к нашим эндпоинтам.

**Статус:** планируемая реализация (плейсхолдеры).

### Планируемая структура файлов

```
06-oauth2/
├── cmd/server/main.go
├── internal/
│   ├── domain/
│   │   ├── user.go              # User, возможно OAuthID / Provider
│   │   └── repository.go
│   ├── repository/memory/
│   │   └── user_repository.go
│   ├── usecase/
│   │   └── auth_usecase.go      # GetOrCreateUserByOAuth(provider, profile), GetUserByID
│   └── delivery/
│       ├── auth_handler.go      # Me, Delete; callback handler
│       ├── oauth_handler.go     # Redirect to provider, Callback (code → token → profile)
│       └── middleware/          # Session или JWT после OAuth login
│           └── auth.go
├── pkg/ или internal/oauth/      # Клиенты провайдеров (Google, GitHub), обмен code→token, UserInfo
├── go.mod
└── go.sum
```

---

## OAuth 2.0: как это работает

**Роли:** Resource Owner (пользователь), Client (наше приложение), Authorization Server (Google/GitHub), Resource Server (API провайдера для профиля).

**Authorization Code flow:**

1. Пользователь нажимает «Войти через Google» → наш сервер редиректит на Google с `client_id`, `redirect_uri`, `scope`, `state` (защита от CSRF).
2. Пользователь логинится на стороне Google и разрешает доступ → Google редиректит на наш `redirect_uri` с `code` и `state`.
3. Наш callback handler проверяет `state`, обменивает `code` на access token (POST к token endpoint провайдера с client_secret).
4. Наш сервер запрашивает профиль пользователя (UserInfo endpoint провайдера с access token), получает email/id.
5. По email или внешнему id мы создаём или находим User в нашем хранилище, выдаём пользователю нашу сессию или JWT (cookie или токен в ответе/редиректе).

Дальнейшие запросы к нашему API защищаются уже нашей сессией/JWT, а не токеном провайдера.

---

## Точка входа: main.go

Планируется: конфигурация OAuth (client_id, client_secret для Google/GitHub), UserRepository, AuthUsecase (get-or-create по OAuth profile), OAuth handlers (redirect + callback), AuthHandler (Me, Delete), middleware (наша сессия/JWT). Роуты: `/health`, `GET /api/v1/auth/oauth/:provider` (redirect), `GET /api/v1/auth/oauth/callback` (callback), защищённые `GET /me`, `DELETE /me`. Сервер и Graceful Shutdown — по шаблону.

---

## Clean Architecture

Domain: User (возможно поле Provider, ExternalID для связи с аккаунтом провайдера). Use Case: GetOrCreateUserByOAuth(provider, profile) — поиск по external_id или email, при отсутствии — создание User; GetUserByID. Delivery: OAuth redirect/callback handlers (HTTP-специфичные: редиректы, query params), вызов use case с профилем; после логина установка нашей сессии/JWT и редирект на фронт или выдача токена.

---

## Domain Layer (Доменный слой)

**User:** ID, Email, возможно Provider (google/github), ExternalID (id у провайдера), CreatedAt. Пароль может быть пустым для OAuth-only пользователей.

**UserRepository:** Create, GetByID, GetByEmail, GetByExternalID(provider, externalID) — для поиска при callback.

---

## Repository Layer (Слой данных)

In-memory UserRepository с поддержкой GetByExternalID (перебор или отдельный индекс по provider+externalID).

---

## Use Case Layer (Бизнес-логика)

**GetOrCreateUserByOAuth(provider, externalID, email, name):** GetByExternalID(provider, externalID); если найден — вернуть user; иначе GetByEmail(email) или создать нового User с данными из профиля, сохранить, вернуть user. Пароль для OAuth-пользователя не задаётся или генерируется случайный.

**GetUserByID:** как в предыдущих проектах. После OAuth логина наш сервер выдаёт сессию или JWT с этим user_id; middleware проверяет нашу сессию/JWT.

---

## Delivery Layer (HTTP и Middleware)

**oauth_handler.go:** RedirectHandler(provider) — формирует URL провайдера (client_id, redirect_uri, scope, state), редирект 302. CallbackHandler(provider) — читает query code и state, проверяет state, обменивает code на access token (HTTP-клиент к провайдеру), запрашивает UserInfo, вызывает UseCase.GetOrCreateUserByOAuth, устанавливает нашу сессию (cookie) или возвращает наш JWT, редирект на клиентский URL или JSON с token.

**auth_handler.go:** Me, Delete — user_id из context (middleware нашей сессии/JWT).

**middleware/auth.go:** проверка нашей сессии или JWT (как в 03 или 04), user_id в context.

---

## Полный Flow запроса

Пользователь → GET /auth/oauth/google → редирект на Google → логин на Google → редирект на наш /auth/oauth/callback?code=...&state=... → обмен code на token → запрос UserInfo → GetOrCreateUserByOAuth → Set-Cookie (наша сессия) или выдача JWT → редирект на фронт. Дальше запросы к GET /me с нашей cookie/JWT → middleware → MeHandler.

---

## Зависимости между слоями

OAuth-клиент (обмен code, UserInfo) — в пакете или delivery; use case получает уже готовый профиль (provider, externalID, email, name). Use case не знает о HTTP и редиректах. Middleware для наших эндпоинтов — только наша сессия/JWT.

---

## Преимущества архитектуры

Чёткое разделение: OAuth flow и редиректы в delivery, создание/поиск пользователя в use case; тестируемость GetOrCreateUserByOAuth без HTTP; возможность добавить несколько провайдеров через один use case.

---

## Резюме

**Планируется:** OAuth 2.0 Authorization Code с Google/GitHub, redirect и callback, получение профиля, GetOrCreateUserByOAuth, выдача нашей сессии/JWT для доступа к API. После реализации — подставить актуальный код и диаграммы.

**Следующий проект:** [07-sso-oidc](../07-sso-oidc/) — SSO / OpenID Connect.
