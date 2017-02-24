package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"jadegong/api.mmsystem.com/g"
	"jadegong/api.mmsystem.com/model"
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

//用户注册
func Register(c echo.Context) error {
	now := time.Now()
	u := &model.User{
		Id:        bson.NewObjectId(),
		CreatedAt: now,
		UpdatedAt: now,
		IsActived: false,
	}
	//TODO bind func but encrypte password???
	if err := c.Bind(u); err != nil {
		return err
	}
	if u.Password == "" || strings.Contains("012", strconv.Itoa(int(u.Type))) == false {
		return c.JSON(http.StatusBadRequest, model.Error{Code: g.ERR_DATA_INVALID, Error: g.GetErrMsg(g.ERR_DATA_INVALID)})
	}

	if errNo := validateEmail(u.Email); errNo != g.SUCCESS {
		return c.JSON(g.GetErrHttpStatus(errNo), model.Error{Code: errNo, Error: g.GetErrMsg(errNo)})
	}

	session := g.Session()
	db := session.DB(g.Conf.DBName)
	defer session.Close()

	err := db.C(g.USER).Insert(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{Code: g.ERR_DB_FAILED, Error: g.GetErrMsg(g.ERR_DB_FAILED)})
	}
	return c.JSON(http.StatusCreated, u)
}

//Must use application/json
func CreateUser(c echo.Context) error {
	fmt.Printf("The header token: %s", c.Request().Header.Get("Authorization"))
	u := &model.User{}
	if err := c.Bind(u); err != nil {
		return err
	}
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

func validateEmail(email string) int {
	if g.IsEmail(email) == false {
		return g.ERR_REG_EMAIL_INVALID
	}

	session := g.Session()
	db := session.DB(g.Conf.DBName)
	defer session.Close()

	cnt, err := db.C(g.USER).Find(bson.M{"email": email}).Count()
	if err != nil {
		return g.ERR_DB_FAILED
	}
	if cnt > 0 {
		return g.ERR_REG_EMAIL_EXISTS
	}
	return g.SUCCESS
}
