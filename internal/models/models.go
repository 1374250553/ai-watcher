package models

import "time"

type News struct {
	ID        int64     `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	URL       string    `json:"url" db:"url"`
	Summary   string    `json:"summary" db:"summary"`
	Source    string    `json:"source" db:"source"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type APIResource struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Provider    string    `json:"provider" db:"provider"`
	Description string    `json:"description" db:"description"`
	Endpoint    string    `json:"endpoint" db:"endpoint"`
	FreeQuota   string    `json:"free_quota" db:"free_quota"`
	DocURL      string    `json:"doc_url" db:"doc_url"`
	Model       string    `json:"model" db:"model"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	LastUpdated time.Time `json:"last_updated" db:"last_updated"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type FetchLog struct {
	ID        int64     `json:"id" db:"id"`
	Source    string    `json:"source" db:"source"`
	Status    string    `json:"status" db:"status"`
	Message   string    `json:"message" db:"message"`
	Count     int       `json:"count" db:"count"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
