package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"jadegong/echo.demo/model"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type JwtCustomClaims struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	jwt.StandardClaims
}

var (
	users = map[int]*model.User{}
	seq   = 1
)

func GetUsers(c echo.Context) error {
	return c.JSON(http.StatusCreated, users)
}

//Must use application/json
func CreateUser(c echo.Context) error {
	fmt.Printf("The header token: %s", c.Request().Header.Get("Authorization"))
	u := &model.User{
		Id: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	users[u.Id] = u
	seq += 1
	return c.JSON(http.StatusCreated, u)
}

func GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id >= seq {
		return fmt.Errorf("Invalid data id: %v", c.Param("id"))
	}
	u := users[id]
	return c.JSON(http.StatusCreated, u)
}

func UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	u := new(model.User)
	err = c.Bind(u)
	if err != nil || id >= seq {
		return fmt.Errorf("Invalid data id: %v", c.Param("id"))
	}
	users[id].Name = u.Name
	return c.JSON(http.StatusOK, users[id])
}

func DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id >= seq {
		return fmt.Errorf("Invalid data id: %v", c.Param("id"))
	}
	delete(users, id)
	return c.NoContent(http.StatusNoContent)
}

//upload file avatar
func SaveAvatar(c echo.Context) error {
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return err
	}
	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dist, err := os.Create(avatar.Filename)
	if err != nil {
		return err
	}
	defer dist.Close()

	if _, err = io.Copy(dist, src); err != nil {
		return err
	}

	ret := map[string]string{
		"Avatar": avatar.Filename,
	}
	return c.JSON(http.StatusOK, ret)
}

func AdminLogin(c echo.Context) error {
	username := c.FormValue("name")
	password := c.FormValue("password")
	if username == "admin" && password == "admin" {
		claims := &JwtCustomClaims{
			username,
			password,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		//generate encoded token
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
	return echo.ErrUnauthorized
}
