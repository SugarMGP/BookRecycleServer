package models

import "time"

// Recycle 旧书回收
type Recycle struct {
	ID         uint
	Seller     uint    // 卖家ID
	SellerName string  // 卖家姓名
	Img        string  // 图片
	Note       string  // 备注
	Weight     float64 // 预估重量
	Address    string  // 上门地址
	Campus     uint    // 校区
	Receiver   uint    // 接单者ID
	Status     uint    // 1等待接单 2已接单 3等待结算 4已完成
	CreatedAt  time.Time
}
