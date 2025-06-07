# EGO
**Ego** is a lightweight Go core framework designed for microservices and backend systems. It provides a structured foundation to rapidly build scalable, maintainable, and observable applications with minimal boilerplate.

## Features
- ⚙️ **Configuration Management** 
- 📦 **Mysql Database Integration**
- 🧠 **Tracing Implementation**
- 🌐 **REST Server**
- 📄 **Standardized API Responses**

## Getting Started
Installation
```
go get github.com/your-org/ego
```

Run REST Server
```
import (
	"github.com/evenyosua18/ego/app"
	"github.com/evenyosua18/ego/http"
)

func main() {
	http.RegisterRouteByGroup("public", []http.RouteFunc{
		// Register your public routes here
	})

	app.GetApp().RunRest()
}
```

## Tech Stack
| Name                 | Reference                                              |
|----------------------|--------------------------------------------------------|
| Routing              | [Fiber](https://github.com/gofiber/fiber)              |
| Tracer               | [Sentry (APM)](https://github.com/getsentry/sentry-go) |
| Config Loader        | [Viper](https://github.com/spf13/viper)                |
| Database Scan Helper | [Scanny](https://github.com/georgysavva/scany)         |

### ONGOING
- [ ] Validator 
- [ ] Benchmark
- [ ] Unit Test
- [ ] Template Project Structure
- [ ] Automation API Documentation