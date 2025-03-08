package main

import (
	"GoSongs/config"
	"GoSongs/internal/handlers"
	"GoSongs/internal/repository"
	"GoSongs/internal/service"
	"GoSongs/migrations"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Songs API
// @version 1.0
// @description API для работы с песнями
// @host localhost:8080
// @BasePath /

// @Summary Получить информацию о песне
// @Description Возвращает текст, дату выхода и ссылку на песню
// @Tags songs
// @Accept  json
// @Produce  json
// @Param group query string true "Название группы"
// @Param song query string true "Название песни"
// @Success 200 {object} SongDetail
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /info [get]

func getSongInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"releaseDate": "16.07.2006",
		"text":        "Ooh baby, don't you know I suffer?\nOoh baby...",
		"link":        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	})
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DBDsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	repo := repository.NewRepository(db)
	svc := service.NewService(repo, cfg.MusicAPIURL)
	handler := handlers.NewHandler(svc)

	r := gin.Default()
	r.GET("/library", handler.GetSongs)
	r.GET("/song/:id", handler.GetSong)
	r.DELETE("/song/:id", handler.DeleteSong)
	r.PUT("/song/:id", handler.UpdateSong)
	r.POST("/song", handler.CreateSong)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/info", getSongInfo)

	log.Println("Starting server on :8080")
	r.Run(":8080")
}
