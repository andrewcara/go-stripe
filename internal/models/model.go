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
	ID        int       `json:"id"`
	Artist_id int       `json:"artist_id"`
	Venue     string    `json:"venue"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Artist struct {
	ID          int       `json:"id"`
	Artist_name string    `json:"artist_name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
type Ticket struct {
	ID           int       `json:"id"`
	Event_ID     string    `json:"event_id"`
	Ticket_Price int       `json:"description"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// Need to implement these models in Postgres
type Order struct {
	ID        int       `json:"id"`
	TicketID  int       `json:"ticket_id"`
	StautsID  int       `json:"status_id"`
	Quantity  int       `json:"quantity"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Type for statuses
type Status struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type TransactionStatus struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
type Transaction struct {
	ID                  int       `json:"id"`
	Amount              int       `json:"amount"`
	Currency            int       `json:"currency"`
	LastFour            string    `json:"last_four"`
	BankReturnCode      string    `json:"bank_return_code"`
	TransactionStatusID int       `json:"transaction_status_id"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`
}
type User struct {
	ID        int       `json:"id"`
	FirstName int       `json:"first_name"`
	LastName  int       `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (m *DBmodel) GetEvent(id int) (TicketEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var event TicketEvent

	row := m.DB.QueryRowContext(ctx, "select * from ticket_event where event_id = $1", id)

	err := row.Scan(&event.ID, &event.Artist_id, &event.Venue, &event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		return event, err
	}
	return event, nil
}
