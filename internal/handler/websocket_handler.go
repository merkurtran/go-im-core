package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/merkurtran/go-im-core/internal/service"
	ws "github.com/merkurtran/go-im-core/internal/websocket"
	"github.com/merkurtran/go-im-core/pkg/jwt"
	wslib "github.com/coder/websocket"
)

type WebSocketHandler struct {
	hub         *ws.Hub
	userService *service.UserService
	msgService  *service.MessageService
	msgRouter   ws.MessageRouter
	jwtSecret   string
	// allowedOrigins 用于 CORS 白名单校验，生产环境应配置具体域名
	allowedOrigins []string
}

var upgrader = wslib.AcceptOptions{
	OriginPatterns: []string{"localhost", "127.0.0.1"}, // 允许的 Origin 前缀
}

func NewWebSocketHandler(hub *ws.Hub, userSvc *service.UserService, msgSvc *service.MessageService, msgRouter ws.MessageRouter, jwtSecret string) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		userService: userSvc,
		msgService:  msgSvc,
		msgRouter:   msgRouter,
		jwtSecret:   jwtSecret,
	}
}

// HandleWebSocket 处理 WebSocket 升级请求
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// 1. 前置校验：Token 不能为空
	token := c.Query("token")
	if strings.TrimSpace(token) == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    2001,
			"message": "token is required",
		})
		return
	}

	// 2. JWT 认证
	claims, err := jwt.ParseToken(token, h.jwtSecret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    2003,
			"message": "invalid or expired token",
		})
		return
	}

	// 3. WebSocket 升级（必须在认证之后）
	conn, err := wslib.Accept(c.Writer, c.Request, &upgrader)
	if err != nil {
		return
	}
	defer conn.Close(wslib.StatusInternalError, "connection closed unexpectedly")

	// 4. 注册客户端到 Hub
	client := ws.NewClient(h.hub, conn, claims.UserID, h.msgRouter)
	h.hub.Register(client)
	defer h.hub.Unregister(client)

	// 5. 更新用户状态为在线
	if err := h.userService.UpdateUserStatus(c.Request.Context(), claims.UserID, "online"); err != nil {
		slog.Warn("failed to update user status to online", "user_id", claims.UserID, "error", err)
	}
	defer func() {
		if err := h.userService.UpdateUserStatus(c.Request.Context(), claims.UserID, "offline"); err != nil {
			slog.Warn("failed to update user status to offline", "user_id", claims.UserID, "error", err)
		}
	}()

	// 6. 推送未读消息数
	unreadCount, _ := h.msgService.GetUnreadCount(c.Request.Context(), claims.UserID)
	if unreadCount > 0 {
		h.hub.SendToUser(claims.UserID, []byte(fmt.Sprintf(`{"event":"unread","count":%d}`, unreadCount)))
	}

	// 7. 启动读写 Pump（阻塞直到连接关闭）
	go client.WritePump()
	client.ReadPump()
}
