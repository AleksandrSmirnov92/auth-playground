# Архитектура: 04 — JWT Bearer Token — Node.js / TypeScript

## Обзор

Stateless авторизация через JWT: Login возвращает токен, клиент передаёт `Authorization: Bearer <token>`. Middleware верифицирует JWT и извлекает user_id из claims.

**Статус:** TODO.

### Планируемая структура

```
04-jwt-bearer/node/
├── src/
│   ├── domain/
│   ├── usecase/          # Login (возврат JWT)
│   ├── delivery/         # Login (выдача JWT), middleware Bearer
│   └── pkg/jwt/          # Подпись и верификация (jsonwebtoken)
└── ...
```

**Следующий (Node):** [05-access-refresh-tokens](../../05-access-refresh-tokens/node/ARCHITECTURE.md)
