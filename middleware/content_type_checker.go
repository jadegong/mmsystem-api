package middleware

import (
	"mime"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type ContentTypeCheckerConfig struct {
	Skipper middleware.Skipper
}

var (
	NoneJsonUrls map[string][]string = map[string][]string{
		"ALL": []string{"/upload", "/download", "/thumb"}, //Requests with file
	}
	DefaultContentTypeCheckerConfig ContentTypeCheckerConfig = ContentTypeCheckerConfig{
		Skipper: contentTypeCheckerSkipper,
	}
)

func contentTypeCheckerSkipper(c echo.Context) bool {
	ret := false
	if len(NoneJsonUrls["ALL"]) > 0 || len(NoneJsonUrls[c.Request().Method]) > 0 {
		urls := append(NoneJsonUrls["ALL"], NoneJsonUrls[c.Request().Method]...)
		for _, url := range urls {
			if strings.Contains(c.Request().URL.Path, url) {
				ret = true
				break
			}
		}
	}
	return ret
}

func ContentTypeCheckerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if DefaultContentTypeCheckerConfig.Skipper(c) {
				return next(c)
			}
			mediaType, params, _ := mime.ParseMediaType(c.Request().Header.Get("Content-Type"))
			charset, ok := params["charset"]
			if !ok {
				charset = "UTF-8"
			}
			if c.Request().ContentLength > 0 && !(mediaType == echo.MIMEApplicationJSON && strings.ToUpper(charset) == "UTF-8") {
				return echo.NewHTTPError(http.StatusUnsupportedMediaType, "Bad Content-Type or charset, expected 'application/json'")
			}
			return next(c)
		}
	}
}
