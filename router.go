package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"jadegong/api.mmsystem.com/handler"
	customMiddleware "jadegong/api.mmsystem.com/middleware"
)

func initRouter() *echo.Echo {
	e := echo.New()
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(customMiddleware.XPoweredByMiddleware())

	//Allow cross origin
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderAcceptEncoding, echo.HeaderAuthorization, echo.HeaderContentType, echo.HeaderOrigin},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	//Content-Type: application/json, charset: UTF-8
	e.Use(customMiddleware.ContentTypeCheckerMiddleware())

	//token authorization
	e.Use(customMiddleware.TokenAuthWithConfig(customMiddleware.TokenAuthConfig{
		Skipper: customMiddleware.TokenAuthSkipper,
	}))

	e.POST("/register", handler.Register) //data-json: name, email, type
	e.POST("/login", handler.Login)

	//User api group
	user := e.Group("/user")
	user.GET("", handler.GetUser)
	user.GET("/notes", handler.GetRegisterNotes)
	user.PUT("/verify/:id", handler.VerifyUser)
	user.POST("/avatar", handler.SaveAvatar) //With upload file
	return e
}
