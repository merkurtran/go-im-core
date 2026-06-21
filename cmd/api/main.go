package main

import (
	"context"
	"fmt"
	"log"
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
	"github.com/merkurtran/go-im-core/pkg/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client, err := database.NewMongoClient(&cfg.MongoDB)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	defer database.Disconnect(client)

	db := client.Database(cfg.MongoDB.Database)

	userRepo := mongo.NewUserRepo(db)
	messageRepo := mongo.NewMessageRepo(db)

	userSvc := service.NewUserService(userRepo, &cfg.JWT)
	msgSvc := service.NewMessageService(messageRepo, userRepo)

	authH := handler.NewAuthHandler(userSvc)
	userH := handler.NewUserHandler(userSvc)
	msgH := handler.NewMessageHandler(msgSvc)

	r := router.Setup(authH, userH, msgH, cfg.JWT.Secret)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.HTTP.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.HTTP.Timeout) * time.Second,
	}

	go func() {
		log.Printf("Server starting on: %d", cfg.HTTP.Port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}
