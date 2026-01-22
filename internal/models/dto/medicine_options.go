package dto

type MedicineCategoryResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	IsActive bool   `json:"is_active"`
}

type MedicineCategoryItemResponse struct {
	ID                string  `json:"id"`
	CategoryID        string  `json:"category_id"`
	DisplayName       string  `json:"display_name"`
	DefaultDosageText *string `json:"default_dosage_text,omitempty"`
	IsActive          bool    `json:"is_active"`
}
