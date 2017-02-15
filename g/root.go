package g

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kless/osutil/user/crypt/sha512_crypt"
	"gopkg.in/mgo.v2/bson"
	"jadegong/api.mmsystem.com/model"
)

func initRoot() string {
	//if len(Conf.RootEmail) == 0 || len(Conf.RootPassword) == 0 {
	//logrus.Error("Please config root user email and password!")
	//return
	//}
	if IsEmail(Conf.RootEmail) == false {
		logrus.Error("Invalid root user email format!")
		return "Invalid root user email format!"
	}

	user := &model.User{}
	c := sha512_crypt.New()
	rootPassword, _ := c.Generate([]byte(Conf.RootPassword), []byte(""))
	now := time.Now()

	session := Session()
	db := session.DB(Conf.DBName)
	defer session.Close()
	//add root user, find root user and user like this email count
	users := []modle.User{}
	db.C(USER).Find(bson.M{"$or": []bson.M{bson.M{"email": Conf.RootEmail}, bson.M{"type": 0}}}).All(&users)

	//First time to add root user
	if len(users) == 0 || (len(users) == 1 && users[0].Email != Conf.RootEmail) {
		user = &model.User{
			Id:        bson.NewObjectId(),
			Email:     Conf.RootEmail,
			Password:  rootPassword,
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
		return nil
	} else if len(users) == 1 && users[0] == Conf.RootEmail { //TODO No root user, but has normal user
	}
	for i, u := range users {
		if u.Type == 0 && u.Email == Conf.RootEmail {
			return nil
		}
	}
	if len(users) == 0 {
		return nil
	}
	err := db.C(USER).Update(bson.M{"type": 0}, bson.M{"$set": bson.M{"type": 2}})
	if err != nil {
		logrus.Errorf("Database error: %v", err.Error())
		return err.Error()
	}

	return nil
}
