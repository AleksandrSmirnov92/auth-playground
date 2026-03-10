# auth-playground — Все виды авторизации

Коллекция независимых мини-проектов, каждый из которых реализует один вид авторизации с нуля. **Каждый вид авторизации доступен на трёх языках: Go, Node.js/TypeScript и C# .NET.**

## Языки

| Язык | Папка | Стек |
|------|-------|------|
| Go | `go/` | stdlib, bcrypt |
| Node.js | `node/` | Express, TypeScript |
| C# | `csharp/` | ASP.NET Core |

## Реализации

| # | Вид авторизации | Go | Node | C# | Сложность |
|---|-----------------|-----|------|-----|-----------|
| [01](01-basic-auth/) | Basic Auth (Base64) | Done | Done | Done | Простая |
| [02](02-api-key/) | API Key | TODO | TODO | TODO | Простая |
| [03](03-session-cookie/) | Session + Cookie | TODO | TODO | TODO | Средняя |
| [04](04-jwt-bearer/) | JWT Bearer Token | TODO | TODO | TODO | Средняя |
| [05](05-access-refresh-tokens/) | Access + Refresh Tokens | TODO | TODO | TODO | Средняя |
| [06](06-oauth2/) | OAuth 2.0 (Google/GitHub) | TODO | TODO | TODO | Сложная |
| [07](07-sso-oidc/) | SSO / OpenID Connect | TODO | TODO | TODO | Сложная |

## Как запустить

```bash
cd 01-basic-auth/go        # Go
go run cmd/server/main.go

cd 01-basic-auth/node      # Node.js
npm install && npm run dev

cd 01-basic-auth/csharp    # C# .NET
dotnet run
```

Каждый сервер запускается на `http://localhost:8080` и имеет endpoint `/health`.

## Архитектура

Все реализации используют Clean Architecture: domain, usecase, delivery, repository.

## Порядок изучения

Проекты расположены от простого к сложному. Рекомендуется изучать по порядку.
