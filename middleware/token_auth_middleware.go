package middleware

import (
	"errors"
	"strings"

	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"jadegong/api.mmsystem.com/g"
	"net/http"
)

type (
	TokenAuthConfig struct {
		Skipper middleware.Skipper

		// Context key to store user information from the token into context.
		// Optional. Default value "remote_user".
		ContextKey string

		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup string

		// AuthScheme to be used in the Authorization header.
		// Optional. Default value "Token".
		AuthScheme string
	}
	jwtExtractor func(echo.Context) (string, error)
)

var (
	tokenEntropy = 32

	DefaultTokenAuthConfig = TokenAuthConfig{
		Skipper:     defaultSkipper,
		ContextKey:  "remote_user",
		TokenLookup: "header:" + echo.HeaderAuthorization,
		AuthScheme:  "Token",
	}
)

// defaultSkipper returns false which processes the middleware.
func defaultSkipper(c echo.Context) bool {
	return false
}

// Default middleware
func TokenAuthMiddleware() echo.MiddlewareFunc {
	c := DefaultTokenAuthConfig
	return TokenAuthWithConfig(c)
}

// Add middleware of token with config
func TokenAuthWithConfig(config TokenAuthConfig) echo.MiddlewareFunc {
	//Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultTokenAuthConfig.Skipper
	}
	if config.ContextKey == "" {
		config.ContextKey = DefaultTokenAuthConfig.ContextKey
	}
	if config.TokenLookup == "" {
		config.TokenLookup = DefaultTokenAuthConfig.TokenLookup
	}
	if config.AuthScheme == "" {
		config.AuthScheme = DefaultTokenAuthConfig.AuthScheme
	}

	//Initialize
	parts := strings.Split(config.TokenLookup, ":")
	extractor := tokenFromHeader(parts[1], config.AuthScheme)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			token, err := extractor(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			userID := TokenIdentify(token)
			if userID != "" {
				c.Set(config.ContextKey, userID)
				return next(c)
			}
			return echo.ErrUnauthorized
		}
	}
}

// Generate a random token
func New() (string, error) {
	bytes := make([]byte, tokenEntropy)
	_, err := rand.Read(bytes[:cap(bytes)])
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// Hash applies a simple MD5 hash over a token, making it safe to store
func Hash(token string) string {
	hashed := md5.Sum([]byte(token))
	return base64.URLEncoding.EncodeToString(hashed[:])
}

// Read from storage -- redis or memory
func TokenIdentify(token string) string {
	var userID string
	if g.Conf.Cache == g.CACHE_REDIS {
		userID, _ = g.Redis.Get(g.AUTH_TOKEN_NS + Hash(token)).Result()
	} else {
		ret, _ := g.Cache.Get(g.AUTH_TOKEN_NS + Hash(token))
		userID, _ = ret.(string)
	}
	return userID
}

// jwtFromHeader returns a `jwtExtractor` that extracts token from the request header.
func tokenFromHeader(header string, authScheme string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", errors.New("Missing or invalid jwt in the request header")
	}
}
