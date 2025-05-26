package dataservice

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"skripsi/helper"
	"skripsi/model/domain"
	"skripsi/model/entity"
	"skripsi/model/web"
	datarepo "skripsi/repository/data"
)

type DataServiceImpl struct {
	datarepo datarepo.DataRepository
}

func NewSektorDataService(datarepo datarepo.DataRepository) *DataServiceImpl {
	return &DataServiceImpl{
		datarepo: datarepo,
	}
}


// Fungsi untuk menyimpan file foto
func saveFile(file multipart.File, destination string) error {
	// Membuat file baru di lokasi tujuan
	out, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	// Salin konten dari file upload ke file baru
	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	return nil
}

func (service *DataServiceImpl) SaveData(request web.PostDataRequest, file multipart.File, filename string) (map[string]interface{}, error) {
	// Simpan file gambar ke folder `public/pelanggaran`
	if file != nil {
		dir := "public/pelanggaran"
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("gagal membuat folder penyimpanan: %v", err)
		}

		// Buat nama file unik
		newFileName := helper.GenerateRandomFileName(filepath.Ext(filename))
		newFilePath := filepath.Join(dir, newFileName)

		// Simpan file
		if err := saveFile(file, newFilePath); err != nil {
			return nil, fmt.Errorf("gagal menyimpan file gambar: %v", err)
		}

		// Simpan path relatif ke database
		request.Gambar = "/public/pelanggaran/" + newFileName
	}

	// Simpan data ke database
	newData := domain.Data{
		JenisKendaraan:   request.JenisKendaraan,
		JenisPelanggaran: request.JenisPelanggaran,
		Lokasi:           request.Lokasi,
		Date:             request.Date,
		Gambar:           request.Gambar,
	}

	savedata, err := service.datarepo.SaveData(newData)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"message":           "data berhasil disimpan",
		"id":                savedata.ID,
		"jenis_kendaraan":   savedata.JenisKendaraan,
		"jenis_pelanggaran": savedata.JenisPelanggaran,
		"lokasi":            savedata.Lokasi,
		"tanggal":           savedata.Date,
		"gambar":            savedata.Gambar,
	}

	return response, nil
}

func (service *DataServiceImpl) GetUser() ([]entity.DataEntity, error) {
	getData, err := service.datarepo.GetListData()
	if err != nil {
		return nil, err
	}

	dataEntities := entity.ToDataListEntity(getData)

	return dataEntities, nil
}

func(service *DataServiceImpl) GetDataById(id int)(entity.DataEntity, error){
	getDataService, err := service.datarepo.GetDataById(id)
	if err != nil {
		return entity.DataEntity{}, err
	}

	return entity.ToDataEntity(getDataService), nil
}

func (service *DataServiceImpl) DeleteData(id int) error {
	// Cari produk berdasarkan ID
	data, err := service.datarepo.GetDataById(id)
	if err != nil {
		return err
	}

	// Periksa apakah Gambar mengandung prefix 'public/pelanggaran/'
imageFileName := filepath.Base(data.Gambar) // hanya ambil nama file dari path apa pun
imagePath := filepath.Join("public", "pelanggaran", imageFileName)

fmt.Println("Path gambar: ", imagePath)

if _, err := os.Stat(imagePath); err == nil {
	if err := os.Remove(imagePath); err != nil {
		return fmt.Errorf("gagal menghapus gambar: %v", err)
	}
}
	// Hapus produk dari database
	return service.datarepo.DeleteDataId(id)
}

//update
// func (service *DataServiceImpl) UpdateDataId(Id int, req web.UpdateDataRequest, file multipart.File) (map[string]interface{}, error) {
// 	// Cek apakah user ada
// 	_, errUser := service.datarepo.GetDataById(Id)
// 	if errUser != nil {
// 		return nil, fmt.Errorf("user not found")
// 	}

// 	// Jika ada foto baru yang diupload
// 	if file != nil {
// 		// Ambil data user lama
// 		oldUser, _ := service.datarepo.GetDataById(Id)
// 		oldProfilePath := filepath.Join("public", "pelanggaran", filepath.Base(oldUser.Gambar))

// 		// Hapus foto lama jika ada
// 		if _, err := os.Stat(oldProfilePath); err == nil {
// 			if err := os.Remove(oldProfilePath); err != nil {
// 				return nil, fmt.Errorf("failed to delete old foto pelabggaran: %v", err)
// 			}
// 		}

// 		// Tentukan nama file baru
// 		newFileName := helper.GenerateRandomFileName(filepath.Ext(req.Gambar))
// 		newProfilePath := filepath.Join("public", "pelanggaran", newFileName)

// 		// Simpan file
// 		if err := saveFile(file, newProfilePath); err != nil {
// 			return nil, fmt.Errorf("failed to save profile picture: %v", err)
// 		}

// 		// Set path untuk disimpan ke database
// 		req.Gambar = "/public/pelanggaran/" + newFileName
// 	}

// 	// Persiapkan data untuk update
// 	updatedUser := domain.Data{
// 		ID:        Id,
// 		JenisKendaraan:  req.JenisKendaraan,
// 		Gambar:   req.Gambar,
// 		JenisPelanggaran:     req.JenisPelanggaran,
// 		Lokasi: req.Lokasi,
// 		Date:    req.Date,
// 	}

// 	// Update ke repo
// 	result, errUpdate := service.datarepo.UpdateDaataId(Id, updatedUser)
// 	if errUpdate != nil {
// 		return nil, errUpdate
// 	}

// 	// Response akhir
// 	response := helper.ResponseToJson{
// 		"name":       result.JenisKendaraan,
// 		"profile":    result.Gambar,
// 		"email":      result.JenisPelanggaran,
// 		"no_telepon": result.Lokasi,
// 		"alamat":     result.Date,
// 	}
// 	return response, nil
// }

func (service *DataServiceImpl) UpdateDataId(Id int, req web.UpdateDataRequest, file multipart.File) (map[string]interface{}, error) {
    // 1. Ambil data lama
    oldData, err := service.datarepo.GetDataById(Id)
    if err != nil {
        return nil, fmt.Errorf("data not found")
    }

    // 2. Siapkan struct baru dari data lama
    updatedData := oldData

    // 3. Update field yang ada di request
    if req.JenisKendaraan != "" {
        updatedData.JenisKendaraan = req.JenisKendaraan
    }
    if req.JenisPelanggaran != "" {
        updatedData.JenisPelanggaran = req.JenisPelanggaran
    }
    if req.Lokasi != "" {
        updatedData.Lokasi = req.Lokasi
    }
    if !req.Date.IsZero() {
        updatedData.Date = req.Date
    }

    // 4. Jika ada file gambar baru, simpan dan update path gambar
    if file != nil {
        // Hapus gambar lama
        oldProfilePath := filepath.Join("public", "pelanggaran", filepath.Base(oldData.Gambar))
        if _, err := os.Stat(oldProfilePath); err == nil {
            _ = os.Remove(oldProfilePath)
        }

        // Generate nama file baru dan simpan file
        newFileName := helper.GenerateRandomFileName(filepath.Ext(req.Gambar))
        newProfilePath := filepath.Join("public", "pelanggaran", newFileName)
        if err := saveFile(file, newProfilePath); err != nil {
            return nil, fmt.Errorf("failed to save image: %v", err)
        }

        updatedData.Gambar = "/public/pelanggaran/" + newFileName
    }

    // 5. Update data di repository
    result, errUpdate := service.datarepo.UpdateDaataId(Id, updatedData)
    if errUpdate != nil {
        return nil, errUpdate
    }

    // 6. Buat response
    response := map[string]interface{}{
        "jenis_kendaraan":   result.JenisKendaraan,
        "gambar":            result.Gambar,
        "jenis_pelanggaran": result.JenisPelanggaran,
        "lokasi":            result.Lokasi,
        "date":              result.Date,
    }

    return response, nil
}
