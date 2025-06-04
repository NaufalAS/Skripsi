package datacontroller

import (
	"log"
	"mime/multipart"
	"net/http"
	"skripsi/helper"
	"skripsi/model"
	"skripsi/model/entity"
	"skripsi/model/web"
	dataservice "skripsi/service/data"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type DataControllerImpl struct {
	DataService dataservice.DataService
}

func NewDataController(DataService dataservice.DataService) *DataControllerImpl {
	return &DataControllerImpl{
		DataService: DataService,
	}
}

func (controller *DataControllerImpl) PostDataController(c echo.Context) error {
	// Ambil file dari form
	fileHeader, err := c.FormFile("gambar")
	if err != nil {
		log.Println("❌ Tidak bisa baca file:", err)
		return c.JSON(http.StatusBadRequest, model.ResponseToClient(http.StatusBadRequest, "file gambar wajib diisi", nil))
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ResponseToClient(http.StatusInternalServerError, "gagal membuka file gambar", nil))
	}
	defer file.Close()

	// Ambil dan parse tanggal dari form
	dateString := c.FormValue("date")
	parsedDate, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ResponseToClient(http.StatusBadRequest, "format tanggal tidak valid, gunakan YYYY-MM-DD", nil))
	}

	// Ambil data form lain
	request := web.PostDataRequest{
		JenisKendaraan:   c.FormValue("jenis_kendaraan"),
		JenisPelanggaran: c.FormValue("jenis_pelanggaran"),
		Lokasi:           c.FormValue("lokasi"),
		Kecepatan:           c.FormValue("kecepatan"),
		Date:             parsedDate,
		Gambar:           fileHeader.Filename,
	}

	// Validasi data form (jika kamu menggunakan validator Echo)
	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.ResponseToClient(http.StatusBadRequest, err.Error(), nil))
	}

	// Kirim ke service untuk proses penyimpanan data dan file
	saveData, errSave := controller.DataService.SaveData(request, file, fileHeader.Filename)
	if errSave != nil {
		return c.JSON(http.StatusInternalServerError, model.ResponseToClient(http.StatusInternalServerError, errSave.Error(), nil))
	}

	return c.JSON(http.StatusOK, model.ResponseToClient(http.StatusOK, "data berhasil disimpan", saveData))
}

func (controller *DataControllerImpl) GetListDataController(c echo.Context) error {
	// Extract query parameters for filters, limit, and page
	filters, limit, page := helper.ExtractFilter(c.QueryParams())

	// Call service to get data list and pagination data
	getDataController, totalCount, currentPage, totalPages, nextPage, prevPage, errGetDataController := controller.DataService.GetDataList(filters, limit, page)
	if errGetDataController != nil {
		return c.JSON(http.StatusInternalServerError, model.ResponseToClient(http.StatusInternalServerError, "error", errGetDataController.Error()))
	}

	if getDataController == nil {
		getDataController = []entity.DataEntity{}
	}

	// Prepare pagination data
	pagination := model.Pagination{
		CurrentPage:  currentPage,
		NextPage:     nextPage,
		PrevPage:     prevPage,
		TotalPages:   totalPages,
		TotalRecords: totalCount,
	}

	// Create response with pagination and data
	response := model.ResponseToClientpagi(http.StatusOK, "true", "Successfully fetched all data", pagination, getDataController)

	return c.JSON(http.StatusOK, response)
}


func (controller *DataControllerImpl) GetDataByIdController(c echo.Context) error{
	IdData := c.Param("id")
	id, _ := strconv.Atoi(IdData)
	getDataIdController, errGetDataController := controller.DataService.GetDataById(id)

	if errGetDataController != nil {
		return c.JSON(http.StatusInternalServerError, model.ResponseToClient(http.StatusInternalServerError, "error", errGetDataController.Error()))
	}

	return c.JSON(http.StatusOK, model.ResponseToClient(http.StatusOK,  "berhasil melihat satu id", getDataIdController))
}

func (controller *DataControllerImpl) DeleteDataId(c echo.Context) error {
	// Ambil ID dari URL dan konversi ke UUID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ResponseToClient(http.StatusBadRequest, "Invalid ID format", "invalid request"))
	}

	if errDeleteProduk := controller.DataService.DeleteData(id); errDeleteProduk != nil {
		return c.JSON(http.StatusBadRequest, model.ResponseToClient(http.StatusBadRequest, "error", errDeleteProduk.Error()))
	}

	return c.JSON(http.StatusOK, model.ResponseToClient(http.StatusOK, "Delete data Success", err))
}

func (controller *DataControllerImpl) UpdateDataByIdController(c echo.Context) error {
	// Ambil ID dari parameter
	userIdStr := c.Param("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ResponseToClient(http.StatusBadRequest, "Invalid user ID", nil))
	}

	// Ambil form values
	name := c.FormValue("jeniskendaraan")
	email := c.FormValue("jenispelanggaran")
	kecepatan := c.FormValue("kecepatan")
	phoneNumber := c.FormValue("lokasi")
	dateString := c.FormValue("date")

	var parsedDate time.Time
	var isDateProvided bool
	if dateString != "" {
		parsedDate, err = time.Parse("2006-01-02", dateString)
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.ResponseToClient(http.StatusBadRequest, "format tanggal tidak valid, gunakan YYYY-MM-DD", nil))
		}
		isDateProvided = true
	}

	// Ambil file gambar
	var openedFile multipart.File
	var fileName string
	formFile, err := c.FormFile("gambar")
	if err == http.ErrMissingFile {
		openedFile = nil
		fileName = ""
	} else if err != nil {
		return c.JSON(http.StatusBadRequest, model.ResponseToClient(http.StatusBadRequest, "Failed to get uploaded file", nil))
	} else {
		openedFile, err = formFile.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.ResponseToClient(http.StatusInternalServerError, "Failed to open uploaded file", nil))
		}
		defer openedFile.Close()
		fileName = formFile.Filename
	}

	// Siapkan request
	req := web.UpdateDataRequest{
		JenisKendaraan:   name,
		JenisPelanggaran: email,
		Lokasi:           phoneNumber,
		Gambar:           fileName,
		Kecepatan: 		  kecepatan,
	}
	if isDateProvided {
		req.Date = parsedDate
	}

	// Panggil service
	result, err := controller.DataService.UpdateDataId(userId, req, openedFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ResponseToClient(http.StatusInternalServerError, err.Error(), nil))
	}

	return c.JSON(http.StatusOK, model.ResponseToClient(http.StatusOK, "Successfully updated", result))
}
