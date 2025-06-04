package query_builder_data

import (
	querybuilder "skripsi/query_builder"

	"gorm.io/gorm"
)

type DataQueryBuilder interface {
	querybuilder.BaseQueryBuilderList
	GetBuilderData(filters string, limit int, page int) (*gorm.DB, error)
	GetBuilderDataListWeb(limit int, page int, filters string) (*gorm.DB, error)
}

type DataQueryBuilderImpl struct {
	querybuilder.BaseQueryBuilderList
	db *gorm.DB
}

func NewDataQueryBuilder(db *gorm.DB) *DataQueryBuilderImpl {
	return &DataQueryBuilderImpl{
		BaseQueryBuilderList: querybuilder.NewBaseQueryBuilderList(db),
		db:                   db,
	}
}

func (DataQueryBuilder *DataQueryBuilderImpl) GetBuilderData(filters string, limit int, page int) (*gorm.DB, error) {
	query := DataQueryBuilder.db

	// Implementasi filter di sini
	if filters != "" {
		searchPattern := "%" + filters + "%"
		query = query.Where("jeniskendaraan ILIKE ?", searchPattern)
	}

	if limit <= 0 {
		limit = 15
	}

	query, err := DataQueryBuilder.GetQueryBuilderList(query, limit, page)
	if err != nil {
		return nil, err
	}

	// query = query.Preload("data")

	return query, nil
}

func (DataQueryBuilder *DataQueryBuilderImpl) GetBuilderDataListWeb(limit int, page int, filters string) (*gorm.DB, error) {
	query := DataQueryBuilder.db

	if filters != "" {
		searchPattern := "%" + filters + "%"
		query = query.Where("nama ILIKE ?", searchPattern)
	}

	if limit <= 0 {
		limit = 15
	}

	query, err := DataQueryBuilder.GetQueryBuilderList(query, limit, page)
	if err != nil {
		return nil, err
	}

	// query = query.Preload("data")

	return query, nil
}
