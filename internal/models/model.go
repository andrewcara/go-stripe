package models

import (
	"context"
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

func (m *DBmodel) GetEvent(id int) (TicketEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var event TicketEvent

	row := m.DB.QueryRowContext(ctx, "select event_id, venue from ticket_event where event_id = $1", id)

	err := row.Scan(&event.ID, &event.Name)

	if err != nil {
		return event, err
	}

	return event, nil
}
