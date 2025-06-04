package entity

import (
	"skripsi/model/domain"
	"time"
)

type DataEntity struct {
	Id int    `json:"id"`
	JenisKendaraan   string `json:"jeniskendaraan"`
	JenisPelanggaran string `json:"jenispelanggaran"`
	Lokasi string `json:"lokasi"`
	Date time.Time `json:"date"`
	Kecepatan string `json:data"kecepatan"`
	Gambar string `json:"gambar"`
}

func ToDataEntity(data domain.Data) DataEntity {
	return DataEntity{
		Id: data.ID,
		JenisKendaraan: data.JenisKendaraan,
		JenisPelanggaran: data.JenisPelanggaran,
		Lokasi: data.Lokasi,
		Date: data.Date,
		Gambar: data.Gambar,
		Kecepatan: data.Kecepatan,
	}
}

func ToDataListEntity(data []domain.Data) []DataEntity {
	var userData []DataEntity

	for _, data := range data {
		userData = append(userData, ToDataEntity(data))

	}
	return userData
}