# Архитектура: 06 — OAuth 2.0 — Node.js / TypeScript

## Обзор

OAuth 2.0 Authorization Code flow с Google/GitHub. Redirect на провайдера → callback с code → обмен на token → профиль пользователя → создание/поиск пользователя → наша сессия/JWT.

**Статус:** TODO.

### Планируемая структура

```
06-oauth2/node/
├── src/
│   ├── domain/
│   ├── usecase/          # GetOrCreateUserByOAuth
│   ├── delivery/         # Redirect, Callback handler, oauth clients
│   └── pkg/oauth/        # Google, GitHub клиенты
└── ...
```

**Следующий (Node):** [07-sso-oidc](../../07-sso-oidc/node/ARCHITECTURE.md)
