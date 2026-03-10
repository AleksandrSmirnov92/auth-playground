# Архитектура: 07 — SSO / OpenID Connect

## Содержание

1. [Обзор проекта](#обзор-проекта)
2. [SSO и OpenID Connect: как это работает](#sso-и-openid-connect-как-это-работает)
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

Этот мини-проект реализует **Single Sign-On (SSO)** на основе **OpenID Connect (OIDC)** — надстройки над OAuth 2.0 для аутентификации. Пользователь входит один раз через Identity Provider (IdP); IdP возвращает **ID Token** (JWT с claims о пользователе: sub, email и т.д.) и при необходимости access token для **UserInfo** endpoint. Наше приложение проверяет ID Token (подпись, issuer, audience), извлекает идентификатор пользователя и создаёт/находит пользователя в своей системе, после чего выдаёт свою сессию или JWT для доступа к нашим API. SSO — единый вход для нескольких приложений через одного провайдера.

**Статус:** планируемая реализация (плейсхолдеры).

### Планируемая структура файлов

```
07-sso-oidc/go/
├── cmd/server/main.go
├── internal/
│   ├── domain/
│   │   ├── user.go
│   │   └── repository.go
│   ├── repository/memory/
│   │   └── user_repository.go
│   ├── usecase/
│   │   └── auth_usecase.go      # GetOrCreateUserByOIDC(claims), GetUserByID
│   └── delivery/
│       ├── auth_handler.go
│       ├── oidc_handler.go      # Redirect, Callback (code → tokens, верификация ID Token)
│       └── middleware/
│           └── auth.go
├── pkg/ или internal/oidc/       # Discovery (JWKS, issuer), верификация ID Token
├── go.mod
└── go.sum
```

---

## SSO и OpenID Connect: как это работает

**OpenID Connect** добавляет к OAuth 2.0:

- **ID Token** — JWT, подписанный IdP, с claims: `sub` (subject, идентификатор пользователя), `email`, `iss` (issuer), `aud` (audience), `exp`, `iat`. Клиент или наш backend верифицирует подпись по ключам IdP (JWKS из discovery endpoint).
- **Discovery** — URL вида `/.well-known/openid-configuration` у IdP возвращает ссылки на authorization endpoint, token endpoint, JWKS URI, UserInfo и т.д.
- **UserInfo** — опционально; при необходимости запрашивается по access token для дополнительных полей.

**Flow (похож на OAuth 2.0):** редирект на IdP → пользователь логинится → редирект на наш callback с `code` → обмен code на **id_token** и access_token → верификация ID Token (подпись по JWKS, проверка iss, aud, exp) → извлечение `sub`/email из claims → GetOrCreateUser в нашем приложении → выдача нашей сессии/JWT. Дальнейшие запросы к нашему API — по нашей сессии/JWT. **SSO:** один вход в IdP даёт доступ ко всем приложениям, доверяющим этому IdP.

---

## Точка входа: main.go

Планируется: конфигурация OIDC (issuer, client_id, client_secret, discovery URL или явные endpoints), UserRepository, AuthUsecase (get-or-create по claims из ID Token), OIDC handlers (redirect + callback с верификацией ID Token), AuthHandler (Me, Delete), middleware нашей сессии/JWT. Роуты: `/health`, `GET /api/v1/auth/oidc/:provider` (redirect), `GET /api/v1/auth/oidc/callback`, защищённые `GET /me`, `DELETE /me`. Сервер и Graceful Shutdown — по шаблону.

---

## Clean Architecture

Domain: User (возможно Provider, ExternalID = sub из OIDC). Use Case: GetOrCreateUserByOIDC(provider, sub, email, ...) — по sub/email ищем или создаём пользователя; GetUserByID. Delivery: OIDC redirect/callback, верификация ID Token (в пакете oidc или в handler), извлечение claims, вызов use case, установка нашей сессии/JWT.

---

## Domain Layer (Доменный слой)

**User:** ID, Email, Provider (oidc provider name), ExternalID (sub из IdP), CreatedAt. Пароль для SSO-only пользователей может быть пустым.

**UserRepository:** Create, GetByID, GetByEmail, GetByExternalID(provider, externalID).

---

## Repository Layer (Слой данных)

In-memory UserRepository с GetByExternalID, как в 06-oauth2.

---

## Use Case Layer (Бизнес-логика)

**GetOrCreateUserByOIDC(provider, sub, email, name):** GetByExternalID(provider, sub); если найден — вернуть user; иначе GetByEmail или создание User с данными из claims, сохранение, возврат user.

**GetUserByID:** как в предыдущих проектах. Middleware проверяет нашу сессию/JWT и передаёт user_id в context.

---

## Delivery Layer (HTTP и Middleware)

**oidc_handler.go:** RedirectHandler — редирект на authorization endpoint IdP (OAuth 2.0 + scope openid). CallbackHandler — приём code и state, обмен code на id_token и access_token; вызов пакета/функции верификации ID Token (JWKS, iss, aud, exp); извлечение claims (sub, email); вызов UseCase.GetOrCreateUserByOIDC; установка нашей сессии или выдача JWT; редирект или JSON.

**auth_handler.go:** Me, Delete — user_id из context.

**middleware/auth.go:** проверка нашей сессии или JWT после OIDC login.

**Пакет oidc (или internal/oidc):** загрузка JWKS по discovery или по URL, верификация подписи ID Token, парсинг claims. Зависит только от конфигурации IdP и криптографии, не от домена.

---

## Полный Flow запроса

Пользователь → GET /auth/oidc/keycloak (или другого IdP) → редирект на IdP → логин → редирект на наш callback с code → обмен code на id_token (и access_token) → верификация ID Token по JWKS → извлечение sub, email → GetOrCreateUserByOIDC → наша сессия/JWT → редирект на фронт. Дальше запросы к GET /me с нашей cookie/JWT → middleware → MeHandler. Один вход в IdP — доступ ко всем приложениям, настроенным на этот IdP (SSO).

---

## Зависимости между слоями

OIDC-клиент и верификация ID Token — в отдельном пакете или delivery; use case получает уже готовые идентификаторы (provider, sub, email, name). Use case не знает о JWT и discovery. Middleware для наших эндпоинтов — только наша сессия/JWT.

---

## Преимущества архитектуры

Разделение: OIDC flow и криптографическая верификация в одном месте, бизнес-логика «кто такой пользователь у нас» в use case; возможность подключать несколько IdP (Keycloak, Auth0, Google как OIDC); тестируемость GetOrCreateUserByOIDC без HTTP.

---

## Резюме

**Планируется:** SSO на базе OpenID Connect, discovery и верификация ID Token (JWKS), извлечение claims (sub, email), GetOrCreateUserByOIDC, выдача нашей сессии/JWT для доступа к API. После реализации — подставить актуальный код и диаграммы.

**Коллекция реализаций:** все виды авторизации описаны в [корневом README](../../../README.md).
