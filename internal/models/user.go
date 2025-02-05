package models

import "gorm.io/gorm"

// User 用户
type User struct {
	gorm.Model
	Username string `gorm:"username"`
	Password string `gorm:"password"`
}
