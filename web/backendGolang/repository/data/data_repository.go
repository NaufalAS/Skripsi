package datarepo

import (
	"skripsi/model/domain"

	"gorm.io/gorm"
)

type DataRepositoryImpl struct {
	DB *gorm.DB
}

func NewDataRepository(db *gorm.DB) *DataRepositoryImpl {
	return &DataRepositoryImpl{
		DB: db,
	}
}

func (repo *DataRepositoryImpl) SaveData(data domain.Data) (domain.Data, error) {
	err := repo.DB.Create(&data).Error
	if err != nil {
		return domain.Data{}, err
	}
	return data, nil
}

func (repo *DataRepositoryImpl) GetListData() ([]domain.Data, error) {
	var data []domain.Data
	err := repo.DB.Order("id asc").Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func(repo *DataRepositoryImpl) GetDataById(id int)(domain.Data, error){
	var data domain.Data

	if err := repo.DB.Find(&data,"id = ?", id).Error; err != nil{
		return data, err
	}

	return data, nil
}

func (repo *DataRepositoryImpl) DeleteDataId(iddata int)error {
	if err := repo.DB.Delete(&domain.Data{}, iddata).Error; err != nil {
        return err
    }
    return  nil
}

func (repo *DataRepositoryImpl) UpdateDaataId(idData int, data domain.Data) (domain.Data, error) {
    if err := repo.DB.Model(&domain.Data{}).Where("id = ?", idData).Updates(data).Error; err != nil {
        return domain.Data{}, err
    }
    return data, nil
}