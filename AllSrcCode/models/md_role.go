package models

type Role struct {
	Id          int
	Name        string
	description string
}

func (m *Role) TableName() string {
	return GetTableName("role")
}
