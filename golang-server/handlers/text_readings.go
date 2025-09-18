package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/example/golang-postgres-crud/db"
	"github.com/example/golang-postgres-crud/models"
	ocr "github.com/example/golang-postgres-crud/ocr_service"
	"github.com/gin-gonic/gin"
)

// CreateTextReading godoc
// @Summary      Upload an image and perform OCR
// @Description  Uploads an image, performs OCR, saves the image to a static folder, and stores the data in the database.
// @Tags         text-readings
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Image file to upload (JPEG/PNG)"
// @Success      201 {object} models.TextReadings
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /api/text-readings [post]
func CreateTextReading(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	imageBytes, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	ocrService, err := ocr.NewOcrService()
	if err != nil {
		log.Printf("Error connecting to gRPC server: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to OCR server"})
		return
	}
	defer ocrService.Close()

	ocrText, err := ocrService.PerformOcr(imageBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not perform OCR operation"})
		return
	}

	hash := sha256.New()
	hash.Write(imageBytes)
	hash.Write([]byte(time.Now().String()))
	hashString := hex.EncodeToString(hash.Sum(nil))[:5]
	filename := fmt.Sprintf("%s_%s", hashString, file.Filename)
	filePath := filepath.Join("static", "images", filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	textReading := models.TextReadings{
		FileSize: file.Size,
		FilePath: filePath,
		OcrText:  ocrText,
	}

	db.DB.Create(&textReading)

	c.JSON(http.StatusCreated, textReading)
}

// GetTextReadings godoc
// @Summary      Get all text readings
// @Description  Retrieves a list of all text reading records from the database
// @Tags         text-readings
// @Produce      json
// @Success      200 {array} models.TextReadings
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
// @Success      200 {object} models.TextReadings
// @Failure      404 {object} map[string]string
// @Router       /api/text-readings/{id} [get]
func GetTextReading(c *gin.Context) {
	var textReading models.TextReadings
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := db.DB.First(&textReading, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TextReading not found"})
		return
	}
	c.JSON(http.StatusOK, textReading)
}

// UpdateTextReading godoc
// @Summary      Update an existing text reading's OCR text
// @Description  Updates the OcrText field of a text reading record by its ID.
// @Tags         text-readings
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Text Reading ID"
// @Param        input body      object true "The new OcrText data"
// @Success      200 {object} models.TextReadings
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api/text-readings/{id} [put]
func UpdateTextReading(c *gin.Context) {
	var textReading models.TextReadings
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := db.DB.First(&textReading, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TextReading not found"})
		return
	}

	var input struct {
		OcrText string `json:"ocrText"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	textReading.OcrText = input.OcrText
	db.DB.Save(&textReading)

	c.JSON(http.StatusOK, textReading)
}

// DeleteTextReading godoc
// @Summary      Delete a text reading
// @Description  Removes a text reading record from the database and deletes the corresponding image file from disk.
// @Tags         text-readings
// @Produce      json
// @Param        id   path      int  true  "Text Reading ID"
// @Success      200 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /api/text-readings/{id} [delete]
func DeleteTextReading(c *gin.Context) {
	var textReading models.TextReadings
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := db.DB.First(&textReading, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TextReading not found"})
		return
	}

	if _, err := os.Stat(textReading.FilePath); err == nil {
		if err := os.Remove(textReading.FilePath); err != nil {
			log.Printf("Failed to delete file %s: %v", textReading.FilePath, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated file"})
			return
		}
	} else if !os.IsNotExist(err) {
		log.Printf("Failed to check file existence %s: %v", textReading.FilePath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check file existence"})
		return
	}

	db.DB.Delete(&textReading)

	c.JSON(http.StatusOK, gin.H{"message": "TextReading and associated file deleted"})
}

// GetTextReadingImage godoc
// @Summary      Get image by text reading ID
// @Description  Retrieves the image file associated with a text reading record.
// @Tags         text-readings
// @Produce      image/jpeg
// @Param        id   path      int  true  "Text Reading ID"
// @Success      200 {file} file
// @Failure      404 {object} map[string]string
// @Router       /api/text-readings/{id}/image [get]
func GetTextReadingImage(c *gin.Context) {
	var textReading models.TextReadings
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := db.DB.First(&textReading, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TextReading not found"})
		return
	}

	if _, err := os.Stat(textReading.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image file not found"})
		return
	}

	c.File(textReading.FilePath)
}
