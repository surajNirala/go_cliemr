package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/surajNirala/go_cliemr/internal/models"
	"github.com/surajNirala/go_cliemr/internal/services"
	"github.com/surajNirala/go_cliemr/pkg/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ExcelHandler struct {
	service *services.ExcelService
	db      *gorm.DB
}

func NewExcelHandler(service *services.ExcelService, db *gorm.DB) *ExcelHandler {
	return &ExcelHandler{service: service, db: db}
}

/* func (s *ExcelHandler) PatientImport(c *gin.Context) {
	utils.RespondSuccess(c, "success", http.StatusOK, "Patient import endpoint is not implemented yet. This is a placeholder response.", nil)
} */

func (s *ExcelHandler) PatientImport(c *gin.Context) {
	// Max file size (5MB)
	const maxFileSize = 5 * 1024 * 1024
	file, err := c.FormFile("file")
	if err != nil {
		utils.RespondError(c, "error", http.StatusBadRequest, "File is required")
		return
	}
	doctorCode := c.PostForm("doctor_code")
	if doctorCode == "" {
		utils.RespondError(c, "error", http.StatusBadRequest, "doctor_code is required")
		return
	}
	// Validate size
	if file.Size > maxFileSize {
		utils.RespondError(c, "error", http.StatusBadRequest, "File too large. Max 5MB allowed")
		return
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		utils.RespondError(c, "error", http.StatusBadRequest, "Invalid file type. Only Excel (.xlsx, .xls) allowed")
		return
	}

	// Save temp file
	tempPath := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		utils.RespondError(c, "error", http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Open Excel
	f, err := excelize.OpenFile(tempPath)
	if err != nil {
		utils.RespondError(c, "error", http.StatusInternalServerError, "Failed to read Excel file")
		return
	}
	defer f.Close()

	// Get first sheet
	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		utils.RespondError(c, "error", http.StatusInternalServerError, "Failed to read rows")
		return
	}

	// Build patients slice
	var patients []models.Patient
	var patientCodes, phones []string

	createdBy := c.GetString("userID")

	for i, row := range rows {
		if i == 0 {
			// skip header
			continue
		}

		code := safeCell(row, 0)
		phone := safeCell(row, 3)

		if code != "" {
			patientCodes = append(patientCodes, code)
		}
		if phone != "" {
			phones = append(phones, phone)
		}

		p := models.Patient{
			CreatedBy:    createdBy,
			DoctorCode:   doctorCode,
			PatientCode:  code,
			Name:         safeCell(row, 1),
			Age:          safeCell(row, 2),
			Phone:        phone,
			Gender:       safeCell(row, 4),
			City:         safeCell(row, 5),
			Address:      safeCell(row, 6),
			RegisterDate: parseDateSafe(row, 7),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			LanguageID:   1,
		}
		patients = append(patients, p)
	}

	// 🔹 Check duplicates in DB
	var existing []models.Patient
	if len(patientCodes) > 0 || len(phones) > 0 {
		query := s.db.Model(&models.Patient{})
		if len(patientCodes) > 0 {
			query = query.Or("patient_code IN ?", patientCodes)
		}
		// if len(phones) > 0 {
		// 	query = query.Or("phone IN ?", phones)
		// }
		if err := query.Find(&existing).Error; err != nil {
			utils.RespondError(c, "error", http.StatusInternalServerError, "Failed to check duplicates")
			return
		}
	}

	// Build sets of existing codes/phones
	existingCodes := make(map[string]bool)
	existingPhones := make(map[string]bool)
	for _, e := range existing {
		if e.PatientCode != "" {
			existingCodes[strings.ToLower(e.PatientCode)] = true
		}
		if e.Phone != "" {
			existingPhones[strings.ToLower(e.Phone)] = true
		}
	}

	// Filter patients (remove duplicates)
	var newPatients []models.Patient
	var skipped []models.Patient
	// Track patient_codes inside the same Excel file
	seenCodes := make(map[string]bool)
	for _, p := range patients {
		code := strings.ToLower(p.PatientCode)
		if existingCodes[strings.ToLower(p.PatientCode)] || existingPhones[strings.ToLower(p.Phone)] {
			skipped = append(skipped, p) // duplicate found
			continue
		}
		// Skip if patient_code already appeared in the same sheet
		if code != "" && seenCodes[code] {
			skipped = append(skipped, p)
			continue
		}
		// Mark as seen and add to newPatients
		seenCodes[code] = true
		newPatients = append(newPatients, p)
	}

	// Bulk Insert only new patients
	if len(newPatients) > 0 {
		if err := s.db.CreateInBatches(&newPatients, 100).Error; err != nil {
			utils.RespondError(c, "error", http.StatusInternalServerError, "Failed to insert patients", err.Error())
			return
		}
	}

	utils.RespondSuccess(c, "success", http.StatusOK, "Patient import completed", gin.H{
		"inserted": len(newPatients),
		"skipped":  len(skipped),
	})
}

// Helpers
func safeCell(row []string, idx int) string {
	if idx < len(row) {
		return strings.TrimSpace(strings.ToLower(row[idx]))
	}
	return ""
}

func parseDateSafe(row []string, idx int) string {
	if idx < len(row) && row[idx] != "" {
		// Try multiple date formats
		formats := []string{"2006-01-02", "02/01/2006", "01/02/2006"}
		for _, f := range formats {
			if t, err := time.Parse(f, row[idx]); err == nil {
				return t.Format("2006-01-02")
			}
		}
	}
	return ""
}
