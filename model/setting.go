package model

import "time"

type Setting struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	Key       string `gorm:"not null; uniqueIndex"`
	Value     string `gorm:"not null;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
