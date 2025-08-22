package main

import (
	"myapp/internal/chat"
	"myapp/pkg/config"
	"myapp/pkg/database"

	"github.com/labstack/echo"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
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
	client := openai.NewClient(option.WithAPIKey(cfg.ApiKey))
	chatService := chat.NewService(chatRepo, &client)
	chatHandler := chat.NewHandler(chatService)
	e.POST("v1/chat", chatHandler.SendMessage)
	e.GET("v1/chat/:sessionId", chatHandler.ShowHistory)

	e.Logger.Fatal(e.Start(":8080"))
}
