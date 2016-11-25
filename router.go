package main

import (
	"dj/echo.demo/handler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func initRouter() *echo.Echo {
	e := echo.New()
	//todo middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	//Allow cross origin
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.GET("/stream", handler.GetStreamResponse) //Streaming response
	e.POST("/login", handler.AdminLogin)        //Admin login: form (name, email)
	e.GET("/users", handler.GetUsers)

	//User api group
	user := e.Group("/user")
	user.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &handler.JwtCustomClaims{},
		SigningKey: []byte("secret"),
	}))
	user.POST("", handler.CreateUser) //data-json: name, email
	user.GET("/:id", handler.GetUser)
	user.PUT("/:id", handler.UpdateUser)
	user.DELETE("/:id", handler.DeleteUser)
	user.POST("/avatar", handler.SaveAvatar) //With upload file
	return e
}
