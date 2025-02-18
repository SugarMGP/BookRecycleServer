package ws

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"bookrecycle-server/internal/apiException"
	"bookrecycle-server/internal/models"
	"bookrecycle-server/internal/services/messageService"
	"bookrecycle-server/internal/utils"
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

func (cm *ConnectionManager) handleMessage(message *models.Message) {
	// 保存消息到数据库
	if err := messageService.SaveMessage(*message); err != nil {
		zap.L().Warn("Error saving message to database", zap.Error(err))
		return
	}

	cm.mutex.RLock()
	receiverConn, exists := cm.connections[message.Receiver]
	senderConn, senderExists := cm.connections[message.Sender]
	cm.mutex.RUnlock()

	// 发送给接收者
	if exists {
		if err := receiverConn.Conn.WriteJSON(*message); err != nil {
			zap.L().Warn("Error writing message", zap.Error(err))
		}
	}

	// 发送给自己
	if senderExists {
		if err := senderConn.Conn.WriteJSON(*message); err != nil {
			zap.L().Warn("Error writing message", zap.Error(err))
		}
	}
}

func (cm *ConnectionManager) registerConnection(conn *websocket.Conn, uid uint) {
	cm.mutex.Lock()
	cm.connections[uid] = &Connection{
		Conn:   conn,
		UserID: uid,
		Mutex:  &sync.Mutex{},
	}

	messages, err := messageService.GetMessagesByUser(uid)
	if err != nil {
		zap.L().Warn("Error getting messages", zap.Error(err))
	}
	for _, msg := range messages {
		if err := conn.WriteJSON(msg); err != nil {
			zap.L().Warn("Error sending history message", zap.Error(err))
		}
	}
	cm.mutex.Unlock()
}

func (cm *ConnectionManager) unregisterConnection(uid uint) {
	cm.mutex.Lock()
	delete(cm.connections, uid)
	cm.mutex.Unlock()
}

// Init 初始化
func Init() {
	// 初始化连接管理器
	cm = &ConnectionManager{
		connections: make(map[uint]*Connection),
		userChannel: make(chan []byte),
		stop:        false,
	}

	// 启动消息处理协程
	go func() {
		for !cm.stop {
			msg := <-cm.userChannel
			var message models.Message
			if err := json.Unmarshal(msg, &message); err != nil {
				zap.L().Warn("Error unmarshaling message", zap.Error(err))
				continue
			}
			cm.handleMessage(&message)
		}
	}()
}

// Stop 停止消息处理
func Stop() {
	cm.stop = true
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

	cm.registerConnection(conn, user.ID)

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			cm.unregisterConnection(user.ID)
			break
		}

		if msgType == websocket.TextMessage {
			var message models.Message
			if err := json.Unmarshal(msg, &message); err != nil {
				zap.L().Warn("Error unmarshaling message", zap.Error(err))
				continue
			}

			// 填充信息
			message.Sender = user.ID
			message.CreatedAt = time.Now()

			msg, err = json.Marshal(message)
			if err != nil {
				zap.L().Warn("Error marshaling message", zap.Error(err))
				continue
			}
			cm.userChannel <- msg
		}
	}
}
