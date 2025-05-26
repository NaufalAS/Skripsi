package web

import "time"

type PostDataRequest struct {
	JenisKendaraan   string    `json:"jeniskendaraan"`
	JenisPelanggaran string    `json:"jenispelanggaran"`
	Lokasi           string    `json:"lokasi"`
	Date             time.Time `json:"date"`
	Gambar           string    `json:"gambar"`
}

type UpdateDataRequest struct {
	JenisKendaraan   string    `json:"jeniskendaraan"`
	JenisPelanggaran string    `json:"jenispelanggaran"`
	Lokasi           string    `json:"lokasi"`
	Date             time.Time `json:"date"`
	Gambar           string    ` json:"gambar"`
}
