package dto

import "time"

type AuditLogResponse struct {
	ID         string    `json:"id"`
	ActorID    *string   `json:"actor_id,omitempty"`
	TargetUserID *string `json:"target_user_id,omitempty"`
	ActionType string    `json:"action_type"`
	IPAddress  *string   `json:"ip_address,omitempty"`
	UserAgent  *string   `json:"user_agent,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}
