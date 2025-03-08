package service

import (
	"GoSongs/internal/models"
	"GoSongs/internal/repository"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type Service struct {
	repo       *repository.Repository
	apiURL     string
	httpClient *http.Client
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func NewService(repo *repository.Repository, apiURL string) *Service {
	return &Service{
		repo:       repo,
		apiURL:     apiURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *Service) GetSongs(filters map[string]interface{}, page, pageSize int) ([]models.Songs, error) {
	return s.repo.GetSongs(filters, page, pageSize)
}

func (s *Service) GetSong(id int64, page, pageSize int) (models.Songs, []models.Verses, error) {
	return s.repo.GetSong(id, page, pageSize)
}

func (s *Service) DeleteSong(id int64) error {
	return s.repo.DeleteSong(id)
}

func (s *Service) UpdateSong(id int64, updatedSong models.Songs) (models.Songs, error) {
	updatedSong.ID = id
	updatedSong.UpdatedAt = time.Now()
	return s.repo.UpdateSong(id, updatedSong)
}

func (s *Service) CreateSong(song models.Songs) (models.Songs, error) {
	log.Printf("Fetching song details from external API, group: %s, song: %s", song.Group, song.Song)

	req, err := http.NewRequest("GET", s.apiURL+"/info", nil)
	if err != nil {
		log.Printf("Failed to create API request: %v", err)
		return song, err
	}

	q := req.URL.Query()
	q.Add("group", song.Group)
	q.Add("song", song.Song)
	req.URL.RawQuery = q.Encode()

	resp, err := s.httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch from external API: %v", err)
		return song, err
	}
	defer resp.Body.Close()

	var songDetail SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		log.Printf("Failed to decode API response: %v", err)
		return song, err
	}

	releaseDate, err := time.Parse("02.01.2006", songDetail.ReleaseDate)
	if err != nil {
		log.Printf("Failed to parse release date: %v", err)
		return song, err
	}

	song.Release = releaseDate
	song.Link = songDetail.Link
	song.CreatedAt = time.Now()
	song.UpdatedAt = time.Now()

	verses := splitTextToVerses(songDetail.Text)
	var verseModels []models.Verses
	for i, text := range verses {
		if text != "" {
			verseModels = append(verseModels, models.Verses{
				Number:    i + 1,
				Text:      text,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}

	createdSong, err := s.repo.CreateSong(song, verseModels)
	return createdSong, err
}

func splitTextToVerses(text string) []string {
	lines := strings.Split(text, "\n\n")
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	return result
}
