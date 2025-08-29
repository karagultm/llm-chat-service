package main

import (
	"myapp/internal/chat"
	"myapp/pkg/config"
	"myapp/pkg/database"

	"github.com/labstack/echo"
)

func main() {
	// confikg yükleme
	cfg := config.Load() // bu da özel fonksiyonmuş

	//logger açma

	//database
	db := database.Connect(cfg.DatabaseURL)
	db.AutoMigrate(&chat.ChatMessage{})
	//echo başlatma
	e := echo.New()

	chatRepo := chat.NewRepository(db)

	client := chat.NewClient(cfg.ApiKey)

	chatService := chat.NewService(chatRepo, client)

	chatHandler := chat.NewHandler(chatService)

	e.POST("v1/chat", chatHandler.Send)
	e.GET("v1/chat/:sessionId", chatHandler.ShowHistory)

	e.Logger.Fatal(e.Start(":8080"))
}
