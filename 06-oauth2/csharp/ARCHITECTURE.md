# Архитектура: 06 — OAuth 2.0 — C# .NET

## Обзор

OAuth 2.0 Authorization Code flow с Google/GitHub. Redirect на провайдера → callback с code → обмен на token → профиль пользователя → создание/поиск пользователя → наша сессия/JWT.

**Статус:** TODO.

### Планируемая структура

```
06-oauth2/csharp/
├── Domain/
├── UseCase/              # GetOrCreateUserByOAuth
├── Delivery/             # Redirect, Callback, OAuth clients
└── pkg/OAuth/            # Google, GitHub клиенты
```

**Следующий (C#):** [07-sso-oidc](../../07-sso-oidc/csharp/ARCHITECTURE.md)
