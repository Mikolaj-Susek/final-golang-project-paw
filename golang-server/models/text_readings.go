package models

import (
	"gorm.io/gorm"
)

type TextReadings struct {
	gorm.Model
	FileSize int64  `json:"fileSize"`
	FilePath string `json:"filePath"`
	OcrText  string `json:"ocrText"`
}
