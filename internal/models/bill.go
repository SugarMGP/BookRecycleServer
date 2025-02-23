package models

import (
	"time"
)

// Bill 账单
type Bill struct {
	ID        uint
	User      uint
	Amount    string
	CreatedAt time.Time
}
