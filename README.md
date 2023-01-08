# go-gae-logger

Simple Logger for Google Cloud Logging that supports structured logging and Google Cloud Trace.  
Currently, This package supports following service and framework:

* Trace service
  * [OpenTelemetry](https://opentelemetry.io/)
* Framework
  * [Echo](https://echo.labstack.com/)
  
# Supporting Structure

* severity
* message
* sourceLocation
* spanId
* trace

# Usage
## Using Echo
```go
import (
	"net/http"
	
	"github.com/brbranch/go-gae-logger/logger"
	"github.com/brbranch/go-gae-logger/logger/framework/echoframework"
	"github.com/brbranch/go-gae-logger/logger/provider/otlm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
    e := a.echo

    provider := otlm.Create("[ProjectName]/trace", os.Getenv("GOOGLE_CLOUD_PROJECT"))
    if /* expression to check local server */ {
        provider.LocalMode()
    }

    e.Use(echoframework.Middleware(provider))
    e.GET("/", func(c echo.Context) error {
        ctx , cf := logger.Span(c.Request().Context(), "index")
        defer cf()
		
		logger.Debug(ctx, "hello %s!", "logging")
		
        return c.String(http.StatusOK, "ok")
    })

    log.Fatal(e.Start(fmt.Sprintf(":%s", "8080")))
}
```

# License
MIT

# References
* [Structured Logging](https://cloud.google.com/logging/docs/structured-logging)
* [Google Cloud Trace](https://cloud.google.com/trace)
