package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/merkurtran/go-im-core/internal/handler"
	"github.com/merkurtran/go-im-core/internal/middlewares"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api/v1")

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		users.PUT("/me/password", userHandler.ChangePassword)
		users.GET("/search", userHandler.SearchUsers)
	}

	messages := authed.Group("/messages")
	{
		messages.POST("", messageHandler.SendMessage)
		messages.GET("", messageHandler.GetConversation)
		messages.PATCH("/read", messageHandler.MarkRead)
		messages.PATCH("/:message_id/recall", messageHandler.RecallMessage)
	}

	return r
}
