# Архитектура проекта: Authorization Backend на Go

## 📋 Содержание

1. [Обзор проекта](#обзор-проекта)
2. [Точка входа: main.go](#точка-входа-maingo)
3. [Clean Architecture](#clean-architecture)
4. [Domain Layer (Доменный слой)](#domain-layer-доменный-слой)
5. [Repository Layer (Слой данных)](#repository-layer-слой-данных)
6. [Use Case Layer (Бизнес-логика)](#use-case-layer-бизнес-логика)
7. [Delivery Layer (HTTP хендлеры)](#delivery-layer-http-хендлеры)
8. [Полный Flow запроса](#полный-flow-запроса)
9. [Зависимости между слоями](#зависимости-между-слоями)

---

## Обзор проекта

### Структура файлов

```
authorization_on_golang/
├── cmd/
│   └── server/
│       └── main.go                          # 🚀 Точка входа приложения
├── internal/
│   ├── domain/                              # 🎯 Бизнес-модели и интерфейсы
│   │   ├── user.go                         # Модель User
│   │   └── repository.go                   # Интерфейс UserRepository
│   ├── repository/                          # 💾 Слой работы с данными
│   │   └── memory/
│   │       └── user_repository.go          # In-memory реализация
│   ├── usecase/                            # 🧠 Бизнес-логика
│   │   └── auth_usecase.go                 # Use cases авторизации
│   └── delivery/                           # 🌐 HTTP слой
│       └── auth_handler.go                 # HTTP handlers
├── go.mod                                   # Зависимости проекта
└── go.sum                                   # Checksums зависимостей
```

---

## Точка входа: main.go

### Что мы создали в main.go

```mermaid
graph TD
    A[Запуск приложения] --> B[Инициализация слоев]
    B --> C[Repository Layer]
    B --> D[UseCase Layer]
    B --> E[Handler Layer]
    E --> F[HTTP Router]
    F --> G[HTTP Server]
    G --> H[Graceful Shutdown]

    style A fill:#e1f5ff
    style B fill:#fff4e1
    style G fill:#ffe1e1
```



### Пошаговое объяснение main.go

#### 1️⃣ **Health Check Handler**

```go
func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
```

**Зачем:**

- Проверка работоспособности сервера
- Для мониторинга в production (Kubernetes, Load Balancers)
- Быстрое тестирование: `curl http://localhost:8080/health`

#### 2️⃣ **Инициализация слоев (Dependency Injection)**

```go
// Создаем слой данных
userRepository := memory.NewUserRepository()

// Создаем бизнес-логику, передаем repository
authUsecase := usecase.NewAuthUsecase(userRepository)

// Создаем HTTP handlers, передаем usecase
authHandler := delivery.NewAuthHandler(authUsecase)
```

**Зачем:**

- **Инверсия зависимостей** (Dependency Injection)
- Каждый слой получает зависимости через конструктор
- Легко заменить `memory` на `postgres` - изменится только 1 строка
- Тестируемость - можно подставить моки

**Диаграмма зависимостей:**

```mermaid
graph LR
    Main[main.go] --> Repo[memory.UserRepository]
    Main --> UseCase[auth.UseCase]
    Main --> Handler[auth.Handler]

    UseCase -.depends on.-> RepoInterface[domain.UserRepository]
    Handler -.depends on.-> UseCase
    Repo -.implements.-> RepoInterface

    style Main fill:#e1f5ff
    style RepoInterface fill:#d4edda
```



#### 3️⃣ **HTTP Server конфигурация**

```go
server := &http.Server{
    Addr:         ":8080",              // Порт
    ReadTimeout:  15 * time.Second,     // Timeout чтения запроса
    WriteTimeout: 15 * time.Second,     // Timeout отправки ответа
    IdleTimeout:  15 * time.Second,     // Timeout idle соединений
    Handler:      mux,                  // Наш роутер
}
```

**Зачем таймауты:**

- **ReadTimeout** - защита от медленных клиентов (Slowloris атака)
- **WriteTimeout** - ограничение времени обработки запроса
- **IdleTimeout** - закрытие неактивных keep-alive соединений

#### 4️⃣ **Регистрация endpoints**

```go
mux.HandleFunc("/health", healthHandler)
mux.HandleFunc("POST /api/v1/auth/register", authHandler.RegisterHandler)
mux.HandleFunc("POST /api/v1/auth/login", authHandler.LoginHandler)
mux.HandleFunc("GET /api/v1/auth/me", authHandler.MeHandler)
```

**Паттерн:** Go 1.22+ поддерживает HTTP методы в роутах

- `POST /api/v1/auth/register` - только POST запросы
- Раньше нужно было: `if r.Method != "POST" { ... }`

#### 5️⃣ **Graceful Shutdown**

```go
// Запуск сервера в goroutine
go func() {
    log.Println("Сервер запущен на http://localhost:8080")
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("Ошибка запуска сервера: %v", err)
    }
}()

// Ожидание сигнала остановки
quit := make(chan os.Signal, 1)
signal.Notify(quit, os.Interrupt)  // Ctrl+C
<-quit

// Плавная остановка
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
if err := server.Shutdown(ctx); err != nil {
    log.Fatalf("Ошибка завершения работы сервера: %v", err)
}
```

**Flow Graceful Shutdown:**

```mermaid
sequenceDiagram
    participant User
    participant Main
    participant Server
    participant Requests

    User->>Main: Ctrl+C (SIGINT)
    Main->>Server: Shutdown(ctx)
    Server->>Server: Перестать принимать новые запросы
    Server->>Requests: Ждать завершения текущих (10 сек)
    Requests-->>Server: Завершены
    Server-->>Main: Shutdown complete
    Main->>Main: Выход
```



**Зачем:**

1. **Целостность данных** - запросы не обрываются посередине
2. **Пользователи не видят ошибок** - все запросы завершаются корректно
3. **Production стандарт** - Kubernetes/Docker ожидают graceful shutdown

---

## Clean Architecture

### Что это и зачем

**Clean Architecture** - это подход к организации кода, где:

- Бизнес-логика **не зависит** от деталей (БД, HTTP, UI)
- Зависимости направлены **внутрь** (к бизнес-логике)
- Легко тестировать и менять компоненты

### Слои нашего приложения

```mermaid
graph TB
    subgraph External["🌐 External Layer (Внешний мир)"]
        HTTP[HTTP Requests]
        CLI[CLI Commands]
        gRPC[gRPC Calls]
    end

    subgraph Delivery["📬 Delivery Layer (Адаптеры)"]
        Handler[auth_handler.go<br/>Обрабатывает HTTP запросы]
    end

    subgraph UseCase["🧠 Use Case Layer (Бизнес-правила)"]
        Auth[auth_usecase.go<br/>Логика регистрации/входа]
    end

    subgraph Domain["🎯 Domain Layer (Сердце приложения)"]
        User[user.go<br/>Модель User]
        Repo[repository.go<br/>Интерфейс UserRepository]
    end

    subgraph Repository["💾 Repository Layer (Работа с данными)"]
        Memory[user_repository.go<br/>In-memory хранилище]
        Postgres[PostgreSQL<br/>будет позже]
    end

    HTTP --> Handler
    CLI -.-> Handler
    gRPC -.-> Handler
    Handler --> Auth
    Auth --> Repo
    Memory -.implements.-> Repo
    Postgres -.implements.-> Repo
    Auth --> User

    style Domain fill:#d4edda
    style UseCase fill:#fff3cd
    style Delivery fill:#f8d7da
    style Repository fill:#d1ecf1
```



### Правила зависимостей

```mermaid
graph TD
    A[Domain Layer<br/>НЕ ЗАВИСИТ ни от чего] --> B[Use Case Layer<br/>Зависит только от Domain]
    B --> C[Delivery/Repository Layers Зависят от Use Case и Domain]

    style A fill:#d4edda
    style B fill:#fff3cd
    style C fill:#f8d7da
```



**Ключевой момент:**

- `UseCase` знает о `UserRepository` **интерфейсе** (domain)
- `UseCase` НЕ знает о `memory.UserRepository` (конкретная реализация)
- Это позволяет менять БД без изменения бизнес-логики

---

## Domain Layer (Доменный слой)

### 🎯 domain/user.go - Модель пользователя

```go
type User struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Password  string    `json:"-"`           // НЕ возвращается в JSON
    CreatedAt time.Time `json:"created_at"`
}
```

**Зачем каждое поле:**

- `ID` - уникальный идентификатор (UUID)
- `Name` - имя пользователя (опционально)
- `Email` - уникальный email для входа
- `Password` - хеш пароля (НЕ отдается клиенту: `json:"-"`)
- `CreatedAt` - дата регистрации

**Почему `json:"-"` для Password:**

```json
// БЕЗ json:"-":
{"id":"123","email":"test@test.com","password":"hash123"}  ❌ Утечка хеша!

// С json:"-":
{"id":"123","email":"test@test.com"}  ✅ Безопасно
```

### 🎯 domain/repository.go - Интерфейс хранилища

```go
type UserRepository interface {
    Create(user *User) error
    GetByID(id string) (*User, error)
    GetByEmail(email string) (*User, error)
}
```

**Зачем интерфейс:**

- **Абстракция** - бизнес-логика не знает ГДЕ хранятся данные
- **Подменяемость** - легко менять реализацию (memory → postgres → redis)
- **Тестируемость** - можно создать mock для тестов

**Пример гибкости:**

```mermaid
graph LR
    UC[UseCase] --> Interface[UserRepository<br/>interface]

    Interface -.-> Memory[Memory<br/>implementation]
    Interface -.-> Postgres[PostgreSQL<br/>implementation]
    Interface -.-> Redis[Redis<br/>implementation]
    Interface -.-> Mock[Mock<br/>for tests]

    style Interface fill:#d4edda
    style Memory fill:#d1ecf1
    style Postgres fill:#d1ecf1
    style Redis fill:#d1ecf1
    style Mock fill:#f8d7da
```



Меняется **1 строка в main.go**, вся остальная логика остается!

---

## Repository Layer (Слой данных)

### 💾 repository/memory/user_repository.go

```go
type UserRepository struct {
    users map[string]*domain.User  // map[UserID]*User
    mu    sync.RWMutex             // Защита от race conditions
}
```

*Почему map[string]User:*

- Ключ = `ID` пользователя
- Значение = указатель на `User`
- Быстрый доступ: O(1) по ID

**Зачем sync.RWMutex:**

```mermaid
sequenceDiagram
    participant G1 as Goroutine 1
    participant G2 as Goroutine 2
    participant Map as users map

    G1->>Map: Write (Register)
    G2->>Map: Read (Login)
    Note over Map: ❌ Race condition!<br/>Одновременное чтение/запись

    rect rgb(255, 200, 200)
        Note over G1,Map: БЕЗ мьютекса = ПАНИКА
    end
```



**С мьютексом:**

```mermaid
sequenceDiagram
    participant G1 as Goroutine 1
    participant G2 as Goroutine 2
    participant Mutex as RWMutex
    participant Map as users map

    G1->>Mutex: Lock() для записи
    Mutex->>G1: OK
    G1->>Map: Write
    G2->>Mutex: RLock() для чтения
    Mutex->>G2: ⏳ Ждите...
    G1->>Map: Write done
    G1->>Mutex: Unlock()
    Mutex->>G2: OK, читайте
    G2->>Map: Read

    rect rgb(200, 255, 200)
        Note over G1,Map: С мьютексом = безопасно ✅
    end
```



### Методы Repository

#### 1️⃣ Create - создание пользователя

```go
func (r *UserRepository) Create(user *domain.User) error {
    r.mu.Lock()           // Эксклюзивная блокировка для записи
    defer r.mu.Unlock()   // Разблокировка при выходе

    if _, exists := r.users[user.ID]; exists {
        return errors.New("user already exists")
    }
    r.users[user.ID] = user
    return nil
}
```

**Зачем Lock (не RLock):**

- `Lock()` - эксклюзивная блокировка (никто не может читать/писать)
- `RLock()` - shared блокировка (можно читать, но не писать)
- Запись требует **полной** блокировки

#### 2️⃣ GetByID - поиск по ID

```go
func (r *UserRepository) GetByID(id string) (*domain.User, error) {
    r.mu.RLock()          // Shared блокировка для чтения
    defer r.mu.RUnlock()

    user, exists := r.users[id]
    if !exists {
        return nil, errors.New("user not found")
    }
    return user, nil
}
```

**Быстрый поиск:**

- Map lookup: O(1)
- Прямой доступ по ключу `id`

#### 3️⃣ GetByEmail - поиск по email

```go
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    for _, user := range r.users {  // Перебираем всех пользователей
        if user.Email == email {
            return user, nil
        }
    }
    return nil, errors.New("user not found")
}
```

**Почему цикл:**

- Map хранит по ID, не по email
- Нужно перебрать все значения: O(n)
- В PostgreSQL будет индекс на email → O(1)

**Оптимизация (не реализована):**

```go
// Можно добавить второй map:
emailToID map[string]string  // map[email]userID
// Тогда GetByEmail будет O(1)
```

---

## Use Case Layer (Бизнес-логика)

### 🧠 usecase/auth_usecase.go

```go
type AuthUsecase struct {
    userRepository domain.UserRepository  // Интерфейс, НЕ конкретная реализация!
}
```

**Ключевой момент:**

- Зависимость от **интерфейса** `domain.UserRepository`
- НЕ от `*memory.UserRepository`
- Можно подставить любую реализацию

### Use Cases (сценарии использования)

#### 1️⃣ Register - регистрация пользователя

```mermaid
flowchart TD
    A[Клиент: Register email, password] --> B{Email уже<br/>существует?}
    B -->|Да| C[❌ Ошибка:<br/>user already exists]
    B -->|Нет| D[Создать User]
    D --> E[Генерировать UUID]
    D --> F[Установить CreatedAt]
    D --> G[Сохранить в Repository]
    G --> H[✅ Вернуть User]

    style C fill:#f8d7da
    style H fill:#d4edda
```



**Код:**

```go
func (u *AuthUsecase) Register(email, password string) (*domain.User, error) {
    // 1. Проверяем, существует ли пользователь
    existingUser, err := u.userRepository.GetByEmail(email)
    if err == nil && existingUser != nil {
        return nil, errors.New("user already exists")
    }

    // 2. Создаем нового пользователя
    user := &domain.User{
        ID:        uuid.New().String(),  // Генерируем UUID
        Email:     email,
        Password:  password,              // Пока plain text (исправим в этапе 3)
        CreatedAt: time.Now(),
    }

    // 3. Сохраняем в repository
    err = u.userRepository.Create(user)
    if err != nil {
        return nil, err
    }

    return user, nil
}
```

**Важные детали:**

**UUID генерация:**

```go
uuid.New().String()
// Пример: "996b2915-2e8a-4bae-98fc-7e3273727cca"
```

- Уникальный ID (вероятность коллизии ~0)
- Не sequential (безопасность - нельзя угадать следующий ID)
- Стандарт для распределенных систем

**Логика проверки существования:**

```go
existingUser, err := u.userRepository.GetByEmail(email)
if err == nil && existingUser != nil {  // Пользователь НАЙДЕН
    return nil, errors.New("user already exists")
}
// err != nil означает "не найден" → можно регистрировать
```

#### 2️⃣ Login - вход в систему

```mermaid
flowchart TD
    A[Клиент: Login email, password] --> B[Найти User по email]
    B --> C{User<br/>найден?}
    C -->|Нет| D[❌ invalid email or password]
    C -->|Да| E{Пароль<br/>совпадает?}
    E -->|Нет| D
    E -->|Да| F[✅ Вернуть User]

    style D fill:#f8d7da
    style F fill:#d4edda
```



**Код:**

```go
func (u *AuthUsecase) Login(email, password string) (*domain.User, error) {
    // 1. Найти пользователя
    user, err := u.userRepository.GetByEmail(email)
    if err != nil {
        return nil, errors.New("invalid email or password")  // Не раскрываем детали
    }

    // 2. Проверить пароль
    if user.Password != password {
        return nil, errors.New("invalid email or password")  // То же сообщение
    }

    return user, nil
}
```

**Почему одинаковое сообщение об ошибке:**

```go
// ❌ ПЛОХО:
"user not found"       // Атакующий знает, что email не зарегистрирован
"invalid password"     // Атакующий знает, что email существует

// ✅ ХОРОШО:
"invalid email or password"  // Неясно, что именно не так
```

Это защита от **user enumeration** атак.

#### 3️⃣ GetUserByID - получить пользователя

```go
func (u *AuthUsecase) GetUserByID(id string) (*domain.User, error) {
    user, err := u.userRepository.GetByID(id)
    if err != nil {
        return nil, err
    }
    return user, nil
}
```

**Простой pass-through** к repository.
В будущем здесь может быть:

- Проверка прав доступа
- Логирование
- Кеширование

---

## Delivery Layer (HTTP хендлеры)

### 🌐 delivery/auth_handler.go

```go
type AuthHandler struct {
    authUsecase *usecase.AuthUsecase
}
```

**Ответственность:**

1. Парсинг HTTP запросов (JSON → Go структуры)
2. Вызов use case
3. Формирование HTTP ответов (Go структуры → JSON)
4. Установка правильных HTTP статус-кодов

### Request/Response DTO

```go
type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

**Зачем отдельные структуры:**

- Валидация входных данных
- Отделение API контракта от доменных моделей
- Можно добавить валидацию (этап 8): `validate:"required,email"`

### HTTP Handlers

#### 1️⃣ RegisterHandler

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant UseCase
    participant Repository

    Client->>Handler: POST /register<br/>{email, password}
    Handler->>Handler: Decode JSON
    Handler->>UseCase: Register(email, password)
    UseCase->>Repository: GetByEmail()
    Repository-->>UseCase: not found
    UseCase->>Repository: Create(user)
    Repository-->>UseCase: success
    UseCase-->>Handler: user
    Handler->>Handler: Encode JSON
    Handler-->>Client: 201 Created<br/>{id, email, ...}
```



**Код:**

```go
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Парсинг JSON
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)  // 400
        return
    }

    // 2. Вызов use case
    user, err := h.authUsecase.Register(req.Email, req.Password)
    if err != nil {
        if err.Error() == "user already exists" {
            http.Error(w, err.Error(), http.StatusConflict)  // 409
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)  // 500
        return
    }

    // 3. Успешный ответ
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)  // 201
    json.NewEncoder(w).Encode(user)
}
```

**HTTP статус-коды:**

- `400 Bad Request` - невалидный JSON
- `409 Conflict` - email уже существует
- `500 Internal Server Error` - неожиданная ошибка
- `201 Created` - успешная регистрация

#### 2️⃣ LoginHandler

```go
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)  // 400
        return
    }

    user, err := h.authUsecase.Login(req.Email, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)  // 401
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)  // 200
}
```

**Статус-коды:**

- `401 Unauthorized` - неверный email или пароль
- `200 OK` - успешный вход

#### 3️⃣ MeHandler

```go
func (h *AuthHandler) MeHandler(w http.ResponseWriter, r *http.Request) {
    // Получаем user_id из query параметра
    userID := r.URL.Query().Get("user_id")
    if userID == "" {
        http.Error(w, "user_id is required", http.StatusBadRequest)  // 400
        return
    }

    user, err := h.authUsecase.GetUserByID(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)  // 404
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)  // 200
}
```

**Пример запроса:**

```bash
GET /api/v1/auth/me?user_id=996b2915-2e8a-4bae-98fc-7e3273727cca
```

**Почему временное решение:**

- В production user_id будет браться из JWT токена (этап 5)
- Сейчас для простоты передаем в query параметре

---

## Полный Flow запроса

### Пример: Регистрация пользователя

```mermaid
sequenceDiagram
    autonumber
    participant Client as 🖥️ Client
    participant Main as main.go
    participant Handler as 📬 AuthHandler
    participant UseCase as 🧠 AuthUsecase
    participant Repo as 💾 UserRepository
    participant Memory as 🗄️ Memory Storage

    Client->>Main: POST /api/v1/auth/register<br/>{email, password}
    Main->>Handler: RegisterHandler()

    rect rgb(255, 240, 240)
        Note over Handler: Delivery Layer
        Handler->>Handler: Decode JSON
        Handler->>UseCase: Register(email, password)
    end

    rect rgb(255, 255, 230)
        Note over UseCase: Use Case Layer
        UseCase->>Repo: GetByEmail(email)
        Repo->>Memory: Search in map
        Memory-->>Repo: not found
        Repo-->>UseCase: error

        UseCase->>UseCase: Create User<br/>ID=uuid.New()<br/>CreatedAt=now()
        UseCase->>Repo: Create(user)
    end

    rect rgb(230, 240, 255)
        Note over Repo,Memory: Repository Layer
        Repo->>Memory: Lock()<br/>users[id] = user
        Memory-->>Repo: success
        Repo-->>UseCase: nil
    end

    rect rgb(255, 240, 240)
        UseCase-->>Handler: user
        Handler->>Handler: Encode JSON<br/>Set status 201
        Handler-->>Main: Response
    end

    Main-->>Client: 201 Created<br/>{id, email, created_at}
```



### Что происходит на каждом шаге:

1. **Client** отправляет POST запрос с JSON
2. **main.go** роутер направляет на `RegisterHandler`
3. **Handler** декодирует JSON в `RegisterRequest`
4. **Handler** вызывает `UseCase.Register()`
5. **UseCase** проверяет существование через `Repository.GetByEmail()`
6. **Repository** ищет в map, не находит → возвращает error
7. **UseCase** создает новый `User` с UUID и временем
8. **UseCase** вызывает `Repository.Create()`
9. **Repository** блокирует map, добавляет пользователя
10. **Memory** сохраняет в `map[id]*User`
11. **Repository** возвращает success
12. **UseCase** возвращает созданного `User`
13. **Handler** кодирует в JSON, ставит статус 201
14. **Client** получает ответ с данными пользователя

---

## Зависимости между слоями

### Направление зависимостей

```mermaid
graph TB
    subgraph External["🌍 Внешний мир"]
        HTTP[HTTP Client]
    end

    subgraph Delivery["📬 Delivery Layer"]
        Handler[auth_handler.go]
    end

    subgraph UseCase["🧠 Use Case Layer"]
        Auth[auth_usecase.go]
    end

    subgraph Domain["🎯 Domain Layer<br/>(Независимый центр)"]
        User[user.go]
        Interface[repository.go<br/>interface]
    end

    subgraph Repository["💾 Repository Layer"]
        Memory[memory/user_repository.go]
    end

    HTTP --> Handler
    Handler --> Auth
    Auth --> Interface
    Auth --> User
    Memory -.implements.-> Interface

    style Domain fill:#d4edda,stroke:#333,stroke-width:3px
    style Interface fill:#ffffcc
```



### Инверсия зависимостей (Dependency Inversion)

**Традиционный подход (❌):**

```mermaid
graph TD
    UC[UseCase] --> Mem[memory.UserRepository]
    UC --> Pg[postgres.UserRepository]

    Note[UseCase зависит от конкретных реализаций]

    style UC fill:#f8d7da
```



**Clean Architecture подход (✅):**

```mermaid
graph TD
    UC[UseCase] --> Interface[domain.UserRepository<br/>interface]
    Mem[memory.UserRepository] -.implements.-> Interface
    Pg[postgres.UserRepository] -.implements.-> Interface

    Note[UseCase зависит только от абстракции]

    style UC fill:#d4edda
    style Interface fill:#ffffcc
```



**Преимущества:**

1. **UseCase не меняется** при смене БД
2. **Легко тестировать** - подставляем mock
3. **Гибкость** - можно иметь несколько реализаций одновременно

### Пример: Смена БД (memory → PostgreSQL)

**Что нужно изменить:**

```go
// 1. Было в main.go:
userRepository := memory.NewUserRepository()

// 2. Станет:
userRepository := postgres.NewUserRepository(db)

// 3. ВСЁ! Больше ничего не меняется 🎉
```

**Остается без изменений:**

- ✅ `auth_usecase.go` - ни одной строки
- ✅ `auth_handler.go` - ни одной строки
- ✅ `domain/user.go` - ни одной строки
- ✅ `domain/repository.go` - ни одной строки

---

## Преимущества нашей архитектуры

### 1. Тестируемость

```go
// Можем создать mock:
type MockUserRepository struct {}

func (m *MockUserRepository) Create(user *domain.User) error {
    // Контролируемое поведение для тестов
    return nil
}

// Тест:
mockRepo := &MockUserRepository{}
usecase := NewAuthUsecase(mockRepo)
user, err := usecase.Register("test@test.com", "pass123")
// Проверяем логику без реальной БД
```

### 2. Гибкость

```mermaid
graph TB
    Main[main.go<br/>Единственное место изменений]

    Main --> Memory[memory.UserRepository<br/>Разработка]
    Main --> Postgres[postgres.UserRepository<br/>Production]
    Main --> Redis[redis.UserRepository<br/>Кеш]
    Main --> Mock[mock.UserRepository<br/>Тесты]

    style Main fill:#ffffcc
    style Memory fill:#d1ecf1
    style Postgres fill:#d1ecf1
    style Redis fill:#d1ecf1
    style Mock fill:#f8d7da
```



### 3. Читаемость

Каждый файл имеет **одну ответственность**:

- `user.go` - знает, ЧТО такое пользователь
- `auth_usecase.go` - знает, КАК регистрировать/входить
- `auth_handler.go` - знает, КАК обрабатывать HTTP
- `user_repository.go` - знает, ГДЕ хранить данные

### 4. Масштабируемость

```mermaid
graph TD
    A[Сейчас: 3 use cases<br/>Register, Login, GetByID]
    A --> B[Добавить: EmailVerification<br/>PasswordReset<br/>OAuth]
    B --> C[Каждый use case<br/>в отдельном файле]
    C --> D[Не трогаем существующий код<br/>Open/Closed Principle]

    style A fill:#d1ecf1
    style D fill:#d4edda
```



---

## Что дальше?

### Этап 3: Хеширование паролей (bcrypt)

Сейчас пароли хранятся в открытом виде:

```go
Password: password  // ❌ Опасно!
```

Станет:

```go
Password: "$2a$10$N9qo8uLO..."  // ✅ Хеш bcrypt
```

### Этап 5: JWT токены

Вместо передачи `user_id` в query:

```go
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Этап 7: PostgreSQL

```go
// Вместо map:
users map[string]*User

// Станет:
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR UNIQUE,
    ...
);
```

---

## Резюме

### Что мы построили:

1. **main.go** - точка входа с Graceful Shutdown
2. **Domain Layer** - независимое ядро (User, интерфейсы)
3. **Repository Layer** - работа с данными (in-memory)
4. **Use Case Layer** - бизнес-логика (Register, Login, GetByID)
5. **Delivery Layer** - HTTP handlers (JSON ↔ Go)

### Почему это Clean Architecture:

✅ **Независимость от фреймворков** - используем стандартную библиотеку  
✅ **Тестируемость** - можем мокировать любой слой  
✅ **Независимость от БД** - легко меняем память на PostgreSQL  
✅ **Независимость от UI** - легко добавить gRPC рядом с HTTP  
✅ **Бизнес-правила не знают о деталях** - UseCase не знает про HTTP/БД

### Ключевые принципы:

1. **Dependency Injection** - зависимости передаются через конструкторы
2. **Interface Segregation** - маленькие, специфичные интерфейсы
3. **Dependency Inversion** - зависимость от абстракций, не реализаций
4. **Single Responsibility** - каждый файл/функция делает что-то одно

---

**Архитектура готова к масштабированию! 🚀**