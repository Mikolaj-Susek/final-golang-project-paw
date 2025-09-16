package models

import (
	"gorm.io/gorm"
)

type TextReadings struct {
	gorm.Model
	OcrText  string `json:"ocrText"`
	FilePath string `json:"filePath"`
}
