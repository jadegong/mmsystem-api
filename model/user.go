package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id            bson.ObjectId `bson:"_id"`
	Email         string        `bson:"email"`
	Password      string        `bson:"password"`
	Type          uint8         `bson:"type"` //用户类型，0: admin(管理manager和普通用户), 1: manager(数据录入人员), 2: user(普通用户)
	Name          string        `bson:"name"`
	Mobile        string        `bson:"mobile"`
	CreatedAt     time.Time     `bson:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at"` //最后更新时间
	IsActived     bool          `bson:"is_actived"` //是否激活
	ActivedAt     time.Time     `bson:"actived_at"` //激活时间
	Avatar        bson.ObjectId `bson:"avatar,omitempty"`
	LastLoginDate time.Time     `bson:"last_login_date"` //上一次登录时间
	LastLoginIp   string        `bson:"last_login_ip"`   //上一次登录IP
}

//用户日志
type UserLog struct {
	Id        bson.ObjectId `bson:"_id"`
	UserId    bson.ObjectId `bson:"user_id"`
	IpAddress bson.ObjectId `bson:"ip_address"`
	LogTime   time.Time     `bson:"log_time"`
	Type      uint8         `bson:"type"`
	Data      interface{}   `bson:"data,omitempty"`
}
