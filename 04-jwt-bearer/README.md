# 04 - JWT Bearer Token

**Быстрый старт (после реализации):**

- **Go:** `cd 04-jwt-bearer/go && go run cmd/server/main.go`
- **Node.js:** `cd 04-jwt-bearer/node && npm install && npm run dev`
- **C# .NET:** `cd 04-jwt-bearer/csharp && dotnet run`

**Swagger:** `http://localhost:8080/swagger`

Stateless авторизация через JWT. Login возвращает токен, защищённые роуты проверяют `Authorization: Bearer <token>`.

> TODO: реализация

## Swagger / OpenAPI

После реализации в каждой версии (Go/Node/C#) будет доступно:
- Swagger UI: `http://localhost:8080/swagger`
- OpenAPI JSON: `http://localhost:8080/openapi.json` (или эквивалент для конкретного стека)
