package models

import (
	"time"
)

// ---------- DB Model (matches table) ----------
type User struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Username        string     `json:"username"`
	CenterCode      string     `json:"center_code"`
	CreatedBy       uint       `json:"created_by"`
	DoctorID        uint       `json:"doctor_id"`
	DoctorCode      int        `json:"doctor_code"`
	ClinicAddress   string     `json:"clinic_address"`
	Image           string     `json:"image"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Phone           string     `json:"phone"`
	Flag            int        `json:"flag"` // 1=>Superadmin, 2=>Admin, 3=>doctor
	EmailVerifiedAt *time.Time `json:"-"`
	Password        string     `json:"-"` // hide
	RememberToken   string     `json:"-"`
	Status          int        `json:"status"` // 1=active, 0=inactive
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ---------- Response Structs ----------
type UserBasicResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type DoctorResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	DoctorCode    int    `json:"doctor_code"`
	ClinicAddress string `json:"clinic_address"`
	Phone         string `json:"phone"`
	Image         string `json:"image"`
}

// ---------- Mapping Functions ----------
func ToUserBasicResponse(u User) UserBasicResponse {
	return UserBasicResponse{
		ID:       u.ID,
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
	}
}

func ToDoctorResponse(u User) DoctorResponse {
	return DoctorResponse{
		ID:            u.ID,
		Name:          u.Name,
		DoctorCode:    u.DoctorCode,
		ClinicAddress: u.ClinicAddress,
		Phone:         u.Phone,
		Image:         u.Image,
	}
}
