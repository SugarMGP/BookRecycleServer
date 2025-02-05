package models

import "gorm.io/gorm"

// Student 学生用户
type Student struct {
	gorm.Model
	Username  string // 学号
	Password  string // 密码
	Name      string // 姓名
	Phone     string // 联系电话
	Campus    uint   // 校区
	Location  string // 宿舍地址
	HasLogged bool   // 之前是否登录过
}

// 校区编号
const (
	CampusZH uint = iota + 1
	CampusPF
	CampusMGS
)
