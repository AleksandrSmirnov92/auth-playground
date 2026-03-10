# Архитектура: 07 — SSO / OpenID Connect — C# .NET

## Обзор

SSO на базе OIDC: IdP возвращает ID Token (JWT) с claims. Верификация подписи (JWKS), извлечение sub/email, GetOrCreateUserByOIDC, выдача нашей сессии/JWT.

**Статус:** TODO.

### Планируемая структура

```
07-sso-oidc/csharp/
├── Domain/
├── UseCase/              # GetOrCreateUserByOIDC
├── Delivery/             # Redirect, Callback, ID Token верификация
└── pkg/Oidc/             # Discovery, JWKS, верификация ID Token
```

[Корневой README](../../../README.md)
