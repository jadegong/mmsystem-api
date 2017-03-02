package handler

import (
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
	"jadegong/api.mmsystem.com/middleware"
	"jadegong/api.mmsystem.com/model"
)

type JwtCustomClaims struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	jwt.StandardClaims
}

//用户注册
func Register(c echo.Context) error {
	now := time.Now()
	u := &model.User{
		Id:         bson.NewObjectId(),
		CreatedAt:  now,
		UpdatedAt:  now,
		IsActived:  false,
		IsRejected: false,
	}
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Code: g.ERR_DATA_INVALID, Error: g.GetErrMsg(g.ERR_DATA_INVALID)})
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

	//encrypt password
	u.Password = g.EncryptPassword(u.Password)
	err := db.C(g.USER).Insert(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{Code: g.ERR_DB_FAILED, Error: g.GetErrMsg(g.ERR_DB_FAILED)})
	}
	return c.JSON(http.StatusCreated, u)
}

func Login(c echo.Context) error {
	var (
		u    = &model.User{}
		user = &model.User{}
		ret  map[string]interface{}
	)
	err := c.Bind(u)
	if err != nil || u.Email == "" || g.IsEmail(u.Email) == false {
		return c.JSON(http.StatusBadRequest, model.Error{Code: g.ERR_DATA_INVALID, Error: g.GetErrMsg(g.ERR_DATA_INVALID)})
	}
	session := g.Session()
	db := session.DB(g.Conf.DBName)
	defer session.Close()

	err = db.C(g.USER).Find(bson.M{"email": u.Email}).One(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{Code: g.ERR_DB_FAILED, Error: g.GetErrMsg(g.ERR_DB_FAILED)})
	}
	if user.IsActived != true {
		return c.JSON(http.StatusForbidden, model.Error{Code: g.ERR_USER_NOT_VERIFIED, Error: g.GetErrMsg(g.ERR_USER_NOT_VERIFIED)})
	}
	if g.EncryptPassword(u.Password) != user.Password {
		return c.JSON(http.StatusBadRequest, model.Error{Code: g.ERR_PASSWORD_INVALID, Error: g.GetErrMsg(g.ERR_PASSWORD_INVALID)})
	}

	//generate encoded token
	token, _ := middleware.New()
	if g.Conf.Cache == g.CACHE_REDIS {
		err := g.Redis.Set(g.AUTH_TOKEN_NS+middleware.Hash(token), user.Id.Hex(), g.Conf.TokenDuration).Err()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.Error{Code: g.ERR_REDIS_NOT_AVAILABLE, Error: g.GetErrMsg(g.ERR_REDIS_NOT_AVAILABLE)})
		}
	} else {
		g.Cache.Set(g.AUTH_TOKEN_NS+middleware.Hash(token), user.Id.Hex(), g.Conf.TokenDuration)
	}

	ipAddress := c.RealIP()
	err = db.C(g.USER).UpdateId(user.Id, bson.M{"$set": bson.M{"last_login_ip": ipAddress,
		"$currentDate": bson.M{"last_login_at": true}}})

	userMap := getUserMap(user)
	ret["user"] = userMap
	ret["token"] = token

	return c.JSON(http.StatusOK, ret)
}

// Get one user's info
func GetUser(c echo.Context) error {
	user := model.User{}
	userID := c.Get("remote_user")

	session := g.Session()
	db := session.DB(g.Conf.DBName)
	defer session.Close()

	err := db.C(g.USER).FindId(bson.ObjectIdHex(userID)).One(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{Code: g.ERR_DB_FAILED, Error: g.GetErrMsg(g.ERR_DB_FAILED)})
	}
	return c.JSON(http.StatusCreated, getUserMap(user))
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
		//generate encoded token
		token, _ := middleware.New()
		if g.Conf.Cache == g.CACHE_REDIS {
			err := g.Redis.Set(g.AUTH_TOKEN_NS+middleware.Hash(token), username, g.Conf.TokenDuration).Err()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, model.Error{Code: g.ERR_REDIS_NOT_AVAILABLE, Error: g.GetErrMsg(g.ERR_REDIS_NOT_AVAILABLE)})
			}
		} else {
			g.Cache.Set(g.AUTH_TOKEN_NS+middleware.Hash(token), username, g.Conf.TokenDuration)
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": token,
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

func getUserMap(user *model.User) map[string]interface{} {
	userName := user.Name
	if userName == "" {
		userName = user.Email
	}

	ret := map[string]interface{}{
		"id":          user.Id.Hex(),
		"email":       user.Email,
		"type":        user.Type,
		"name":        user.Name,
		"mobile":      user.Mobile,
		"createdAt":   user.CreatedAt,
		"updatedAt":   user.UpdatedAt,
		"isActived":   user.IsActived,
		"activedAt":   user.ActivedAt,
		"isRejected":  user.IsRejected,
		"rejectedAt":  user.RejectedAt,
		"avatar":      user.Avatar.Hex(),
		"lastLoginAt": user.LastLoginAt,
		"lastLoginIp": user.LastLoginIp,
	}

	return ret
}
