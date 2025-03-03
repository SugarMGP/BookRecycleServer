package models

import "gorm.io/gorm"

type Withdrawal struct {
	gorm.Model

	Money   string `json:"money"`
	Account string `json:"account"`
	Name    string `json:"name"`
}
