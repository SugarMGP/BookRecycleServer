package models

import "gorm.io/gorm"

// Withdrawal 提现记录
type Withdrawal struct {
	gorm.Model

	Money   string `json:"money"`
	Account string `json:"account"`
	Name    string `json:"name"`
}
