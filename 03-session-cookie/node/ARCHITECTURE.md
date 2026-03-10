# Архитектура: 03 — Session + Cookie — Node.js / TypeScript

## Обзор

Серверная сессионная авторизация: Login создаёт сессию, session_id в Cookie. Middleware по cookie находит пользователя. Logout — удаление сессии.

**Статус:** TODO.

### Планируемая структура

```
03-session-cookie/node/
├── src/
│   ├── domain/           # User, Session, SessionRepository
│   ├── usecase/          # Login (создание сессии), Logout
│   └── delivery/         # Login (Set-Cookie), Logout, middleware по cookie
└── ...
```

**Следующий (Node):** [04-jwt-bearer](../../04-jwt-bearer/node/ARCHITECTURE.md)
