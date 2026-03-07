# EGO

**Ego** is a lightweight Go core framework designed for microservices and backend systems. It provides a structured foundation to rapidly build scalable, maintainable, and observable applications with minimal boilerplate.

## Features

- ⚙️ **Configuration Management** — Viper-powered TOML config with env override
- 🌐 **REST Server** — Fiber-based HTTP router with groups, CORS, rate limiting, and auth middleware
- 📦 **Database Integration** — MySQL wrapper with transaction management and query helpers
- 🗄️ **Cache Manager** — Redis-backed cache with singleflight deduplication
- 🔐 **Authentication** — Service-to-service token management with auto-refresh
- 🧠 **Tracing** — Sentry APM integration with configurable sample rate
- 📄 **Error Codes** — Standardized codes mapping to HTTP/gRPC status
- 📝 **Structured Logging** — Leveled logger with field support
- 📡 **HTTP Client** — Request wrapper with retry, headers, and auth token injection
- 🔑 **Cryptography** — AES-GCM, RSA, bcrypt, SHA-256, and PKCE utilities
- ✅ **Validator** — Struct validation using `go-playground/validator`
- 🎲 **Generator** — Random and cryptographically secure string generation
- 🧪 **Stubs** — Function type stubs for testability

---

## Getting Started

### Installation

```bash
go get github.com/evenyosua18/ego
```

### Minimal REST Server

```go
package main

import (
    "github.com/evenyosua18/ego/app"
    "github.com/evenyosua18/ego/http"
)

func main() {
    http.RegisterRouteByGroup("public", []http.RouteFunc{
        func(r http.IHttpRouter) {
            r.Get("/ping", func(ctx http.Context) error {
                return ctx.ResponseSuccess(map[string]string{"message": "pong"})
            })
        },
    })

    app.GetApp().RunRest()
}
```

---

## Configuration

Ego reads from a `config.toml` file using [Viper](https://github.com/spf13/viper). By default, it looks in `./config/config.toml`. Override the path with env vars:

| Env Variable    | Description                      | Default      |
|-----------------|----------------------------------|--------------|
| `CONFIG_PATH`   | Path to config directory         | `./config`   |
| `CONFIG_ROOT`   | Root path (overrides CONFIG_PATH)| —            |
| `CONFIG_DIR`    | Subdirectory appended to path    | —            |

Environment variables can also override config keys using `_` as separator (e.g. `SERVICE_NAME` overrides `service.name`).

### Sample `config.toml`

```toml
[service]
name = "my-service"
env = "local"                     # default: "local"
shutdown_timeout = "30s"          # default: 30s

[logger]
level = "info"                    # debug | info | warn | error | fatal

[database]
driver = "mysql"                  # default: mysql
name = "mydb"                     # required (leave empty to skip DB)
host = "localhost"                # default: localhost
port = "3306"                     # default: 3306
user = "root"                     # default: root
password = "secret"               # required if env != "local"
protocol = "tcp"                  # default: tcp
params = ""                       # extra DSN params
max_open_conns = 100              # default: 100
max_idle_conns = 20               # default: 20
conn_max_lifetime = "30m"         # default: 30m
conn_max_idle_time = "5m"         # default: 5m

[tracer]
dsn = "https://xxx@sentry.io/123" # leave empty to disable
sample_rate = 1.0                 # default: 1.0
flush_time = "1"                  # seconds, default: 1

[router]
prefix = ""                       # URL prefix
port = ":8080"                    # default: :8080
show_registered = true            # print routes on startup
html_path = ""                    # path for HTML templates
doc_path = ""                     # swagger JSON path
disable_auth_checker = false      # skip auth validation
max_connection = 5000             # default: 5000
rate_limit.max_limit = 100        # default: 100 (0 = disabled)
read_timeout = "30s"              # default: 30s
write_timeout = "0s"              # default: 0 (no timeout)
idle_timeout = "0s"               # default: 0 (no timeout)
allow_origins = ["*"]
allow_methods = ["GET", "POST"]
allow_headers = ["Content-Type"]
allow_credentials = false

[cache.redis]
addr = "localhost:6379"
password = ""
db = 0
max_retries = 3
min_idle_conns = 10
pool_size = 10
idle_timeout = "5m"
read_timeout = "5s"
write_timeout = "5s"
dial_timeout = "5s"
pool_timeout = "5s"

[code]
filename = "codes.yaml"           # custom error codes file

[auth_svc]
base_url = "http://auth-svc:8080/"
client_id = "my-client-id"
client_secret = "my-client-secret"
```

---

## Package Reference

### `config` — Configuration Management

Provides a singleton Viper wrapper. Reads `config.toml` on first access.

```go
import "github.com/evenyosua18/ego/config"

// Read values
name := config.GetConfig().GetString("service.name")
port := config.GetConfig().GetInt("database.port")
flag := config.GetConfig().GetBool("router.show_registered")

// For unit tests — inject config values without a file
config.SetTestConfig(map[string]any{
    "service.name": "test-svc",
    "database.name": "testdb",
})
```

---

### `http` — REST Server & Router

Fiber-based HTTP router with built-in middleware (panic recovery, logging, rate limiting, CORS) and optional auth validation.

#### Registering Routes

Routes are organized into **groups** (`"public"`, `"svc"`, or custom). Public routes are always registered at `/`, other groups at `/<group-name>`.

```go
import "github.com/evenyosua18/ego/http"

// Register public routes
http.RegisterRouteByGroup("public", []http.RouteFunc{
    func(r http.IHttpRouter) {
        r.Get("/users", listUsers)
        r.Post("/users", createUser)
        r.Get("/users/:id", getUser)
        r.Put("/users/:id", updateUser)
        r.Delete("/users/:id", deleteUser)
    },
})

// Register service-to-service routes
http.RegisterRouteByGroup("svc", []http.RouteFunc{
    func(r http.IHttpRouter) {
        r.Post("/sync", syncData)
    },
})
```

#### Handler Function

Handlers receive an `http.Context` which wraps Fiber's context:

```go
func getUser(ctx http.Context) error {
    id := ctx.Param("id")         // path param
    page := ctx.Query("page")     // query param

    // Bind all query params to struct
    var filter struct {
        Page  int    `query:"page"`
        Limit int    `query:"limit"`
        Name  string `query:"name"`
    }
    ctx.BindQuery(&filter)

    // Parse JSON body
    var req CreateUserRequest
    ctx.RequestBody(&req)

    // Response helpers
    return ctx.ResponseSuccess(user)     // 200 JSON
    return ctx.ResponseError(err)        // auto maps code.Code to HTTP status
    return ctx.JSON(201, data)           // custom status
    return ctx.ResponseRedirect("/login")
}
```

#### Route Options (Roles & Middleware)

```go
r.Get("/admin/dashboard", adminHandler,
    http.SetRouterRolesOption([]string{"admin", "superadmin"}),
    http.SetMiddleware(myCustomMiddleware),
)
```

#### Sub-groups

```go
func(r http.IHttpRouter) {
    v1 := r.Group("/v1")
    v1.Get("/items", listItems)
    v1.Post("/items", createItem)
}
```

---

### `code` — Standardized Error Codes

Maps custom error codes to HTTP/gRPC statuses with user-facing messages.

#### Built-in Codes

| Code | HTTP | Description |
|------|------|-------------|
| `internal_error` | 500 | Internal server error |
| `not_found_error` | 404 | Data not found |
| `database_error` | 500 | Database error |
| `bad_request` | 400 | Bad request |
| `unauthorized` | 401 | Unauthorized |
| `rate_limit_error` | 429 | Too many requests |
| `panic_error` | 500 | Panic recovery |
| `encryption_error` | 500 | Crypto error |
| `cache_error` | 500 | Cache error |
| `cache_not_found` | 404 | Cache miss |

#### Usage

```go
import "github.com/evenyosua18/ego/code"

// Wrap an error with a code
err := code.Wrap(dbErr, code.DatabaseError)

// Get a code (copy)
c := code.Get(code.BadRequestError).SetMessage("email is required")

// Extract code from error
c := code.Extract(err)
fmt.Println(c.Code(), c.CodeHTTP(), c.Message(), c.Error())

// Check if error matches a code
if code.Is(err, code.NotFoundError) {
    // handle not found
}
```

#### Custom Codes via YAML

Create a `codes.yaml` file and set `code.filename` in config:

```yaml
codes:
  - code: "duplicate_email"
    message: "Email already registered"
    error: "duplicate email"
    http_code: 409
    grpc_code: 6
```

---

### `sqldb` — Database Layer

MySQL wrapper with interfaces for testability, transaction management, and query helpers using [Scany](https://github.com/georgysavva/scany).

#### Querying

```go
import "github.com/evenyosua18/ego/sqldb"

// Get the DB manager (handles db/tx from context)
mgr := sqldb.GetDbManager()

// Inject DB into context (required before GetExecutor)
ctx, err := mgr.SetDBContext(ctx)

// Get executor (returns tx if in transaction, otherwise db)
exec, err := mgr.GetExecutor(ctx)

// Single row
row := exec.QueryRow("SELECT name FROM users WHERE id = ?", userID)
err := row.Scan(&name)

// Multiple rows with Scany
rows, err := exec.Query("SELECT * FROM users WHERE active = ?", true)
var users []User
err = rows.ScanAll(&users)

// Single row with Scany
var user User
err = rows.ScanOne(&user)

// Execute (INSERT, UPDATE, DELETE)
result, err := exec.Exec("INSERT INTO users (name) VALUES (?)", name)
id, _ := result.LastInsertId()
affected, _ := result.RowsAffected()
```

#### Transactions

```go
// Manual transaction
tx, txCtx, err := mgr.BeginTx(ctx)
defer tx.EndTx(&err)  // auto commit/rollback based on err

exec, _ := mgr.GetExecutor(txCtx)
_, err = exec.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromID)
_, err = exec.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toID)

// Wrapped transaction (recommended)
err := mgr.WrappedBeginTrx(ctx, func(txCtx context.Context) error {
    exec, _ := mgr.GetExecutor(txCtx)
    _, err := exec.Exec("INSERT INTO orders (user_id) VALUES (?)", userID)
    return err
})
```

---

### `cache` — Cache Manager

Redis-backed caching with [singleflight](https://pkg.go.dev/golang.org/x/sync/singleflight) to prevent thundering herd.

```go
import "github.com/evenyosua18/ego/cache"

// Get-or-set with type safety (generic)
user, err := cache.GetOrSet(ctx, cache.GetManager(), "user:123", 10*time.Minute,
    func() (User, error) {
        return repo.FindUserByID(ctx, 123)
    },
)

// Invalidate cache keys
err := cache.GetManager().Invalidate(ctx, "user:123", "user:list")
```

---

### `request` — HTTP Client

Wraps `net/http` with JSON marshaling, retry logic, and functional options.

```go
import "github.com/evenyosua18/ego/request"

client := request.NewClient(nil) // uses http.DefaultClient

// Simple GET
resp, body, err := client.Get(ctx, "https://api.example.com/data",
    request.WithHeaders(map[string]string{"X-Api-Key": "abc"}),
    request.WithQueryParams(map[string]string{"page": "1"}),
)

// POST with JSON body and retry
resp, body, err := client.Post(ctx, "https://api.example.com/resource",
    request.WithBody(CreateRequest{Name: "test"}),
    request.WithRetry(3, 2*time.Second, 500, 502, 503),
    request.WithAuthToken(auth.GetAccessToken()),
)

// Low-level Do
resp, body, err := client.Do(ctx, request.Request{
    Method:      http.MethodPut,
    URL:         "https://api.example.com/item/1",
    Body:        updatePayload,
    Headers:     map[string]string{"X-Request-ID": "abc"},
    QueryParams: map[string]string{"version": "2"},
})
```

---

### `auth` — Token Management

Manages service-to-service access tokens with automatic background refresh.

```go
import "github.com/evenyosua18/ego/auth"

// Auto-managed by app.RunRest() if auth_svc config is set.

// Get current token for outgoing requests
token := auth.GetAccessToken()

// Validate incoming token
claims, err := auth.ValidateToken(ctx, "Bearer xxx...")
// claims.UserID, claims.Roles, claims.Scopes, etc.
```

---

### `logger` — Structured Logging

Leveled logger with structured fields. Replaceable via the `Logger` interface.

```go
import "github.com/evenyosua18/ego/logger"

logger.Debug("processing request", logger.Field{Key: "user_id", Value: 42})
logger.Info("server started")
logger.Warn("connection pool high")
logger.Error(fmt.Errorf("failed to connect"))

// Create a child logger with persistent fields
l := logger.L().With(logger.Field{Key: "component", Value: "auth"})
l.Info("validating token")

// Set custom logger implementation
logger.SetLogger(myCustomLogger)
```

---

### `tracer` — Distributed Tracing

Sentry APM integration with span management. Auto-configured by `app.RunRest()`.

```go
import "github.com/evenyosua18/ego/tracer"

// Start a span (auto-starts a transaction if none exists)
sp := tracer.StartSpan(ctx, "UserService.GetUser",
    tracer.WithRequest(req),
    tracer.WithAttribute("user_id", userID),
)
defer sp.End()

// Log errors and responses
if err != nil {
    return sp.LogError(err)
}
sp.LogResponse(result)

// Start span and get new context
sp, newCtx := tracer.StartSpanWithContext(ctx, "FetchData")
defer sp.End()
```

---

### `cryptox` — Cryptography Utilities

```go
import "github.com/evenyosua18/ego/cryptox"

// AES-GCM (key must be base64-encoded)
encrypted, err := cryptox.EncryptAES(base64Key, "plaintext")
decrypted, err := cryptox.DecryptAES(base64Key, encrypted)

// Bcrypt password hashing
hashed, err := cryptox.HashPassword("mypassword")
err := cryptox.VerifyHashedPassword(hashed, "mypassword")

// SHA-256 hashing
base64Val, hashVal := cryptox.HashValue("some-token")
isValid := cryptox.IsHashValid(hashVal, base64Val)

// RSA key parsing (from base64-encoded PEM)
privKey, err := cryptox.GetRSAPrivateKey(base64PrivPEM)
pubKey, err := cryptox.GetRSAPublicKey(base64PubPEM)

// PKCE validation
valid := cryptox.IsPKCEValid(codeVerifier, codeChallenge, cryptox.PKCEMethodS256)
```

---

### `validator` — Request Validation

Uses [go-playground/validator](https://github.com/go-playground/validator) for struct validation.

```go
import "github.com/evenyosua18/ego/validator"

type CreateUserRequest struct {
    Name  string `validate:"required,min=2"`
    Email string `validate:"required,email"`
    Age   int    `validate:"gte=18"`
}

if err := validator.Validate(req); err != nil {
    return err // returns code.BadRequestError with field details
}
```

---

### `generator` — String Generation

```go
import "github.com/evenyosua18/ego/generator"

// Random alphanumeric string
str := generator.RandomString(32) // e.g. "aB3kZ9..."

// Cryptographically secure code (URL-safe base64)
code, err := generator.SecureCode(32) // e.g. "xK9m2..."
```

## Tech Stack

| Component            | Library                                                    |
|----------------------|------------------------------------------------------------|
| HTTP Router          | [Fiber v3](https://github.com/gofiber/fiber)               |
| Tracer               | [Sentry APM](https://github.com/getsentry/sentry-go)      |
| Config Loader        | [Viper](https://github.com/spf13/viper)                   |
| Database Scan Helper | [Scany](https://github.com/georgysavva/scany)              |
| Validation           | [Validator](https://github.com/go-playground/validator)    |
| Cache                | [go-redis](https://github.com/redis/go-redis)             |