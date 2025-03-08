package handlers

import (
	"GoSongs/internal/models"
	"GoSongs/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetSongs(c *gin.Context) {
	filters := map[string]interface{}{
		"group": c.Query("group"),
		"song":  c.Query("song"),
		"link":  c.Query("link"),
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	songs, err := h.service.GetSongs(filters, page, pageSize)
	if err != nil {
		log.Printf("Failed to handle GetSongs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, songs)
}

func (h *Handler) GetSong(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Invalid ID in GetSong: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	song, verses, err := h.service.GetSong(id, page, pageSize)
	if err != nil {
		log.Printf("Song not found in GetSong, id: %d, error: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"song": song, "verses": verses})
}

func (h *Handler) DeleteSong(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Invalid ID in DeleteSong: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.service.DeleteSong(id); err != nil {
		log.Printf("Failed to delete song, id: %d, error: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdateSong(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Invalid ID in UpdateSong: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var updatedSong models.Songs
	if err := c.ShouldBindJSON(&updatedSong); err != nil {
		log.Printf("Failed to bind JSON in UpdateSong: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	song, err := h.service.UpdateSong(id, updatedSong)
	if err != nil {
		log.Printf("Failed to update song, id: %d, error: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, song)
}

func (h *Handler) CreateSong(c *gin.Context) {
	var song models.Songs
	if err := c.ShouldBindJSON(&song); err != nil {
		log.Printf("Failed to bind JSON in CreateSong: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdSong, err := h.service.CreateSong(song)
	if err != nil {
		log.Printf("Failed to create song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdSong)
}
