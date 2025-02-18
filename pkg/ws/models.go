package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Connection 用户连接
type Connection struct {
	Conn   *websocket.Conn
	UserID uint
	*sync.Mutex
}

// ConnectionManager 连接管理器
type ConnectionManager struct {
	connections map[uint]*Connection
	mutex       sync.RWMutex
	userChannel chan []byte
	stop        bool
}
