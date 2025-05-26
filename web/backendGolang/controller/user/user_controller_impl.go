package usercontroller

import "github.com/labstack/echo/v4"

type UserController interface {
	GetListUserController(c echo.Context) error
	GetUserByIdController(c echo.Context) error
	PostUserController(c echo.Context) error
	DeleteProdukId(c echo.Context) error 
	 UpdateUserByIdController(c echo.Context) error
	 UpdatePasswordController(c echo.Context) error 
}