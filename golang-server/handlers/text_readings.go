package handlers

import (
	"net/http"

	"github.com/example/golang-postgres-crud/db"
	"github.com/example/golang-postgres-crud/models"
	"github.com/gin-gonic/gin"
)

// CreateTextReading godoc
// @Summary      Create a new text reading
// @Description  Adds a new text reading record to the database
// @Tags         text-readings
// @Accept       json
// @Produce      json
// @Param        reading  body      models.TextReadings  true  "Text Reading to add"
// @Success      201      {object}  models.TextReadings
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/text-readings [post]
func CreateTextReading(c *gin.Context) {
	var textReading models.TextReadings
	if err := c.ShouldBindJSON(&textReading); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Create(&textReading)
	c.JSON(http.StatusCreated, textReading)
}

func GetTextReadings(c *gin.Context) {
	var textReadings []models.TextReadings
	db.DB.Find(&textReadings)
	c.JSON(http.StatusOK, textReadings)
}

func GetTextReading(c *gin.Context) {
	var textReading models.TextReadings
	if err := db.DB.First(&textReading, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TextReading not found"})
		return
	}
	c.JSON(http.StatusOK, textReading)
}

func UpdateTextReading(c *gin.Context) {
	var textReading models.TextReadings
	if err := db.DB.First(&textReading, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TextReading not found"})
		return
	}

	var input models.TextReadings
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Model(&textReading).Updates(input)
	c.JSON(http.StatusOK, textReading)
}

func DeleteTextReading(c *gin.Context) {
	var textReading models.TextReadings
	if err := db.DB.First(&textReading, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TextReading not found"})
		return
	}
	db.DB.Delete(&textReading)
	c.JSON(http.StatusOK, gin.H{"message": "TextReading deleted"})
}
