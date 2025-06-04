package analiticsrepo

import (
	"errors"
	"skripsi/model/domain"
	"time"

	"gorm.io/gorm"
)

type AnaliticsRepoImpl struct {
	DB *gorm.DB
}

func NewAnalticsRepository(DB *gorm.DB) *AnaliticsRepoImpl{
	return &AnaliticsRepoImpl{
		DB: DB,
	}
}

func (repo *AnaliticsRepoImpl) GetAnalitics() (int64, error) {
	var count int64

	// Menghitung jumlah total baris dalam tabel "data"
	err := repo.DB.Model(&domain.Data{}).Select("COUNT(*)").Scan(&count).Error
	if err != nil {
		return 0, err
	}

	// Jika tidak ada data yang dihitung
	if count == 0 {
		return 0, errors.New("no data found")
	}

	return count, nil
}

func (repo *AnaliticsRepoImpl) GetAnaliticsPerDay() (map[string]int64, error) {
	// Map untuk menyimpan hasil jumlah data per hari
	result := make(map[string]int64)

	// Mengambil tanggal 5 hari dari sekarang hingga 5 hari ke belakang
	currentDate := time.Now()
	dateRange := []string{}
	for i := -5; i <= 0; i++ {
		// Hitung setiap hari dalam rentang waktu (5 hari ke belakang hingga hari ini)
		dateRange = append(dateRange, currentDate.Add(time.Hour*24*time.Duration(i)).Format("2006-01-02"))
	}

	// Query untuk menghitung jumlah data per hari
	for _, date := range dateRange {
		var count int64
		err := repo.DB.Model(&domain.Data{}).
			Where("DATE(waktu) = ?", date).  // Filter berdasarkan tanggal
			Select("COUNT(*)").
			Scan(&count).Error
		if err != nil {
			return nil, err
		}
		result[date] = count
	}

	// Mengembalikan hasil dalam bentuk map yang berisi jumlah data per tanggal
	return result, nil
}

