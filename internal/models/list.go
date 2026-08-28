package models

import (
	"time"

	uuid "github.com/google/uuid"
)

type List struct {
	ID          uuid.UUID `db:"id" json:"id,omitempty"`
	Title       string    `db:"title" json:"title,omitempty"`
	Description string    `db:"description" json:"description,omitempty"`
	Price       float64   `db:"price" json:"price,omitempty"`
	City        string    `db:"city" json:"city,omitempty"`
	Status      string    `db:"status" json:"status,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at,omitempty"`
}
