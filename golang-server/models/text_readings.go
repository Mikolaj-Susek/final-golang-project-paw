package models

import "time"

type TextReadings struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	InsertDate time.Time `json:"insertDate"`
	OcrText    string    `json:"ocrText"`
	FilePath   string    `json:"filePath"`
}
