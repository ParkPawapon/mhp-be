package dto

import "time"

type SupportChatRequestCreateRequest struct {
	Message       string  `json:"message" validate:"required"`
	Category      string  `json:"category" validate:"required"`
	AttachmentURL *string `json:"attachment_url"`
}

type SupportChatRequestResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type SupportChatRequestItem struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Message       string    `json:"message"`
	Category      string    `json:"category"`
	AttachmentURL *string   `json:"attachment_url,omitempty"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type SupportEmergencyResponse struct {
	Hotline     string `json:"hotline"`
	DisplayName string `json:"display_name"`
}
