package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/merkurtran/go-im-core/internal/config"
	"github.com/merkurtran/go-im-core/internal/handler"
	"github.com/merkurtran/go-im-core/internal/repository/mongo"
	"github.com/merkurtran/go-im-core/internal/router"
	"github.com/merkurtran/go-im-core/internal/service"
	"github.com/merkurtran/go-im-core/internal/websocket"
	"github.com/merkurtran/go-im-core/pkg/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	setupLogger(cfg)

	client, err := database.NewMongoClient(&cfg.MongoDB)
	if err != nil {
		slog.Error("failed to connect to mongo", "error", err)
		os.Exit(1)
	}
	defer database.Disconnect(client)

	db := client.Database(cfg.MongoDB.Database)

	userRepo := mongo.NewUserRepo(db)
	messageRepo := mongo.NewMessageRepo(db)

	tokenGen := service.NewJWTTokenGenerator(cfg.JWT.Secret, cfg.JWT.Expire)
	userSvc := service.NewUserService(userRepo, tokenGen)
	msgSvc := service.NewMessageService(messageRepo, userRepo)

	wsHub := websocket.NewHub()
	go wsHub.Run()

	msgRouter := websocket.NewDefaultMessageRouter(msgSvc, userSvc, wsHub)
	wsH := handler.NewWebSocketHandler(wsHub, userSvc, msgSvc, msgRouter, cfg.JWT.Secret)

	authH := handler.NewAuthHandler(userSvc)
	userH := handler.NewUserHandler(userSvc)
	msgH := handler.NewMessageHandler(msgSvc)

	r := router.Setup(authH, userH, msgH, wsH, cfg.JWT.Secret)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.HTTP.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.HTTP.Timeout) * time.Second,
	}

	go func() {
		slog.Info("server starting", "port", cfg.HTTP.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	slog.Info("server shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}
	slog.Info("server exited")
}

func setupLogger(cfg *config.Config) {
	level := slog.LevelInfo
	switch cfg.Log.Level {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	opts := &slog.HandlerOptions{Level: level}
	handler := slog.NewTextHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
}
