package domain

import "time"

// Document represents a document stored in DynamoDB
type Document struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type" validate:"required"`
	Data      map[string]interface{} `json:"data" validate:"required"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// DocumentCreate represents document creation request
type DocumentCreate struct {
	Type string                 `json:"type" validate:"required"`
	Data map[string]interface{} `json:"data" validate:"required"`
}

// DocumentUpdate represents document update request
type DocumentUpdate struct {
	Data map[string]interface{} `json:"data" validate:"required"`
}
