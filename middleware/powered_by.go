package middleware

import (
	"github.com/labstack/echo"
)

func XPoweredByMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Add("X-Powered-By", "jadegong.com.cn")
			return next(c)
		}
	}
}
