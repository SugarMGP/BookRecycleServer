package ws

import (
	"encoding/json"
	"net/http"
	"sync"

	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/utils"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// CM 全局连接管理器
var CM *ConnectionManager

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

// Connection 用户连接
type Connection struct {
	Conn   *websocket.Conn
	UserID uint
	*sync.Mutex
}

// ConnectionManager 连接管理器
type ConnectionManager struct {
	connections     map[uint]*Connection
	offlineMessages map[uint][]*Message
	mutex           sync.RWMutex
	userChannel     chan []byte
}

// Message 消息结构体
type Message struct {
	Sender   uint   `json:"sender"`
	Receiver uint   `json:"receiver"`
	Content  string `json:"content"`
}

func (cm *ConnectionManager) handleMessage(message *Message) {
	cm.mutex.RLock()
	receiverConn, exists := cm.connections[message.Receiver]
	cm.mutex.RUnlock()

	if exists {
		// 如果接收者在线，直接发送消息
		if err := receiverConn.Conn.WriteJSON(*message); err != nil {
			zap.L().Warn("Error writing message", zap.Error(err))
		}
	} else {
		// 如果接收者不在线，存储消息
		cm.mutex.Lock()
		cm.offlineMessages[message.Receiver] = append(cm.offlineMessages[message.Receiver], message)
		cm.mutex.Unlock()
	}
}

func (cm *ConnectionManager) registerConnection(conn *websocket.Conn, uid uint) {
	cm.mutex.Lock()
	cm.connections[uid] = &Connection{
		Conn:   conn,
		UserID: uid,
		Mutex:  &sync.Mutex{},
	}
	// 检查是否有离线消息
	if messages, exists := cm.offlineMessages[uid]; exists {
		for _, msg := range messages {
			if err := conn.WriteJSON(*msg); err != nil {
				zap.L().Warn("Error sending offline message", zap.Error(err))
			}
		}
		// 清除离线消息队列
		delete(cm.offlineMessages, uid)
	}
	cm.mutex.Unlock()
}

func (cm *ConnectionManager) unregisterConnection(uid uint) {
	cm.mutex.Lock()
	delete(cm.connections, uid)
	cm.mutex.Unlock()
}

// Init 初始化连接管理器
func Init() {
	// 初始化连接管理器
	CM = &ConnectionManager{
		connections:     make(map[uint]*Connection),
		offlineMessages: make(map[uint][]*Message),
		userChannel:     make(chan []byte),
	}

	// 启动消息处理协程
	go func() {
		for {
			msg := <-CM.userChannel
			var message Message
			if err := json.Unmarshal(msg, &message); err != nil {
				zap.L().Warn("Error unmarshaling message", zap.Error(err))
				continue
			}
			CM.handleMessage(&message)
		}
	}()
}

// HandleWebSocket 处理 WebSocket 请求
func HandleWebSocket(c *gin.Context) {
	user, err := utils.GetUser(c)
	if err != nil {
		response.AbortWithException(c, apiException.ServerError, err)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		response.AbortWithException(c, apiException.WebSocketError, err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			zap.L().Warn("Error closing connection", zap.Error(err))
		}
	}(conn)

	CM.registerConnection(conn, user.ID)

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			CM.unregisterConnection(user.ID)
			break
		}

		if msgType == websocket.TextMessage {
			CM.userChannel <- msg
		}
	}
}
