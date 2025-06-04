package analiticscontroller

import (
	"net/http"
	"skripsi/model"
	analiticsservice "skripsi/service/analitics"

	"github.com/labstack/echo/v4"
)

type AnaliticsControllerImpl struct {
	analiticsService analiticsservice.AnaliticsService
}

func NewAnaliticsController(analiticsService analiticsservice.AnaliticsService) *AnaliticsControllerImpl {
	return &AnaliticsControllerImpl{
		analiticsService: analiticsService,
	}
}

func (controller *AnaliticsControllerImpl) GetAnalyticsDataController(c echo.Context) error {
	// Panggil service untuk mendapatkan data analitik
	analyticsTotal, err := controller.analiticsService.GetAnalyticsData()
	if err != nil {
		// Return error response if something went wrong
		return c.JSON(http.StatusInternalServerError, model.ResponseToClient(http.StatusInternalServerError, "error", err.Error()))
	}

	// Return successful response
	return c.JSON(http.StatusOK, model.ResponseToClient(http.StatusOK, "Berhasil mendapatkan data analitik", analyticsTotal))
}