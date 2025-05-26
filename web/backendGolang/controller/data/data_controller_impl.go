package datacontroller

import "github.com/labstack/echo/v4"

type DataController interface {
	PostDataController(c echo.Context) error
	GetListDataController(c echo.Context) error
	GetDataByIdController(c echo.Context) error
	DeleteDataId(c echo.Context) error
	UpdateDataByIdController(c echo.Context) error 
}