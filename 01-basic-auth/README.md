# 01 - Basic Auth (Base64)

HTTP Basic Authentication — простейший стандарт авторизации (RFC 7617).

Клиент передаёт `email:password` в заголовке `Authorization`, закодированные в Base64.

## Как это работает

```
Client                              Server
  |                                    |
  |  POST /register {email, pass}      |
  |----------------------------------->|  (пароль хешируется bcrypt)
  |  201 Created                       |
  |<-----------------------------------|
  |                                    |
  |  GET /me                           |
  |  Authorization: Basic base64(...)  |
  |----------------------------------->|  middleware декодирует Base64
  |                                    |  → email:password
  |                                    |  → проверяет через bcrypt
  |  200 OK {user}                     |
  |<-----------------------------------|
```

## Endpoints

| Метод | URL | Auth | Описание |
|-------|-----|------|----------|
| POST | `/api/v1/auth/register` | Нет | Регистрация нового пользователя |
| GET | `/api/v1/auth/me` | Basic | Получить профиль текущего пользователя |
| DELETE | `/api/v1/auth/me` | Basic | Удалить текущего пользователя |
| GET | `/health` | Нет | Проверка работоспособности сервера |

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
npm run build && npm start
# или для разработки: npm run dev
```

**C# .NET:**
```bash
cd 01-basic-auth/csharp
dotnet run
```

## Примеры curl

### Регистрация

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "secret123"}'
```

### Получить профиль (Basic Auth)

```bash
curl -u user@example.com:secret123 http://localhost:8080/api/v1/auth/me
```

Флаг `-u` автоматически кодирует `email:password` в Base64 и добавляет заголовок:
```
Authorization: Basic dXNlckBleGFtcGxlLmNvbTpzZWNyZXQxMjM=
```

### Удалить пользователя

```bash
curl -X DELETE -u user@example.com:secret123 http://localhost:8080/api/v1/auth/me
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
- Basic Auth передаёт credentials в каждом запросе — используйте только через HTTPS
- Заголовок `WWW-Authenticate` возвращается при неудачной авторизации
