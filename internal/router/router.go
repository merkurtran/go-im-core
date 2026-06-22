package router

import (
	"github.com/gin-gonic/gin"
	"github.com/merkurtran/go-im-core/internal/handler"
	"github.com/merkurtran/go-im-core/internal/middlewares"
)

func Setup(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	messageHandler *handler.MessageHandler,
	wsHandler *handler.WebSocketHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	api := r.Group("/api/v1")

	ws := api.Group("/ws")
	{
		ws.GET("", wsHandler.HandleWebSocket)
	}

	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	authed := api.Group("")
	authed.Use(middlewares.AuthMiddleware(jwtSecret))

	users := authed.Group("/users")
	{
		users.GET("/me", userHandler.GetProfile)
		users.PUT("/me", userHandler.UpdateProfile)
	}

	messages := authed.Group("/messages")
	{
		messages.POST("", messageHandler.SendMessage)
		messages.GET("", messageHandler.GetConversation)
		messages.PATCH("/read", messageHandler.MarkRead)
	}

	return r
}
