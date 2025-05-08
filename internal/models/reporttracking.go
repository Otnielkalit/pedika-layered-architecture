package models

import (
	"time"

	"gorm.io/datatypes"
)

type ReportTracking struct {
	ID           uint            `gorm:"primaryKey"`
	NoRegistrasi string          `gorm:"not null"`
	Keterangan   string          `gorm:"type:text"`
	Document     datatypes.JSONMap `gorm:"type:json"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
