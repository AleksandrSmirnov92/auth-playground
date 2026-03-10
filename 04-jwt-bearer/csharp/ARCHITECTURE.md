# Архитектура: 04 — JWT Bearer Token — C# .NET

## Обзор

Stateless авторизация через JWT: Login возвращает токен, клиент передаёт `Authorization: Bearer <token>`. Middleware верифицирует JWT и извлекает user_id из claims.

**Статус:** TODO.

### Планируемая структура

```
04-jwt-bearer/csharp/
├── Domain/
├── UseCase/              # Login (возврат JWT)
├── Delivery/             # Login (выдача JWT), Bearer middleware
└── pkg/Jwt/              # Подпись и верификация (System.IdentityModel.Tokens.Jwt)
```

**Следующий (C#):** [05-access-refresh-tokens](../../05-access-refresh-tokens/csharp/ARCHITECTURE.md)
