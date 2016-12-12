package xPoweredByMiddleware

import "github.com/labstack/echo"

func XPoweredByMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("X-Powered-By", "jadegong.com.cn")
		return next(c)
	}
}
