package model

type Bookmark struct {
	Base
	Description string `gorm:"column:description" json:"description"`
	Url         string `gorm:"column:url" json:"url"`
	CodeInt     int    `gorm:"column:code_int;autoIncrement" json:"code_int"`
	Code        string `gorm:"column:code" json:"code"`
	UserId      string `gorm:"column:user_id" json:"user_id"`
}
