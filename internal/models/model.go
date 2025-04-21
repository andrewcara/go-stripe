package models

import (
	"database/sql"
	"time"
)

type DBmodel struct {
	DB *sql.DB
}

// wrapper for models
type Models struct {
	DB DBmodel
}

// returns a model type with a db connection pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBmodel{DB: db},
	}
}

// Event is the event for all events
type TicketEvent struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
