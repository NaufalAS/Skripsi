package datarepo

import "skripsi/model/domain"

type DataRepository interface {
	SaveData(data domain.Data) (domain.Data, error)
	GetListData(filters string, limit int, page int) ([]domain.Data, int, int, int, *int, *int, error)
	GetDataById(id int)(domain.Data, error)
	DeleteDataId(iddata int)error
	UpdateDaataId(idData int, data domain.Data) (domain.Data, error) 
}
