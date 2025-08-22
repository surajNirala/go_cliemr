package services

import (
	"github.com/surajNirala/go_cliemr/internal/repository"
)

type ExcelService struct {
	repo *repository.ExcelRepository
}

func NewExcelService(repo *repository.ExcelRepository) *ExcelService {
	return &ExcelService{repo: repo}
}
