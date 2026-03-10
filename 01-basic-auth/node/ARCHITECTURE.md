# Архитектура: 01 — Basic Auth (Base64) — Node.js / TypeScript

## Содержание

1. [Обзор проекта](#обзор-проекта)
2. [Basic Auth: как это работает](#basic-auth-как-это-работает)
3. [Точка входа: index.ts](#точка-входа-indexts)
4. [Clean Architecture](#clean-architecture)
5. [Domain Layer (Доменный слой)](#domain-layer-доменный-слой)
6. [Repository Layer (Слой данных)](#repository-layer-слой-данных)
7. [Use Case Layer (Бизнес-логика)](#use-case-layer-бизнес-логика)
8. [Delivery Layer (HTTP и Middleware)](#delivery-layer-http-и-middleware)
9. [Зависимости между слоями](#зависимости-между-слоями)
10. [Резюме](#резюме)

---

## Обзор проекта

Реализация **HTTP Basic Authentication** (RFC 7617) на **Node.js / TypeScript** (Express). Клиент передаёт логин и пароль в заголовке `Authorization` в каждом запросе. Регистрация — по JSON (email + password); доступ к `/me` и удаление аккаунта — только с Basic Auth. Пароли хешируются через bcrypt.

### Структура файлов (Node.js / TypeScript)

```
01-basic-auth/node/
├── openapi.json                         # OpenAPI спецификация (для Swagger UI)
├── src/
│   ├── index.ts                           # Точка входа, DI, роуты
│   ├── domain/
│   │   ├── user.ts                        # Интерфейс User
│   │   └── repository.ts                  # Интерфейс UserRepository
│   ├── repository/
│   │   └── memoryUserRepository.ts        # In-memory реализация
│   ├── usecase/
│   │   └── authUsecase.ts                 # Регистрация, логин, bcrypt
│   └── delivery/
│       ├── authHandler.ts                 # HTTP handlers
│       └── middleware/
│           └── basicAuth.ts               # Basic Auth middleware
├── package.json
├── tsconfig.json
└── ARCHITECTURE.md
```

---

## Swagger / OpenAPI (удобно для тестирования)

- Swagger UI: `http://localhost:8080/swagger`
- OpenAPI JSON: `http://localhost:8080/openapi.json`

## Basic Auth: как это работает

**RFC 7617:** клиент кодирует `email:password` в Base64 и отправляет:

```
Authorization: Basic <base64(email:password)>
```

Отдельного `POST /login` нет — авторизация при каждом запросе. Middleware декодирует заголовок, вызывает use case Login (bcrypt), при успехе кладёт user_id в request и передаёт в handler.

---

## Точка входа: index.ts

```typescript
const userRepository = new MemoryUserRepository();
const authUsecase = new AuthUsecase(userRepository);
const authHandler = createAuthHandler(authUsecase);
const basicAuth = basicAuthMiddleware(authUsecase);

const app = express();
app.use(express.json());

app.get('/health', (_req, res) => res.json({ status: 'ok' }));
app.post('/api/v1/auth/register', authHandler.register);
app.get('/api/v1/auth/me', basicAuth, authHandler.me);
app.delete('/api/v1/auth/me', basicAuth, authHandler.deleteUser);
```

Dependency Injection через конструкторы; middleware Express для защищённых роутов.

---

## Clean Architecture

Те же слои: Domain (интерфейсы), Use Case (бизнес-логика), Repository (данные), Delivery (Express handlers + middleware). Domain не зависит от фреймворка.

---

## Domain Layer (Доменный слой)

### domain/user.ts

```typescript
export interface User {
  id: string;
  email: string;
  password: string;
  createdAt: Date;
}
```

### domain/repository.ts

```typescript
export interface UserRepository {
  create(user: User): Promise<void>;
  getById(id: string): Promise<User | null>;
  getByEmail(email: string): Promise<User | null>;
  delete(id: string): Promise<void>;
}
```

---

## Repository Layer (Слой данных)

**memoryUserRepository.ts:** `Map<string, User>` для хранения. Методы create, getById, getByEmail, delete. Асинхронные для единообразия с будущими БД.

---

## Use Case Layer (Бизнес-логика)

**authUsecase.ts:**

- **register:** проверка по email, bcrypt.hash, создание User с uuid, сохранение
- **login:** getByEmail, bcrypt.compare — одинаковое сообщение при ошибке (защита от enumeration)
- **getUserById, deleteUserById:** вызов repository

---

## Delivery Layer (HTTP и Middleware)

**authHandler.ts:** register (JSON body), me (userId из req), deleteUser (userId из req).

**middleware/basicAuth.ts:** читает `Authorization: Basic ...`, декодирует Base64, вызывает authUsecase.login, при успехе кладёт userId в `req` и вызывает next().

---

## Зависимости между слоями

HTTP → Middleware → Handler → UseCase → Domain; Repository реализует интерфейс Domain. Use Case не знает про Express и заголовки.

---

## Резюме

**Реализовано:** Basic Auth (RFC 7617), bcrypt, Express middleware, Clean Architecture на TypeScript.

**Следующий проект (Node):** [02-api-key](../../02-api-key/node/ARCHITECTURE.md) — авторизация по X-API-Key.
