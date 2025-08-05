package user

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Email string `gorm:"unique;size:32"`
}

func (UserModel) TableName() string {
	return "users"
}
