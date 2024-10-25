package common

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ydssx/kratos-kit/models"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"
)

// 定义WebSocket配置常量
const (
	maxMessageSize = 512                 // 允许的最大消息大小
	writeWait      = 10 * time.Second    // 写入消息的最大允许时间
	pongWait       = 60 * time.Second    // 等待下一个pong消息的最大时间
	pingPeriod     = (pongWait * 9) / 10 // 发送ping消息的周期，必须小于pongWait
)

var newline = []byte{'\n'}

type WsService struct {
	logger   *log.Helper
	upgrader websocket.Upgrader
	conns    sync.Map
	send     chan []byte
	stop     chan struct{}
}

func NewWsService(ctx context.Context, logger log.Logger) *WsService {
	ws := &WsService{
		logger: log.NewHelper(logger),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 注意：在生产环境中应该进行适当的源检查
			},
		},
		send: make(chan []byte, 256),
		stop: make(chan struct{}),
	}
	go ws.writePump()
	go ws.gracefulShutdown(ctx)

	return ws
}

func (s *WsService) AddConn(id string, conn *websocket.Conn) {
	s.conns.Store(id, conn)
}

func (s *WsService) RemoveConn(id string) {
	if conn, ok := s.conns.LoadAndDelete(id); ok {
		conn.(*websocket.Conn).Close()
	}
}

func (s *WsService) NotifyClients(messageType int, data []byte) {
	var wg sync.WaitGroup
	s.conns.Range(func(key, value interface{}) bool {
		wg.Add(1)
		go func(key interface{}, client *websocket.Conn) {
			defer wg.Done()
			if err := client.WriteMessage(messageType, data); err != nil {
				s.logger.Errorf("发送消息失败: %v", err)
				s.RemoveConn(key.(string))
			}
		}(key, value.(*websocket.Conn))
		return true
	})
	wg.Wait()
}

func (s *WsService) NotifyUser(userId string, messageType string, data interface{}) {
	if conn, ok := s.conns.Load(userId); ok {
		client := conn.(*websocket.Conn)

		// 构造消息
		message := struct {
			Type string      `json:"type"`
			Data interface{} `json:"data"`
		}{
			Type: messageType,
			Data: data,
		}

		// 发送消息
		if err := client.WriteJSON(message); err != nil {
			s.logger.Errorf("发送消息给用户 %s 失败: %v", userId, err)
			s.RemoveConn(userId)
		}
	} else {
		s.logger.Warnf("用户 %s 不在线，无法发送消息", userId)
	}
}

func (s *WsService) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	user, err := s.authenticateUser(r)
	if err != nil {
		s.logger.Errorf("用户认证失败: %v", err)
		http.Error(w, "未授权", http.StatusUnauthorized)
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Errorf("WebSocket连接失败: %v", err)
		return
	}

	// 使用用户ID而不是UUID
	userId := strconv.Itoa(int(user.ID))
	defer s.RemoveConn(userId)

	s.AddConn(userId, conn)
	s.configureConnection(conn)

	s.handleMessages(conn, userId)
}

func (s *WsService) authenticateUser(r *http.Request) (*models.User, error) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	return models.NewUserModel().SetUUIds(token).FirstOne()
}

func (s *WsService) configureConnection(conn *websocket.Conn) {
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})
}

func (s *WsService) handleMessages(conn *websocket.Conn, userUUID string) {
	defer s.logger.Infof("用户 %s 的WebSocket连接已关闭", userUUID)

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure, websocket.CloseUnsupportedData, websocket.CloseNoStatusReceived, websocket.CloseTLSHandshake) {
				s.logger.Errorf("WebSocket读取消息错误: %v", err)
			}
			break
		}
		s.logger.Infof("收到来自用户 %s 的消息: %s", userUUID, string(msg))
		s.NotifyClients(messageType, msg) // 使用原始消息类型
	}
}

func (ws *WsService) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-ws.send:
			if !ok {
				ws.logger.Info("WebSocket 发送通道已关闭")
				return
			}
			ws.broadcastMessage(message)
		case <-ticker.C:
			ws.NotifyClients(websocket.PingMessage, nil)
		case <-ws.stop:
			return
		}
	}
}

func (ws *WsService) broadcastMessage(message []byte) {
	ws.conns.Range(func(_, value interface{}) bool {
		client := value.(*websocket.Conn)
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			ws.logger.Errorf("广播消息失败: %v", err)
		}
		return true
	})
}

func (ws *WsService) gracefulShutdown(ctx context.Context) {
	<-ctx.Done()
	close(ws.stop)
	ws.conns.Range(func(_, value interface{}) bool {
		client := value.(*websocket.Conn)
		closeMsg := websocket.FormatCloseMessage(websocket.CloseGoingAway, "服务正在关闭")
		_ = client.WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(writeWait))
		_ = client.Close()
		return true
	})
	close(ws.send)
	ws.logger.Info("WebSocket 服务已完全关闭")
}
