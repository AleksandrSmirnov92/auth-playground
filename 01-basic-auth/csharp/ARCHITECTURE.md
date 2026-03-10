# Архитектура: 01 — Basic Auth (Base64) — C# .NET

## Содержание

1. [Обзор проекта](#обзор-проекта)
2. [Basic Auth: как это работает](#basic-auth-как-это-работает)
3. [Точка входа: Program.cs](#точка-входа-programcs)
4. [Clean Architecture](#clean-architecture)
5. [Domain Layer (Доменный слой)](#domain-layer-доменный-слой)
6. [Repository Layer (Слой данных)](#repository-layer-слой-данных)
7. [Use Case Layer (Бизнес-логика)](#use-case-layer-бизнес-логика)
8. [Delivery Layer (HTTP и Middleware)](#delivery-layer-http-и-middleware)
9. [Зависимости между слоями](#зависимости-между-слоями)
10. [Резюме](#резюме)

---

## Обзор проекта

Реализация **HTTP Basic Authentication** (RFC 7617) на **C# / ASP.NET Core**. Клиент передаёт логин и пароль в заголовке `Authorization` в каждом запросе. Регистрация — по JSON (email + password); доступ к `/me` и удаление аккаунта — только с Basic Auth. Пароли хешируются через BCrypt.Net-Next.

### Структура файлов (C# .NET)

```
01-basic-auth/csharp/
├── Program.cs                              # Точка входа, DI, middleware
├── Domain/
│   ├── User.cs                             # Модель User
│   └── IUserRepository.cs                  # Интерфейс репозитория
├── Repository/
│   └── MemoryUserRepository.cs             # In-memory реализация
├── UseCase/
│   └── AuthUsecase.cs                      # Регистрация, логин, bcrypt
├── Delivery/
│   ├── AuthEndpoints.cs                    # Minimal API endpoints
│   └── Middleware/
│       └── BasicAuthMiddleware.cs          # Basic Auth middleware
├── BasicAuth.csproj
├── appsettings.json
└── ARCHITECTURE.md
```

---

## Basic Auth: как это работает

**RFC 7617:** клиент кодирует `email:password` в Base64 и отправляет:

```
Authorization: Basic <base64(email:password)>
```

Отдельного `POST /login` нет — авторизация при каждом запросе к защищённым путям. Middleware декодирует заголовок, вызывает AuthUsecase.LoginAsync, при успехе кладёт UserId в HttpContext.Items и передаёт дальше.

---

## Точка входа: Program.cs

```csharp
builder.Services.AddSingleton<IUserRepository, MemoryUserRepository>();
builder.Services.AddSingleton<AuthUsecase>();

var app = builder.Build();
var authUsecase = app.Services.GetRequiredService<AuthUsecase>();
app.UseMiddleware<BasicAuthMiddleware>(authUsecase);

app.MapAuthEndpoints(authUsecase);
app.Run("http://localhost:8080");
```

Dependency Injection через ServiceCollection; middleware применяется только к защищённым путям (`/api/v1/auth/me`).

---

## Clean Architecture

Те же слои: Domain (модели и интерфейсы), Use Case (бизнес-логика), Repository (данные), Delivery (Minimal API + middleware). Domain не зависит от ASP.NET Core.

---

## Domain Layer (Доменный слой)

### Domain/User.cs

```csharp
public class User
{
    public required string Id { get; set; }
    public required string Email { get; set; }
    public required string Password { get; set; }
    public DateTime CreatedAt { get; set; }
}
```

### Domain/IUserRepository.cs

```csharp
public interface IUserRepository
{
    Task CreateAsync(User user);
    Task<User?> GetByIdAsync(string id);
    Task<User?> GetByEmailAsync(string email);
    Task DeleteAsync(string id);
}
```

---

## Repository Layer (Слой данных)

**MemoryUserRepository.cs:** `Dictionary<string, User>` для хранения. Методы CreateAsync, GetByIdAsync, GetByEmailAsync, DeleteAsync. Singleton в DI.

---

## Use Case Layer (Бизнес-логика)

**AuthUsecase.cs:**

- **RegisterAsync:** проверка по email, BCrypt.HashPassword, создание User с Guid, сохранение
- **LoginAsync:** GetByEmailAsync, BCrypt.Verify — при ошибке UnauthorizedAccessException
- **GetUserByIdAsync, DeleteUserByIdAsync:** вызов repository

---

## Delivery Layer (HTTP и Middleware)

**AuthEndpoints.cs:** MapGet/MapPost/MapDelete для /health, /api/v1/auth/register, /api/v1/auth/me. UserId берётся из HttpContext.Items (установлен middleware).

**BasicAuthMiddleware.cs:** проверяет путь (только /api/v1/auth/me); читает Authorization, декодирует Base64, вызывает LoginAsync; при успехе записывает user.Id в context.Items["UserId"].

---

## Зависимости между слоями

HTTP → Middleware → Endpoints → UseCase → Domain; Repository реализует IUserRepository. Use Case не знает про HTTP и Basic Auth — только проверяет email+password.

---

## Резюме

**Реализовано:** Basic Auth (RFC 7617), BCrypt.Net-Next, ASP.NET Core middleware, Minimal API, Clean Architecture.

**Следующий проект (C#):** [02-api-key](../../02-api-key/csharp/ARCHITECTURE.md) — авторизация по X-API-Key.
