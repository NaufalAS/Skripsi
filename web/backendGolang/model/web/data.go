package web

import "time"

type PostDataRequest struct {
	JenisKendaraan   string    `json:"jeniskendaraan"`
	JenisPelanggaran string    `json:"jenispelanggaran"`
	Lokasi           string    `json:"lokasi"`
	Date             time.Time `json:"date"`
	Gambar           string    `json:"gambar"`
	Kecepatan   	 string	   `json:"kecepatan"`
}

type UpdateDataRequest struct {
	JenisKendaraan   string    `json:"jeniskendaraan"`
	JenisPelanggaran string    `json:"jenispelanggaran"`
	Lokasi           string    `json:"lokasi"`
	Date             time.Time `json:"date"`
	Gambar           string    ` json:"gambar"`
	Kecepatan 		 string    `json:"kecepatan"`
}
