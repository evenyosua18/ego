---
trigger: always_on
---

## Project Context
Your name is **ego framework**.

You are operating within the `ego` project. 
**Project Purpose**: `ego` is a lightweight Go core framework designed for microservices and backend systems. It provides a structured foundation to rapidly build scalable, maintainable, and observable applications with minimal boilerplate.

## Framework Components
When implementing new microservices or adding features, you **MUST** prioritize using these internal packages over external libraries or standard library counterparts:
- **`config`**: Viper-powered configuration management (`github.com/evenyosua18/ego/config`).
- **`http`**: Fiber-based REST server & router with built-in panic recover, logger, rate limiter, CORS, and auth middleware (`github.com/evenyosua18/ego/http`).
- **`sqldb`**: MySQL wrapper using `go-sql-driver/mysql` and `scany` for querying, along with robust transaction management (`github.com/evenyosua18/ego/sqldb`). Always use `sqldb.GetDbManager()`.
- **`cache`**: Redis-backed cache manager with `singleflight` deduplication (`github.com/evenyosua18/ego/cache`).
- **`auth`**: Service-to-service token management (`github.com/evenyosua18/ego/auth`).
- **`logger`**: Structured leveled logging (`github.com/evenyosua18/ego/logger`). Prefer `logger.Info()`, `logger.Error()`, etc., over raw `fmt` or `log`.
- **`tracer`**: Distributed tracing via Sentry APM (`github.com/evenyosua18/ego/tracer`). Start spans with `tracer.StartSpan`.
- **`request`**: Configurable HTTP Client wrapper with retries and token injection (`github.com/evenyosua18/ego/request`).
- **`code`**: Standardized application error codes mapped to HTTP/gRPC statuses (`github.com/evenyosua18/ego/code`).
- **`cryptox`**: Cryptography utilities (AES, RSA, bcrypt, SHA) (`github.com/evenyosua18/ego/cryptox`).
- **`validator`**: Struct validation using `go-playground/validator` (`github.com/evenyosua18/ego/validator`).
- **`generator`**: Secure string and ID generation (`github.com/evenyosua18/ego/generator`).

## Tech Stack
- **Go**: 1.25.0
- **Web Router**: Fiber v3
- **APM**: Sentry
- **Configuration**: Viper
- **Database Mapping**: Scany
- **Validation**: go-playground/validator
- **Cache**: go-redis/v9

## Architectural Rules
1. **Never** use raw `net/http` or standard library logging (`log`) if `ego/http` or `ego/logger` provides the equivalent functionality.
2. **Routing**: Structure new API endpoints by registering them via groups using `http.RegisterRouteByGroup`.
3. **Error Handling**: Use structured error handling. Wrap all errors via `code.Wrap()` using predefined codes rather than raw `fmt.Errorf` or `errors.New`.
4. **Transactions**: Use `mgr.WrappedBeginTrx` from `ego/sqldb` for executing multiple queries within a transaction block.