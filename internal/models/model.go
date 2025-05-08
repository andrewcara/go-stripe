package models

import (
	"context"
	"database/sql"
	"strings"
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
type Event struct {
	ID          int       `json:"id"`
	Artist_ID   int       `json:"artist_id"`
	Venue       string    `json:"venue"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
type TicketEvent struct {
	ID          int    `json:"id"`
	Venue       string `json:"venue"`
	Price       int    `json:"ticket_price"`
	Description string `json:"description,omitempty"`
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
	Ticket_Price int       `json:"ticket_price"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// Need to implement these models in Postgres
type Order struct {
	ID            int       `json:"id"`
	TicketID      int       `json:"ticket_id"`
	StatusID      int       `json:"status_id"`
	TransactionID int       `json:"transaction_id"`
	CustomerID    int       `json:"customer_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
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
	Currency            string    `json:"currency"`
	LastFour            string    `json:"last_four"`
	BankReturnCode      string    `json:"bank_return_code"`
	TransactionStatusID int       `json:"transaction_status_id"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`
	ExpiryMonth         int       `json:"expiry_month"`
	ExpiryYear          int       `json:"expiry_year"`
	PaymentIntent       string    `json:"payment_intent"`
	PaymentMethod       string    `json:"payment_method"`
}
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
type Customer struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (m *DBmodel) GetEvent(id int) (Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var event Event

	row := m.DB.QueryRowContext(ctx, "select * from events where event_id = $1", id)

	err := row.Scan(&event.ID, &event.Artist_ID, &event.Venue, &event.CreatedAt, &event.UpdatedAt, &event.Description)

	if err != nil {
		return event, err
	}
	return event, nil
}

func (m *DBmodel) GetTicketEvents(id int) (TicketEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var event TicketEvent

	row := m.DB.QueryRowContext(ctx,
		`select 
		events.event_id, venue, ticket_price, description 
		from events 
			inner join 
		tickets on events.event_id = tickets.event_id
			where events.event_id = $1`, id)

	err := row.Scan(&event.ID, &event.Venue, &event.Price, &event.Description)

	if err != nil {
		return event, err
	}
	return event, nil
}

// inserting transaction
func (m *DBmodel) InsertTransaction(txn Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO transactions
	(amount, currency, last_four, bank_return_code,
		transaction_status_id, expiry_month, expiry_year, payment_intent, payment_method)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	var id int
	err := m.DB.QueryRowContext(ctx, stmt,
		txn.Amount,
		txn.Currency,
		txn.LastFour,
		txn.BankReturnCode,
		txn.TransactionStatusID,
		txn.ExpiryMonth,
		txn.ExpiryYear,
		txn.PaymentIntent,
		txn.PaymentMethod).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *DBmodel) InsertOrder(order Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO orders
	(ticket_id, status_id, quantity, amount, transaction_id, customer_id)
		VALUES($1, $2, $3, $4, $5, $6) RETURNING id`

	var id int
	err := m.DB.QueryRowContext(ctx, stmt,
		order.TicketID,
		order.StatusID,
		order.Quantity,
		order.Amount,
		order.TransactionID,
		order.CustomerID).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
func (m *DBmodel) InsertCustomer(c Customer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO customers
	(first_name, last_name, email)
		VALUES($1, $2, $3) RETURNING id`

	var id int
	err := m.DB.QueryRowContext(ctx, stmt,
		c.FirstName,
		c.LastName,
		c.Email).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Gets a user by email address
func (m *DBmodel) GetUserByEmail(email string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	email = strings.ToLower(email)
	var u User
	row := m.DB.QueryRowContext(ctx, "select id, first_name, last_name, email, password, updated_at, created_at  from users where email = $1", email)

	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.UpdatedAt, &u.CreatedAt)

	if err != nil {
		return u, err
	}
	return u, nil
}
