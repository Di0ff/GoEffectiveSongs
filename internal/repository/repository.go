package repository

import (
	"GoSongs/internal/models"
	"log"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	conn *gorm.DB
}

func NewRepository(conn *gorm.DB) *Repository {
	return &Repository{conn: conn}
}

func (r *Repository) GetSongs(filters map[string]interface{}, page, pageSize int) ([]models.Songs, error) {
	var songs []models.Songs
	query := r.conn.Model(&models.Songs{})

	for key, value := range filters {
		if value != "" {
			if key == "song" {
				query = query.Where("song LIKE ?", "%"+value.(string)+"%")
			} else {
				query = query.Where(key+" = ?", value)
			}
		}
	}

	if page > 0 && pageSize > 0 {
		query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	err := query.Find(&songs).Error
	if err != nil {
		log.Printf("Failed to get songs: %v", err)
	} else {
		log.Printf("Songs retrieved, count: %d", len(songs))
	}
	return songs, err
}

func (r *Repository) GetSong(id int64, page, pageSize int) (models.Songs, []models.Verses, error) {
	var song models.Songs
	if err := r.conn.First(&song, id).Error; err != nil {
		log.Printf("Song not found, id: %d, error: %v", id, err)
		return song, nil, err
	}

	var verses []models.Verses
	query := r.conn.Where("song_id = ?", id)
	if page > 0 && pageSize > 0 {
		query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	err := query.Find(&verses).Error
	if err != nil {
		log.Printf("Failed to get verses, song_id: %d, error: %v", id, err)
	} else {
		log.Printf("Verses retrieved, song_id: %d, count: %d", id, len(verses))
	}
	return song, verses, err
}

func (r *Repository) DeleteSong(id int64) error {
	err := r.conn.Delete(&models.Songs{}, id).Error
	if err != nil {
		log.Printf("Failed to delete song, id: %d, error: %v", id, err)
	} else {
		log.Printf("Song deleted, id: %d", id)
	}
	return err
}

func (r *Repository) UpdateSong(id int64, updatedSong models.Songs) (models.Songs, error) {
	var song models.Songs
	if err := r.conn.First(&song, id).Error; err != nil {
		log.Printf("Song not found, id: %d, error: %v", id, err)
		return song, err
	}

	updatedSong.UpdatedAt = time.Now()
	err := r.conn.Model(&song).Updates(updatedSong).Error
	if err != nil {
		log.Printf("Failed to update song, id: %d, error: %v", id, err)
	} else {
		log.Printf("Song updated, id: %d", id)
	}
	return updatedSong, err
}

func (r *Repository) CreateSong(song models.Songs, verses []models.Verses) (models.Songs, error) {
	tx := r.conn.Begin()
	if err := tx.Create(&song).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to create song: %v", err)
		return song, err
	}

	for i := range verses {
		verses[i].SongID = song.ID
		if err := tx.Create(&verses[i]).Error; err != nil {
			tx.Rollback()
			log.Printf("Failed to create verse, song_id: %d, error: %v", song.ID, err)
			return song, err
		}
	}

	tx.Commit()
	log.Printf("Song and verses created, id: %d, verses_count: %d", song.ID, len(verses))
	return song, nil
}
