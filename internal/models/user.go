package models

import "gorm.io/gorm"

// User 用户
type User struct {
	gorm.Model

	Username string // 用户名
	Password string // 密码
	Type     uint   // 用户类型 1：学生 2：收书员 3：管理员

	Name      string // 姓名
	StudentID string // 学号
	Phone     string // 联系电话
	Campus    uint   // 校区
	Address   string // 地址
}

// 校区编号
const (
	CampusZH uint = iota + 1
	CampusPF
	CampusMGS
)
