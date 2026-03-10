# Архитектура: 01 — Basic Auth (Base64) — C# .NET

## Содержание

1. [Обзор проекта](#обзор-проекта)
2. [Как работает авторизация в этом проекте](#как-работает-авторизация-в-этом-проекте)
3. [Точка входа: Program.cs](#точка-входа-programcs)
4. [Clean Architecture](#clean-architecture)
5. [Domain Layer (Доменный слой)](#domain-layer-доменный-слой)
6. [Repository Layer (Слой данных)](#repository-layer-слой-данных)
7. [Use Case Layer (Бизнес-логика)](#use-case-layer-бизнес-логика)
8. [Delivery Layer (HTTP endpoints)](#delivery-layer-http-endpoints)
9. [Contracts / DTO](#contracts--dto)
10. [Swagger / OpenAPI](#swagger--openapi)
11. [Полный цикл запроса](#полный-цикл-запроса)
12. [Зависимости между слоями](#зависимости-между-слоями)
13. [Резюме](#резюме)

---

## Обзор проекта

Реализация аутентификации по email + паролю на **C# / ASP.NET Core Minimal APIs**.

Пароли хешируются через **BCrypt.Net-Next**. Клиент передаёт email и пароль в теле каждого запроса (JSON). Никаких сессий и токенов — это базовый пример "stateless" аутентификации.

### Структура файлов (C# .NET)

```
01-basic-auth/csharp/
├── Program.cs                              # Точка входа, DI, запуск сервера
├── Domain/
│   ├── User.cs                             # Модель User (доменная сущность)
│   └── IUserRepository.cs                  # Интерфейс репозитория
├── Repository/
│   └── MemoryUserRepository.cs             # In-memory реализация репозитория
├── UseCase/
│   └── AuthUsecase.cs                      # Бизнес-логика: register / login / delete
├── Delivery/
│   ├── AuthEndpoints.cs                    # Minimal API endpoints + RegisterRequest DTO
│   ├── Contracts/
│   │   ├── UserResponse.cs                 # Публичный ответ: id, email, created_at
│   │   └── AuthResponse.cs                 # Ответ с сообщением: { message, user }
│   └── Swagger/
│       └── ExamplesSchemaFilter.cs         # Реальные примеры для Swagger UI
├── BasicAuth.csproj
├── appsettings.json
└── ARCHITECTURE.md
```

---

## Как работает авторизация в этом проекте

Клиент отправляет `email` и `password` в **теле запроса** (JSON):

- `POST /api/v1/auth/register` — создать аккаунт
- `POST /api/v1/auth/login` — войти (проверить пароль)
- `DELETE /api/v1/auth/delete` — удалить аккаунт (нужен email + пароль)

> **Почему не в заголовке?** В этом проекте мы изучаем хранение и проверку паролей (BCrypt).
> Передача в теле JSON — самый прозрачный способ показать эту логику.

---

## Точка входа: Program.cs

```csharp
builder.Services.AddSingleton<IUserRepository, MemoryUserRepository>();
builder.Services.AddSingleton<AuthUsecase>();
builder.Services.AddSwaggerGen(c => { c.SchemaFilter<ExamplesSchemaFilter>(); });

var app = builder.Build();
app.UseSwagger();
app.UseSwaggerUI(c => { c.RoutePrefix = "swagger"; });

app.MapAuthEndpoints(authUsecase);
app.Run("http://localhost:8080");
```

**Dependency Injection** (DI) — встроенный механизм ASP.NET Core для связывания зависимостей.
- `AddSingleton` — один экземпляр на всё время работы приложения.
- `GetRequiredService<T>()` — получить экземпляр из контейнера.

BasicAuth middleware удалён: каждый endpoint теперь сам принимает email+пароль в теле.

---

## Clean Architecture

```
Delivery  →  UseCase  →  Domain
                 ↑
            Repository
```

- **Domain** — модели и интерфейсы. Не знает про HTTP и базы данных.
- **Use Case** — бизнес-правила. Не знает про ASP.NET Core.
- **Repository** — реализует `IUserRepository` (сейчас in-memory, можно заменить на SQL).
- **Delivery** — Minimal API endpoints, маппинг JSON ↔ Use Case, Swagger.

---

## Domain Layer (Доменный слой)

### Domain/User.cs

```csharp
public class User
{
    public required string Id { get; set; }      // UUID — уникальный идентификатор
    public required string Email { get; set; }   // логин
    public required string Password { get; set; } // bcrypt-хеш (не plain text!)
    public DateTime CreatedAt { get; set; }
}
```

`required` — ключевое слово C# 11: компилятор требует заполнить поле при создании объекта.

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

`Task<T>` — асинхронное значение (аналог Promise в JS). Все методы возвращают Task, чтобы в будущем можно было подключить настоящую БД (async I/O).

---

## Repository Layer (Слой данных)

**MemoryUserRepository.cs:**
- Хранение: `Dictionary<string, User>` (ключ — id).
- Все методы — асинхронные, но внутри пока синхронные (готово к замене на EF Core).
- `GetByEmailAsync` — линейный перебор `O(n)`; в реальной БД это индекс по email.

---

## Use Case Layer (Бизнес-логика)

**AuthUsecase.cs:**

| Метод | Что делает |
|-------|-----------|
| `RegisterAsync` | Проверяет, нет ли пользователя с таким email; хеширует пароль BCrypt; сохраняет User с Guid |
| `LoginAsync` | Находит пользователя по email; вызывает `BCrypt.Verify` — при ошибке `UnauthorizedAccessException` |
| `DeleteByCredentialsAsync` | Вызывает `LoginAsync` для проверки пароля, затем `DeleteUserByIdAsync` |
| `GetUserByIdAsync` / `DeleteUserByIdAsync` | Передают вызов в repository |

**Защита от user enumeration:** `LoginAsync` при неверном email и при неверном пароле бросает одно и то же исключение (`invalid email or password`) — злоумышленник не может проверить, существует ли email.

---

## Delivery Layer (HTTP endpoints)

**AuthEndpoints.cs** — Minimal API, extension-метод `MapAuthEndpoints`:

| Метод | URL | Тело запроса | Ответ |
|-------|-----|-------------|-------|
| GET | `/health` | — | `{ status: "ok" }` |
| POST | `/api/v1/auth/register` | `{ email, password }` | 201 `AuthResponse` |
| POST | `/api/v1/auth/login` | `{ email, password }` | 200 `AuthResponse` |
| DELETE | `/api/v1/auth/delete` | `{ email, password }` | 200 `{ message }` |

```csharp
app.MapDelete("/api/v1/auth/delete", async (RegisterRequest req) =>
{
    await authUsecase.DeleteByCredentialsAsync(req.Email, req.Password);
    return Results.Json(new { message = "Пользователь успешно удалён" });
});
```

---

## Contracts / DTO

> DTO (Data Transfer Object) — объект для передачи данных между слоями. Здесь: между UseCase и HTTP-ответом.

**UserResponse.cs** — то, что клиент получает в JSON (без пароля):

```csharp
public record UserResponse(
    [property: JsonPropertyName("id")]         string Id,
    [property: JsonPropertyName("email")]      string Email,
    [property: JsonPropertyName("created_at")] DateTime CreatedAt
);
```

`JsonPropertyName` — атрибут .NET для управления именем поля в JSON. Здесь: `created_at` вместо `CreatedAt` (snake_case как в Go/Node).

**AuthResponse.cs** — ответ с сообщением и пользователем:

```csharp
public record AuthResponse(
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("user")]    UserResponse User
);
```

---

## Swagger / OpenAPI

- Swagger UI: `http://localhost:8080/swagger`
- OpenAPI JSON (Swashbuckle): `http://localhost:8080/swagger/v1/swagger.json`

`ExamplesSchemaFilter` — кастомный `ISchemaFilter` Swashbuckle. Вместо `"string"` подставляет реальные значения (`ivanov@example.com`, `1234`) в описание схемы в UI.

---

## Полный цикл запроса

Пример: `DELETE /api/v1/auth/delete`

```
HTTP DELETE /api/v1/auth/delete
Body: { "email": "ivanov@example.com", "password": "1234" }
   ↓
AuthEndpoints.MapDelete
   — парсит RegisterRequest из тела
   — вызывает authUsecase.DeleteByCredentialsAsync(email, password)
   ↓
AuthUsecase.DeleteByCredentialsAsync
   — вызывает LoginAsync → GetByEmailAsync → BCrypt.Verify
   — вызывает DeleteUserByIdAsync → repository.DeleteAsync
   ↓
MemoryUserRepository.DeleteAsync
   — удаляет из Dictionary
   ↓
200 OK: { "message": "Пользователь успешно удалён" }
```

---

## Зависимости между слоями

```
HTTP-клиент
  → AuthEndpoints (Delivery)
    → AuthUsecase (UseCase)
      → IUserRepository (Domain/Repository)
        → MemoryUserRepository (Repository)
```

UseCase не знает про HTTP (нет `HttpContext`).
Domain не знает про BCrypt — хеширование в UseCase.

---

## Резюме

**Реализовано:** аутентификация по email+паролю, BCrypt.Net-Next, Minimal API, Clean Architecture, Swagger.

**Следующий проект (C#):** [02-api-key](../../02-api-key/csharp/ARCHITECTURE.md) — авторизация по X-API-Key.
