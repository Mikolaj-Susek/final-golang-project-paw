package models

import (
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	FileSize int64  `json:"fileSize"`
	FilePath string `json:"filePath"`
}
