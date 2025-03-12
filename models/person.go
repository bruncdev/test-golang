package models

import (
	"time"
)

type Person struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Phone     string          `json:"phone"`
	Stack     string          `json:"stack"`
	CreatedAt time.Time       `json:"created_at"`
}

type CreatePersonRequest struct {
	Name  string          `json:"name"`
	Email string          `json:"email"`
	Phone string          `json:"phone"`
	Stack string          `json:"stack"`
} 