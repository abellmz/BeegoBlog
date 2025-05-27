package models

import "time"

type User struct {
	Id         int64  `orm:"column(id);pk;type(bigint);" json:"id" bson:"_id"`
	Username   string `orm:"column(username);size(50);unique;" json:"username" bson:"username"`
	Password   string `orm:"column(password);size(255);"`
	Email      string `orm:"column(email);size(100);null" json:"email" bson:"email"`
	Age        int    `orm:"column(age);"`
	LoginCount int
	LastTime   time.Time
	LastIp     string
	State      int8 `orm:"column(state);"`
	Created    time.Time
	Updated    time.Time
}

func (m *User) TableName() string {
	return GetTableName("user")
}
