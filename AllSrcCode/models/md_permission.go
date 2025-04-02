package models

type Permission struct {
	Id          int
	Name        string
	description string
}

func (m *Permission) TableName() string {
	return GetTableName("permission")
}
