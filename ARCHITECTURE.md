# Архитектура проекта auth-playground

## Структура репозитория

Каждый вид авторизации реализован на **трёх языках**: Go, Node.js/TypeScript, C# .NET. Подробная архитектура каждого проекта описана в своей папке.

```
auth-playground/
├── 01-basic-auth/
│   ├── go/          → ARCHITECTURE.md (Go)
│   ├── node/        → ARCHITECTURE.md (Node.js/TypeScript)
│   └── csharp/      → ARCHITECTURE.md (C# .NET)
├── 02-api-key/
│   ├── go/
│   ├── node/
│   └── csharp/
├── ... (03–07)
└── ARCHITECTURE.md  ← этот файл
```

## Общие принципы

Все реализации используют **Clean Architecture**:
- **Domain** — модели и интерфейсы
- **Use Case** — бизнес-логика (регистрация, логин и т.д.)
- **Repository** — хранение данных (in-memory)
- **Delivery** — HTTP handlers, middleware

## Где искать детали

| Вид авторизации | Go | Node | C# |
|-----------------|-----|------|-----|
| 01 Basic Auth | [go/ARCHITECTURE.md](01-basic-auth/go/ARCHITECTURE.md) | [node/ARCHITECTURE.md](01-basic-auth/node/ARCHITECTURE.md) | [csharp/ARCHITECTURE.md](01-basic-auth/csharp/ARCHITECTURE.md) |
| 02 API Key | [go/ARCHITECTURE.md](02-api-key/go/ARCHITECTURE.md) | [node/ARCHITECTURE.md](02-api-key/node/ARCHITECTURE.md) | [csharp/ARCHITECTURE.md](02-api-key/csharp/ARCHITECTURE.md) |
| 03 Session+Cookie | [go/ARCHITECTURE.md](03-session-cookie/go/ARCHITECTURE.md) | [node/ARCHITECTURE.md](03-session-cookie/node/ARCHITECTURE.md) | [csharp/ARCHITECTURE.md](03-session-cookie/csharp/ARCHITECTURE.md) |
| 04 JWT Bearer | [go/ARCHITECTURE.md](04-jwt-bearer/go/ARCHITECTURE.md) | [node/ARCHITECTURE.md](04-jwt-bearer/node/ARCHITECTURE.md) | [csharp/ARCHITECTURE.md](04-jwt-bearer/csharp/ARCHITECTURE.md) |
| 05 Access+Refresh | [go/ARCHITECTURE.md](05-access-refresh-tokens/go/ARCHITECTURE.md) | [node/ARCHITECTURE.md](05-access-refresh-tokens/node/ARCHITECTURE.md) | [csharp/ARCHITECTURE.md](05-access-refresh-tokens/csharp/ARCHITECTURE.md) |
| 06 OAuth 2.0 | [go/ARCHITECTURE.md](06-oauth2/go/ARCHITECTURE.md) | [node/ARCHITECTURE.md](06-oauth2/node/ARCHITECTURE.md) | [csharp/ARCHITECTURE.md](06-oauth2/csharp/ARCHITECTURE.md) |
| 07 SSO/OIDC | [go/ARCHITECTURE.md](07-sso-oidc/go/ARCHITECTURE.md) | [node/ARCHITECTURE.md](07-sso-oidc/node/ARCHITECTURE.md) | [csharp/ARCHITECTURE.md](07-sso-oidc/csharp/ARCHITECTURE.md) |
