# Архитектура: 03 — Session + Cookie — C# .NET

## Обзор

Серверная сессионная авторизация: Login создаёт сессию, session_id в Cookie. Middleware по cookie находит пользователя. Logout — удаление сессии.

**Статус:** TODO.

### Планируемая структура

```
03-session-cookie/csharp/
├── Domain/               # User, Session, ISessionRepository
├── UseCase/              # Login (создание сессии), Logout
├── Delivery/             # Login (Set-Cookie), Logout, middleware по cookie
└── ...
```

**Следующий (C#):** [04-jwt-bearer](../../04-jwt-bearer/csharp/ARCHITECTURE.md)
