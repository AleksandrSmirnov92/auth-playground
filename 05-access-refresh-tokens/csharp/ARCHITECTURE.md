# Архитектура: 05 — Access + Refresh Tokens — C# .NET

## Обзор

Access Token (JWT, короткоживущий) + Refresh Token (opaque, на сервере). При истечении access — POST /refresh с refresh token для получения новой пары.

**Статус:** TODO.

### Планируемая структура

```
05-access-refresh-tokens/csharp/
├── Domain/               # RefreshToken, IRefreshTokenRepository
├── UseCase/              # Login (пара токенов), Refresh
├── Delivery/             # Login, Refresh, middleware только по Access JWT
└── ...
```

**Следующий (C#):** [06-oauth2](../../06-oauth2/csharp/ARCHITECTURE.md)
