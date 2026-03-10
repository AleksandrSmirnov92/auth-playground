# 05 - Access + Refresh Tokens

**Быстрый старт (после реализации):**

- **Go:** `cd 05-access-refresh-tokens/go && go run cmd/server/main.go`
- **Node.js:** `cd 05-access-refresh-tokens/node && npm install && npm run dev`
- **C# .NET:** `cd 05-access-refresh-tokens/csharp && dotnet run`

**Swagger:** `http://localhost:8080/swagger`

Короткоживущий Access Token + долгоживущий Refresh Token. Endpoint `/refresh` для обновления.

> TODO: реализация

## Swagger / OpenAPI

После реализации в каждой версии (Go/Node/C#) будет доступно:
- Swagger UI: `http://localhost:8080/swagger`
- OpenAPI JSON: `http://localhost:8080/openapi.json` (или эквивалент для конкретного стека)
