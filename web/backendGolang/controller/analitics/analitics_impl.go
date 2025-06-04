package analiticscontroller

import "github.com/labstack/echo/v4"

type AnaliticsController interface {
	GetAnalyticsDataController(c echo.Context) error
}