package querybuilder

import (
	"gorm.io/gorm"
)

type BaseQueryBuilderList interface {
	GetQueryBuilderList(query *gorm.DB, limit int, page int) (*gorm.DB, error)
}

type BaseQueryBuilderListImpl struct {
	db *gorm.DB
}

func NewBaseQueryBuilderList(db *gorm.DB) *BaseQueryBuilderListImpl {
	return &BaseQueryBuilderListImpl{
		db: db,
	}
}

func (baseQueryBuilder *BaseQueryBuilderListImpl) GetQueryBuilderList(query *gorm.DB, limit int, page int) (*gorm.DB, error) {
	// Set limit dan pagination
	if limit == 0 {
		limit = 100
	}
	query = query.Limit(limit)

	if page == 0 {
		page = 1
	}
	query = query.Offset((page - 1) * limit)

	return query, nil
}