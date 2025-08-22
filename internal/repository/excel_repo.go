package repository

import (
	"gorm.io/gorm"
)

type ExcelRepository struct {
	db *gorm.DB
}

func NewExcelRepository(db *gorm.DB) *ExcelRepository {
	return &ExcelRepository{db: db}
}
