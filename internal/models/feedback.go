package models

import (
	"time"
)

// Feedback 意见反馈
type Feedback struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"time"`
}
