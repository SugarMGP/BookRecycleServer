package ws

import (
	"encoding/json"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/models"
	"bookrecycle-server/internal/services/messageService"
	"bookrecycle-server/internal/services/userService"
	"bookrecycle-server/internal/utils/jwt"
	"bookrecycle-server/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var cm *ConnectionManager

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

// ConnectionManager 连接管理器
type ConnectionManager struct {
	connections map[uint]*websocket.Conn
	mutex       sync.RWMutex
	userChannel chan models.Message
	stop        atomic.Bool
}

func (cm *ConnectionManager) handleMessage(message *models.Message) {
	// 保存消息到数据库
	if err := messageService.CreateMessage(message); err != nil {
		zap.L().Warn("Error saving message to database", zap.Error(err))
		return
	}

	cm.mutex.RLock()
	receiverConn, exists := cm.connections[message.Receiver]
	senderConn, senderExists := cm.connections[message.Sender]
	cm.mutex.RUnlock()

	// 发送给接收者
	if exists {
		if err := receiverConn.WriteJSON(message); err != nil {
			zap.L().Warn("Error writing message", zap.Error(err))
			_ = receiverConn.Close()
			cm.unregisterConnection(message.Receiver)
		}
	}

	// 发送给自己
	if senderExists {
		if err := senderConn.WriteJSON(message); err != nil {
			zap.L().Warn("Error writing message", zap.Error(err))
			_ = senderConn.Close()
			cm.unregisterConnection(message.Sender)
		}
	}
}

func (cm *ConnectionManager) registerConnection(conn *websocket.Conn, uid uint) {
	// 获取历史消息
	messages, err := messageService.GetMessagesByUser(uid)
	if err != nil {
		zap.L().Warn("Error getting messages", zap.Error(err))
	}

	// 注册连接并推送历史消息
	cm.mutex.Lock()
	cm.connections[uid] = conn
	cm.mutex.Unlock()
	for _, msg := range messages {
		if err := conn.WriteJSON(msg); err != nil {
			zap.L().Warn("Error sending history message", zap.Error(err))
		}
	}
}

func (cm *ConnectionManager) unregisterConnection(uid uint) {
	// 删除连接
	cm.mutex.Lock()
	delete(cm.connections, uid)
	cm.mutex.Unlock()
}

// Init 初始化
func Init() {
	// 初始化连接管理器
	cm = &ConnectionManager{
		connections: make(map[uint]*websocket.Conn),
		userChannel: make(chan models.Message),
		stop:        atomic.Bool{},
	}

	// 启动消息处理协程
	go func() {
		for !cm.stop.Load() {
			msg := <-cm.userChannel
			cm.handleMessage(&msg)
		}
	}()
}

// Stop 停止消息处理
func Stop() {
	cm.stop.Store(true)
	cm.mutex.Lock()
	for uid, conn := range cm.connections {
		_ = conn.Close()
		delete(cm.connections, uid)
	}
	cm.mutex.Unlock()
}

// HandleWebSocket 处理 WebSocket 请求
func HandleWebSocket(c *gin.Context) {
	token := c.Query("token")
	claims, err := jwt.ParseToken(token)
	if err != nil {
		response.AbortWithException(c, apiException.NoAccessPermission, err)
		return
	}

	uid := claims.UserID
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

	cm.registerConnection(conn, uid)
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			cm.unregisterConnection(uid)
			break
		}

		if msgType == websocket.TextMessage {
			var message models.Message
			if err := json.Unmarshal(msg, &message); err != nil {
				zap.L().Warn("Error unmarshaling message", zap.Error(err))
				continue
			}

			// 填充信息
			message.Time = time.Now().Format(time.DateTime)
			message.Sender = uid
			if message.Sender == message.Receiver {
				continue
			}

			// 填充用户信息
			sender, err := userService.GetUserByID(message.Sender)
			if err != nil {
				response.AbortWithException(c, apiException.ServerError, err)
				return
			}
			message.SenderName = sender.Name

			receiver, err := userService.GetUserByID(message.Receiver)
			if err != nil {
				response.AbortWithException(c, apiException.ServerError, err)
				return
			}
			message.ReceiverName = receiver.Name

			cm.userChannel <- message
		}
	}
}
