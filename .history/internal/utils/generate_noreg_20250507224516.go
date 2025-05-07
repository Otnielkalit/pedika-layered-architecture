package utils

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func GenerateNomorRegistrasi(db *gorm.DB) (string, error) {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	// Hitung total laporan pada bulan dan tahun yang sama
	var count int64
	err := db.Table("reports").Where("EXTRACT(YEAR FROM tanggal_pelaporan) = ? AND EXTRACT(MONTH FROM tanggal_pelaporan) = ?", year, month).Count(&count).Error
	if err != nil {
		return "", err
	}

	urutan := fmt.Sprintf("%03d", count+1)
	format := fmt.Sprintf("%s-DPMDPPA-%d-%d", urutan, month, year)
	return format, nil
}
