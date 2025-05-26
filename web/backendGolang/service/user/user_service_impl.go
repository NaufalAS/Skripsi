package userservice

import (
	"mime/multipart"
	"skripsi/helper"
	"skripsi/model/entity"
	"skripsi/model/web"
)

type UserService interface {
	SaveUser(request web.LoginUserRequest) (map[string]interface{}, error)
	GetUser() ([]entity.UserEntity, error)
	GetUserById(id int) (entity.UserEntity, error)
	UpdateUserId(Id int, req web.UpdateUserRequest, file multipart.File) (map[string]interface{}, error)
	DeleteProduk(id int) error
	Login(name, password string) (helper.ResponseToJson, error)
	UpdatePassword(Id int, oldpassword string, newpassword string, req web.UpdatePasswordRequest) (map[string]interface{}, error)
}
