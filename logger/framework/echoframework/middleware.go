// Package echoframework supports trace using Echo (https://echo.labstack.com/).
package echoframework

import (
	"github.com/labstack/echo/v4"

	"github.com/brbranch/go-gae-logger/logger/provider"
)

// Middleware create echo.MiddlewareFunc that start Trace of requests.
func Middleware(p provider.Provider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			newCtx, cf := p.StartSpan(c.Request(), c.Path())
			defer cf()

			newCtx = provider.Set(newCtx, p)
			c.SetRequest(c.Request().WithContext(newCtx))
			return next(c)
		}
	}
}
