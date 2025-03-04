package models

import "time"

// Report 举报
type Report struct {
	ID           uint
	Reporter     uint   // 举报人
	ReporterName string // 举报人名字
	Seller       uint   // 卖家
	SellerName   string // 卖家名字
	Book         uint   // 举报的书
	BookName     string // 举报的书名
	Title        string // 举报标题
	Status       uint   // 状态 1未处理 2已撤销 3已处理
	CreatedAt    time.Time
}
