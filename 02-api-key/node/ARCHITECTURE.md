# Архитектура: 02 — API Key — Node.js / TypeScript

## Обзор

Авторизация по заголовку **X-API-Key**: клиент передаёт ключ при каждом запросе. Ключ выдаётся при регистрации. Middleware проверяет ключ и передаёт user_id в request.

**Статус:** TODO.

### Планируемая структура

```
02-api-key/node/
├── src/
│   ├── index.ts
│   ├── domain/           # User + ApiKey, UserRepository (GetByAPIKey)
│   ├── repository/
│   ├── usecase/          # Register (генерация ключа), GetUserByAPIKey
│   └── delivery/
│       ├── authHandler.ts
│       └── middleware/   # Проверка X-API-Key
├── package.json
└── tsconfig.json
```

**Следующий (Node):** [03-session-cookie](../../03-session-cookie/node/ARCHITECTURE.md)
