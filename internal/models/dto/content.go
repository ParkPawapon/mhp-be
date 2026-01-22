package dto

import "time"

type HealthContentResponse struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	BodyContent      *string   `json:"body_content,omitempty"`
	ThumbnailURL     *string   `json:"thumbnail_url,omitempty"`
	ExternalVideoURL *string   `json:"external_video_url,omitempty"`
	Category         *string   `json:"category,omitempty"`
	IsPublished      bool      `json:"is_published"`
	CreatedAt        time.Time `json:"created_at"`
}

type CreateHealthContentRequest struct {
	Title            string  `json:"title" validate:"required"`
	BodyContent      *string `json:"body_content"`
	ThumbnailURL     *string `json:"thumbnail_url"`
	ExternalVideoURL *string `json:"external_video_url"`
	Category         *string `json:"category"`
}

type UpdateHealthContentRequest struct {
	Title            *string `json:"title"`
	BodyContent      *string `json:"body_content"`
	ThumbnailURL     *string `json:"thumbnail_url"`
	ExternalVideoURL *string `json:"external_video_url"`
	Category         *string `json:"category"`
}

type PublishHealthContentRequest struct {
	IsPublished bool `json:"is_published"`
}
