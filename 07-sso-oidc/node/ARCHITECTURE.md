# Архитектура: 07 — SSO / OpenID Connect — Node.js / TypeScript

## Обзор

SSO на базе OIDC: IdP возвращает ID Token (JWT) с claims. Верификация подписи (JWKS), извлечение sub/email, GetOrCreateUserByOIDC, выдача нашей сессии/JWT.

**Статус:** TODO.

### Планируемая структура

```
07-sso-oidc/node/
├── src/
│   ├── domain/
│   ├── usecase/          # GetOrCreateUserByOIDC
│   ├── delivery/         # Redirect, Callback, ID Token верификация
│   └── pkg/oidc/         # Discovery, JWKS, верификация ID Token
└── ...
```

[Корневой README](../../../README.md)
