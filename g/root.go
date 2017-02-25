package g

import (
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"jadegong/api.mmsystem.com/model"
)

func initRoot() string {
	var errMsg string
	if IsEmail(Conf.RootEmail) == false {
		errMsg = "Invalid root user email format!"
		logrus.Error(errMsg)
		return errMsg
	}

	user := &model.User{}
	now := time.Now()

	session := Session()
	db := session.DB(Conf.DBName)
	defer session.Close()

	//add root user, find root user and user like this email count
	cnt, _ := db.C(USER).Find(bson.M{"$or": []bson.M{bson.M{"email": Conf.RootEmail}, bson.M{"type": 0}}}).Count()

	//add root user if no user found
	if cnt == 0 {
		user = &model.User{
			Id:        bson.NewObjectId(),
			Email:     Conf.RootEmail,
			Password:  EncryptPassword(Conf.RootPassword),
			Type:      0,
			CreatedAt: now,
			UpdatedAt: now,
			IsActived: true,
			ActivedAt: now,
		}
		err := db.C(USER).Insert(user)
		if err != nil {
			logrus.Errorf("Database error: %v", err.Error())
			return err.Error()
		}
		return ""
	} else {
		errMsg = "There already exists root user!"
		logrus.Error(errMsg)
		return errMsg
	}

	return ""
}
