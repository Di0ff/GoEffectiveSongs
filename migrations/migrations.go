package migrations

import (
	"GoSongs/internal/models"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Songs{}, &models.Verses{}); err != nil {
		log.Fatalf("Error migrating: %v", err)
		return err
	}

	log.Println("Database migration succeeded")
	return nil
}
