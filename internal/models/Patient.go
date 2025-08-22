package models

import "time"

type Patient struct {
	ID           uint      `gorm:"primaryKey"`
	CreatedBy    string    `gorm:"column:created_by"`
	DoctorCode   string    `gorm:"column:doctor_code"`
	PatientCode  string    `gorm:"column:patient_code"`
	Name         string    `gorm:"column:name"`
	Age          string    `gorm:"column:age"`
	Phone        string    `gorm:"column:phone"`
	Gender       string    `gorm:"column:gender"`
	City         string    `gorm:"column:city"`
	Address      string    `gorm:"column:address"`
	RegisterDate string    `gorm:"column:register_date"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	LanguageID   int       `gorm:"column:language_id"`
}

func (Patient) TableName() string {
	return "patients"
}
