# 01 - Basic Auth

**Быстрый старт:**

- **Go:** `cd 01-basic-auth/go && go run cmd/server/main.go`
- **Node.js:** `cd 01-basic-auth/node && npm install && npm run dev`
- **C# .NET:** `cd 01-basic-auth/csharp && dotnet run`

**Swagger:** открой `http://localhost:8080/swagger`

Простейший пример аутентификации по email + паролю. Клиент передаёт email и пароль в теле каждого запроса (JSON). Пароли хешируются через bcrypt.

## Endpoints

| Метод | URL | Тело запроса | Описание |
|-------|-----|-------------|----------|
| POST | `/api/v1/auth/register` | `{ email, password }` | Регистрация нового пользователя |
| POST | `/api/v1/auth/login` | `{ email, password }` | Вход — проверка email + пароля |
| DELETE | `/api/v1/auth/delete` | `{ email, password }` | Удалить аккаунт |
| GET | `/health` | — | Проверка работоспособности сервера |
| GET | `/swagger` | — | Swagger UI (документация и тестирование) |
| GET | `/openapi.json` | — | OpenAPI спецификация (Go / Node) |

## Как это работает

```
Client                              Server
  |                                    |
  |  POST /register {email, pass}      |
  |----------------------------------->|  (пароль хешируется bcrypt)
  |  201 { message, user }             |
  |<-----------------------------------|
  |                                    |
  |  POST /login {email, pass}         |
  |----------------------------------->|  bcrypt.Compare(pass, hash)
  |  200 { message, user }             |
  |<-----------------------------------|
  |                                    |
  |  DELETE /delete {email, pass}      |
  |----------------------------------->|  логин → удаление
  |  200 { message }                   |
  |<-----------------------------------|
```

## Запуск

Каждая реализация запускается на `http://localhost:8080`.

**Go:**
```bash
cd 01-basic-auth/go
go run cmd/server/main.go
```

**Node.js / TypeScript:**
```bash
cd 01-basic-auth/node
npm install
npm run dev
```

**C# .NET:**
```bash
cd 01-basic-auth/csharp
dotnet run
```

## Swagger / OpenAPI

После запуска сервера открой в браузере:

- **Swagger UI**: `http://localhost:8080/swagger`
- **OpenAPI JSON**: `http://localhost:8080/openapi.json` (Go / Node)  
  или `http://localhost:8080/swagger/v1/swagger.json` (C#)

### Как протестировать через Swagger (шаги)

1. **Открыть Swagger UI**
   - Перейди в браузере на `http://localhost:8080/swagger`.

2. **Зарегистрировать пользователя**
   - Найди `POST /api/v1/auth/register`.
   - Нажми **Try it out**.
   - В теле уже будут примеры:
     ```json
     {
       "email": "ivanov@example.com",
       "password": "1234"
     }
     ```
   - Нажми **Execute** и убедись, что вернулся `201 Created`.

3. **Войти**
   - Найди `POST /api/v1/auth/login`.
   - Нажми **Try it out → Execute** с теми же данными.
   - Ожидаемый ответ — `200 OK` и JSON с пользователем.

4. **Удалить пользователя**
   - Найди `DELETE /api/v1/auth/delete`.
   - Нажми **Try it out → Execute** (данные уже подставлены).
   - Ожидаемый ответ — `200 OK` с `{ "message": "Пользователь успешно удалён" }`.

## Примеры curl

### Регистрация

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "ivanov@example.com", "password": "1234"}'
```

### Вход

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "ivanov@example.com", "password": "1234"}'
```

### Удалить аккаунт

```bash
curl -X DELETE http://localhost:8080/api/v1/auth/delete \
  -H "Content-Type: application/json" \
  -d '{"email": "ivanov@example.com", "password": "1234"}'
```

## Структура проекта

```
01-basic-auth/
├── go/                     # Реализация на Go
│   ├── ARCHITECTURE.md     # Подробная архитектура (Go)
│   ├── cmd/server/main.go
│   ├── internal/
│   └── go.mod
├── node/                   # Реализация на Node.js / TypeScript
│   ├── ARCHITECTURE.md     # Подробная архитектура (Node)
│   ├── src/
│   └── package.json
├── csharp/                 # Реализация на C# .NET
│   ├── ARCHITECTURE.md     # Подробная архитектура (C#)
│   ├── Domain/
│   ├── UseCase/
│   └── Delivery/
└── README.md
```

## Безопасность

- Пароли хешируются через **bcrypt** (не хранятся в открытом виде)
- Одинаковое сообщение `invalid email or password` при неверном email и неверном пароле — **защита от user enumeration**
