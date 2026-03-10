# Архитектура: 05 — Access + Refresh Tokens — Node.js / TypeScript

## Обзор

Access Token (JWT, короткоживущий) + Refresh Token (opaque, на сервере). При истечении access — POST /refresh с refresh token для получения новой пары.

**Статус:** TODO.

### Планируемая структура

```
05-access-refresh-tokens/node/
├── src/
│   ├── domain/           # RefreshToken, RefreshTokenRepository
│   ├── usecase/          # Login (пара токенов), Refresh
│   └── delivery/         # Login, Refresh, middleware только по Access JWT
└── ...
```

**Следующий (Node):** [06-oauth2](../../06-oauth2/node/ARCHITECTURE.md)
