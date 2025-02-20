package models

import (
	"time"
)

// Message 消息结构体
type Message struct {
	ID           uint      `json:"id"`
	Sender       uint      `json:"sender"`
	SenderName   string    `json:"sender_name"`
	Receiver     uint      `json:"receiver"`
	ReceiverName string    `json:"receiver_name"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"time"`
}
