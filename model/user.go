package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Email       string        `bson:"email" json:"email"`
	Password    string        `bson:"password" json:"password"`
	Type        uint8         `bson:"type" json:"type"` //用户类型，0: admin(管理manager和普通用户), 1: manager(数据录入人员), 2: user(普通用户)
	Name        string        `bson:"name" json:"name"`
	Mobile      string        `bson:"mobile" json:"mobile"`
	CreatedAt   time.Time     `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updatedAt"`   //最后更新时间
	IsActived   bool          `bson:"is_actived" json:"isActived"`   //是否激活
	ActivedAt   time.Time     `bson:"actived_at" json:"activedAt"`   //激活时间
	IsRejected  bool          `bson:"is_rejected" json:"isRejected"` //是否被拒绝
	RejectedAt  time.Time     `bson:"rejected_at" json:"rejectedAt"` //拒绝时间
	Avatar      bson.ObjectId `bson:"avatar,omitempty" json:"avatar"`
	LastLoginAt time.Time     `bson:"last_login_at" json:"lastLoginAt"` //上一次登录时间
	LastLoginIp string        `bson:"last_login_ip" json:"lastLoginIp"` //上一次登录IP
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
