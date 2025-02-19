package models

import (
	"time"
)

// Message 消息结构体
type Message struct {
	ID        uint   `json:"id"`
	Sender    uint   `json:"sender"`
	Receiver  uint   `json:"receiver"`
	Content   string `json:"content"`
	CreatedAt time.Time
}
