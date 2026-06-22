package websocket

import (
	"context"
	"time"

	"github.com/coder/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10 // 54s
	maxMessageSize = 512 << 10
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string
}

func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close(websocket.StatusInternalError, "read pump exit")
	}()

	c.conn.SetReadLimit(maxMessageSize)

	ctx, cancel := context.WithTimeout(context.Background(), pongWait)
	defer cancel()

	for {
		msgType, data, err := c.conn.Read(ctx)
		if err != nil {
			break
		}

		if msgType == websocket.MessageText && len(data) > 0 {
			cancel()
			ctx, cancel = context.WithTimeout(context.Background(), pongWait)

			// TODO: 将原始字节转发到业务层进行 JSON 解析和路由分发
			// 示例: c.hub.OnMessage(c.userID, data)
			_ = data // 当前仅消费数据防止阻塞，后续接入业务处理
		}
	}
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
