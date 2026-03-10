# Архитектура: 02 — API Key — C# .NET

## Обзор

Авторизация по заголовку **X-API-Key**: клиент передаёт ключ при каждом запросе. Ключ выдаётся при регистрации. Middleware проверяет ключ и передаёт user_id в HttpContext.Items.

**Статус:** TODO.

### Планируемая структура

```
02-api-key/csharp/
├── Program.cs
├── Domain/               # User + ApiKey, IUserRepository (GetByAPIKey)
├── Repository/
├── UseCase/              # Register (генерация ключа), GetUserByAPIKey
├── Delivery/
│   └── Middleware/       # Проверка X-API-Key
└── BasicAuth.csproj
```

**Следующий (C#):** [03-session-cookie](../../03-session-cookie/csharp/ARCHITECTURE.md)
