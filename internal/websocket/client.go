package websocket

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/coder/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10 // 54s
	maxMessageSize = 512 << 10           // 512 KB
)

// WSMessage 定义客户端与服务端之间的 WebSocket 消息协议。
type WSMessage struct {
	Event   string          `json:"event"`   // message, ack, ping, status
	Payload json.RawMessage `json:"payload"` // 事件具体数据
}

// MessagePayload 聊天消息载荷
type MessagePayload struct {
	ReceiverID  string `json:"receiver_id"`
	Content     string `json:"content"`
	MessageType string `json:"message_type"` // text, image, file
}

// AckPayload 消息回执载荷
type AckPayload struct {
	MessageID   string `json:"message_id"`
	OtherUserID string `json:"other_user_id"`
	Status      string `json:"status"` // delivered, read
}

// MessageRouter 处理 WebSocket 消息的业务接口。
type MessageRouter interface {
	OnMessage(ctx context.Context, senderID string, msg *WSMessage) (*WSMessage, error)
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string
	router MessageRouter
}

func NewClient(hub *Hub, conn *websocket.Conn, userID string, router MessageRouter) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
		router: router,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close(websocket.StatusInternalError, "read pump exit")
	}()

	c.conn.SetReadLimit(maxMessageSize)

	readCtx, cancel := context.WithTimeout(context.Background(), pongWait)
	defer cancel()

	for {
		msgType, data, err := c.conn.Read(readCtx)
		if err != nil {
			break
		}

		if msgType == websocket.MessageText && len(data) > 0 {
			var msg WSMessage
			if err := json.Unmarshal(data, &msg); err != nil {
				slog.Warn("invalid ws message format", "user_id", c.userID, "error", err)
				continue
			}

			if c.router != nil {
				resp, err := c.router.OnMessage(readCtx, c.userID, &msg)
				if err != nil {
					slog.Error("message router error", "user_id", c.userID, "event", msg.Event, "error", err)
					continue
				}
				if resp != nil {
					respBytes, _ := json.Marshal(resp)
					c.send <- respBytes
				}
			}

			// 重置读超时（模拟心跳检测）
			cancel()
			readCtx, cancel = context.WithTimeout(context.Background(), pongWait)
		}
	}
	cancel()
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close(websocket.StatusNormalClosure, "write pump exit")
	}()

	for {
		select {
		case msg, ok := <-c.send:
			writeCtx, writeCancel := context.WithTimeout(context.Background(), writeWait)
			if !ok {
				c.conn.Close(websocket.StatusNormalClosure, "hub closed")
				writeCancel()
				return
			}

			if err := c.conn.Write(writeCtx, websocket.MessageText, msg); err != nil {
				writeCancel()
				return
			}
			writeCancel()

		case <-ticker.C:
			pingCtx, pingCancel := context.WithTimeout(context.Background(), writeWait)
			if err := c.conn.Ping(pingCtx); err != nil {
				pingCancel()
				return
			}
			pingCancel()
		}
	}
}

func (c *Client) UserID() string {
	return c.userID
}
