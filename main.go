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

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

	log.Println("Starting server on :8080")
	r.Run(":8080")
}
