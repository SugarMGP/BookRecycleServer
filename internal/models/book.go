package models

import "gorm.io/gorm"

// Book 书籍
type Book struct {
	gorm.Model
	UserID       uint
	Name         string // 书名
	Course       string // 适用课程
	Edition      string // 版次
	Publisher    string // 出版社
	Completeness string // 完好程度
	Img1         string // 封面图
	Img2         string
	Img3         string
	Price        string // 价格
	Note         string // 备注
}
