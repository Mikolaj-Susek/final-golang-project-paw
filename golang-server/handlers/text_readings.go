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

// GetTextReadings godoc
// @Summary      Get all text readings
// @Description  Retrieves a list of all text reading records from the database
// @Tags         text-readings
// @Produce      json
// @Success      200  {array}   models.TextReadings
// @Router       /api/text-readings [get]
func GetTextReadings(c *gin.Context) {
	var textReadings []models.TextReadings
	db.DB.Find(&textReadings)
	c.JSON(http.StatusOK, textReadings)
}

// GetTextReading godoc
// @Summary      Get a single text reading by ID
// @Description  Retrieves a text reading record based on its primary key
// @Tags         text-readings
// @Produce      json
// @Param        id   path      int  true  "Text Reading ID"
// @Success      200  {object}  models.TextReadings
// @Failure      404  {object}  map[string]string
// @Router       /api/text-readings/{id} [get]
func GetTextReading(c *gin.Context) {
	var textReading models.TextReadings
	if err := db.DB.First(&textReading, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TextReading not found"})
		return
	}
	c.JSON(http.StatusOK, textReading)
}

// UpdateTextReading godoc
// @Summary      Update an existing text reading
// @Description  Updates a text reading record in the database by its ID
// @Tags         text-readings
// @Accept       json
// @Produce      json
// @Param        id       path      int                  true  "Text Reading ID"
// @Param        input    body      models.TextReadings  true  "The new data for the text reading"
// @Success      200      {object}  models.TextReadings
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Router       /api/text-readings/{id} [put]
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

// DeleteTextReading godoc
// @Summary      Delete a text reading
// @Description  Removes a text reading record from the database by its ID
// @Tags         text-readings
// @Produce      json
// @Param        id   path      int  true  "Text Reading ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/text-readings/{id} [delete]
func DeleteTextReading(c *gin.Context) {
	var textReading models.TextReadings
	if err := db.DB.First(&textReading, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TextReading not found"})
		return
	}
	db.DB.Delete(&textReading)
	c.JSON(http.StatusOK, gin.H{"message": "TextReading deleted"})
}
