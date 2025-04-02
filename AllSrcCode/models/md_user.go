package models

import "time"

type User struct {
	Id         int    `orm:"column(id);auto"`
	Username   string `orm:"column(username);size(50);unique;"`
	Password   string `orm:"column(password);size(255);"`
	Email      string `orm:"column(email);size(100);null"`
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
