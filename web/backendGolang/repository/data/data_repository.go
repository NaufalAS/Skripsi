package datarepo

import (
	"skripsi/model/domain"
	query_builder_data "skripsi/query_builder/data_qeury"

	"gorm.io/gorm"
)

type DataRepositoryImpl struct {
	DB               *gorm.DB
	dataQueryBuilder query_builder_data.DataQueryBuilder
}

func NewDataRepository(db *gorm.DB, dataQueryBuilder query_builder_data.DataQueryBuilder) *DataRepositoryImpl {
	return &DataRepositoryImpl{
		DB:               db,
		dataQueryBuilder: dataQueryBuilder,
	}
}

func (repo *DataRepositoryImpl) SaveData(data domain.Data) (domain.Data, error) {
	err := repo.DB.Create(&data).Error
	if err != nil {
		return domain.Data{}, err
	}
	return data, nil
}

func (repo *DataRepositoryImpl) GetListData(filters string, limit int, page int) ([]domain.Data, int, int, int, *int, *int, error) {
	var data []domain.Data
	var totalcount int64

	// Get the query with filters, limit, and pagination
	query, err := repo.dataQueryBuilder.GetBuilderData(filters, limit, page)
	if err != nil {
		return nil, 0, 0, 0, nil, nil, err
	}

	// Execute the query with pagination
	err = query.Order("id ASC").Find(&data).Error
	if err != nil {
		return nil, 0, 0, 0, nil, nil, err
	}

	// Count the total number of records for pagination
	query, err = repo.dataQueryBuilder.GetBuilderData(filters, 0, 0) // Reset limit and page for count
	if err != nil {
		return nil, 0, 0, 0, nil, nil, err
	}

	err = query.Model(&domain.Data{}).Count(&totalcount).Error
	if err != nil {
		return nil, 0, 0, 0, nil, nil, err
	}

	// Calculate total pages
	totalPages := 1
	if limit > 0 {
		totalPages = int((totalcount + int64(limit) - 1) / int64(limit))
	}

	// If the requested page is higher than the total pages, return empty data
	if page > totalPages {
		return nil, int(totalcount), page, totalPages, nil, nil, nil
	}

	// Determine next and previous pages
	currentPage := page
	var nextPage *int
	if currentPage < totalPages {
		np := currentPage + 1
		nextPage = &np
	}

	var prevPage *int
	if currentPage > 1 {
		pp := currentPage - 1
		prevPage = &pp
	}

	return data, int(totalcount), currentPage, totalPages, nextPage, prevPage, nil
}

func (repo *DataRepositoryImpl) GetDataById(id int) (domain.Data, error) {
	var data domain.Data

	if err := repo.DB.Find(&data, "id = ?", id).Error; err != nil {
		return data, err
	}

	return data, nil
}

func (repo *DataRepositoryImpl) DeleteDataId(iddata int) error {
	if err := repo.DB.Delete(&domain.Data{}, iddata).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DataRepositoryImpl) UpdateDaataId(idData int, data domain.Data) (domain.Data, error) {
	if err := repo.DB.Model(&domain.Data{}).Where("id = ?", idData).Updates(data).Error; err != nil {
		return domain.Data{}, err
	}
	return data, nil
}
