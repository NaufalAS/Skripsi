package domain

import "time"

type Data struct {
	ID        int    `gorm:"column:id;primaryKey;autoIncrement"`
	JenisKendaraan      string `gorm:"column:jeniskendaraan"`
	JenisPelanggaran      string `gorm:"column:jenispelanggaran"`
	Lokasi  string `gorm:"column:lokasi"`
	Date time.Time `gorm:"column:waktu"`
	Gambar string `gorm:"column:gambar"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Data) TableName() string {
	return "data"
}