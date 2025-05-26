package usercontroller

import (
	"mime/multipart"
	"net/http"
	"skripsi/model"
	"skripsi/model/web"
	userservice "skripsi/service/user"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserControllerImpl struct {
	UserService userservice.UserService
}

func NewSektorUsahaController(UserService userservice.UserService) *UserControllerImpl {
	return &UserControllerImpl{
		UserService: UserService,
	}
}

func (controller *UserControllerImpl) PostUserController(c echo.Context) error {
	user := new(web.LoginUserRequest)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, model.ResponseToClient(http.StatusBadRequest, err.Error(), nil))
	}

	if err := c.Validate(user); err != nil {
		return err
	}

	saverUser, errSaveuser := controller.UserService.SaveUser(*user)
	if errSaveuser != nil {
		return c.JSON(http.StatusInternalServerError, model.ResponseToClient(http.StatusInternalServerError, errSaveuser.Error(), nil))
	}
	return c.JSON(http.StatusOK, model.ResponseToClient(http.StatusOK, "Register Sukses", saverUser))
}

func(controller *UserControllerImpl) LoginUserController(c echo.Context)error{
	loginUser := new(web.LoginUserRequest)

	if err := c.Bind(loginUser); err != nil {
		return c.JSON(http.StatusBadRequest, model.ResponseToClient(http.StatusBadRequest, err.Error(), nil))
	}

	if err := c.Validate(loginUser); err != nil {
		return err
	}

	userLogin, errLogin := controller.UserService.Login(loginUser.Name, loginUser.Password)

	if errLogin != nil {
		return c.JSON(http.StatusBadRequest, model.ResponseToClient(http.StatusBadRequest, errLogin.Error(), nil))
	}

	return c.JSON(http.StatusOK, model.ResponseToClient(http.StatusOK, "Login Success", userLogin))
}


func (controller *UserControllerImpl) GetListUserController(c echo.Context) error {
	getUserController, errGetUserController := controller.UserService.GetUser()

	if errGetUserController != nil {
		return c.JSON(http.StatusInternalServerError, model.ResponseToClient(http.StatusInternalServerError, "error", errGetUserController.Error()))
	}

	return c.JSON(http.StatusOK, model.ResponseToClient(http.StatusOK,  "berhasil melihat seluruh list User", getUserController))
}

func (controller *UserControllerImpl) GetUserByIdController(c echo.Context) error{
	IdUser := c.Param("id")
	id, _ := strconv.Atoi(IdUser)
	getUserIdController, errGetUserController := controller.UserService.GetUserById(id)

	if errGetUserController != nil {
		return c.JSON(http.StatusInternalServerError, model.ResponseToClient(http.StatusInternalServerError, "error", errGetUserController.Error()))
	}

	return c.JSON(http.StatusOK, model.ResponseToClient(http.StatusOK,  "berhasil melihat satu id", getUserIdController))
}

func (controller *UserControllerImpl)  UpdateUserByIdController(c echo.Context) error {
    // Ambil ID dari parameter dan ubah ke int
    userIdStr := c.Param("id")
    userId, err := strconv.Atoi(userIdStr)
    if err != nil {
        return c.JSON(http.StatusBadRequest, model.ResponseToClient(http.StatusBadRequest, "Invalid user ID", nil))
    }
    // Ambil form values
    name := c.FormValue("fullname")
    email := c.FormValue("email")
    phoneNumber := c.FormValue("phone_number")
    address := c.FormValue("alamat")

    // Ambil file potoprofile
    formFile, err := c.FormFile("potoprofile")
    var openedFile multipart.File
    if err == http.ErrMissingFile {
        openedFile = nil
    } else if err != nil {
        return c.JSON(http.StatusBadRequest, model.ResponseToClient(http.StatusBadRequest,  "Failed to get uploaded file", nil))
    } else {
        openedFile, err = formFile.Open()
        if err != nil {
            return c.JSON(http.StatusInternalServerError, model.ResponseToClient(http.StatusInternalServerError, "Failed to open uploaded file", nil))
        }
        defer openedFile.Close()
    }

    // Buat struct request
    req := web.UpdateUserRequest{
        Name:      name,
        Email:     email,
        NoTelepon: phoneNumber,
        Alamat:    address,
        Foto:      formFile.Filename, // Nama file
    }

    // Panggil service
    result, err := controller.UserService.UpdateUserId(userId, req, openedFile)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, model.ResponseToClient(http.StatusInternalServerError, err.Error(), nil))
    }

    return c.JSON(http.StatusOK, model.ResponseToClient(http.StatusOK, "Successfully updated", result))
}



func (controller *UserControllerImpl) DeleteProdukId(c echo.Context) error {
	// Ambil ID dari URL dan konversi ke UUID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ResponseToClient(http.StatusBadRequest, "Invalid ID format", "invalid request"))
	}

	if errDeleteProduk := controller.UserService.DeleteProduk(id); errDeleteProduk != nil {
		return c.JSON(http.StatusBadRequest, model.ResponseToClient(http.StatusBadRequest, "error", errDeleteProduk.Error()))
	}

	return c.JSON(http.StatusOK, model.ResponseToClient(http.StatusOK, "Delete User Success", err))
}

func (controller *UserControllerImpl) UpdatePasswordController(c echo.Context) error {
    // Ambil ID user dari URL parameter
    IdUser := c.Param("id")
    id, _ := strconv.Atoi(IdUser)

    // Bind request body ke struct UpdatePasswordRequest
    request := new(web.UpdatePasswordRequest)
    if err := c.Bind(request); err != nil {
        return c.JSON(http.StatusBadRequest, model.ResponseToClient(http.StatusBadRequest, "error", "invalid request"))
    }

    // Panggil service untuk memperbarui password
    updatedUser, err := controller.UserService.UpdatePassword(id, request.Password, request.NewPassword, *request)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, model.ResponseToClient(http.StatusInternalServerError, "error", err.Error()))
    }

    // Kembalikan response sukses
    return c.JSON(http.StatusOK, model.ResponseToClient(http.StatusOK, "berhasil update password", updatedUser))
}