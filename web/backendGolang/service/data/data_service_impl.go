package dataservice

import (
	"mime/multipart"
	"skripsi/model/entity"
	"skripsi/model/web"
)

type DataService interface {
	SaveData(request web.PostDataRequest, file multipart.File, filename string) (map[string]interface{}, error)
	GetUser() ([]entity.DataEntity, error) 
	GetDataById(id int)(entity.DataEntity, error)
	DeleteData(id int) error
	UpdateDataId(Id int, req web.UpdateDataRequest, file multipart.File) (map[string]interface{}, error) 
}
