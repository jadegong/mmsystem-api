package handler

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
	"jadegong/api.mmsystem.com/g"
	"jadegong/api.mmsystem.com/model"
	"net/http"
)

//获取申请注册信息
//参数传递：type: 1,2
func GetRegisterNotes(c echo.Context) error {
	var (
		userList = []model.User{}
		ret      = []map[string]string{}
		data     = map[string]int{}
	)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Code: g.ERR_DATA_INVALID, Error: g.GetErrMsg(g.ERR_DATA_INVALID)})
	}
	if data["type"] != 1 && data["type"] != 2 {
		return c.JSON(http.StatusBadRequest, model.Error{Code: g.ERR_DATA_INVALID, Error: g.GetErrMsg(g.ERR_DATA_INVALID)})
	}

	session := g.Session()
	db := session.DB(g.Conf.DBName)
	defer session.Close()

	err := db.C(g.USER).Find(bson.M{"type": data["type"], "is_actived": false, "is_rejected": false}).All(&userList)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{Code: g.ERR_DB_FAILED, Error: g.GetErrMsg(g.ERR_DB_FAILED)})
	}
	for index := 0; index < len(userList); index++ {
		ret = append(ret, getUserMap(&userList[index]))
	}
	return c.JSON(http.StatusOK, ret)
}

//通过审核
func VerifyUser(c echo.Context) error {
	var (
		userID = c.Param("id")
		user   = &model.User{}
	)
	if bson.IsObjectIdHex(userID) == false {
		return c.JSON(http.StatusBadRequest, model.Error{Code: g.ERR_DATA_INVALID, Error: g.GetErrMsg(g.ERR_DATA_INVALID)})
	}

	session := g.Session()
	db := session.DB(g.Conf.DBName)
	defer session.Close()

	err := db.C(g.USER).FindId(bson.ObjectIdHex(userID)).One(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{Code: g.ERR_DB_FAILED, Error: g.GetErrMsg(g.ERR_DB_FAILED)})
	}
	if user.IsActived || user.IsRejected { //已激活或已拒绝
		return c.JSON(http.StatusForbidden, model.Error{Code: g.ERR_PERMISSION_DENIED, Error: g.GetErrMsg(g.ERR_PERMISSION_DENIED)})
	}
	err = db.C(g.USER).UpdateId(bson.ObjectIdHex(userID), bson.M{"$set": bson.M{"is_actived": true,
		"$currentDate": bson.M{"actived_at": true, "updated_at": true}}})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{Code: g.ERR_DB_FAILED, Error: g.GetErrMsg(g.ERR_DB_FAILED)})
	}

	return nil
}
