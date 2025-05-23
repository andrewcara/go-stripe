package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"log"
	"time"
)

const (
	ScopeAuthentication = "authentication"
)

type Token struct {
	PlainText string    `json:"token"`
	UserID    int64     `json:"-"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

// Generates a token that lasts for ttl and returns it
func GenerateToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: int64(userID),
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)

	if err != nil {
		return nil, err
	}

	//save to database

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256(([]byte(token.PlainText)))

	token.Hash = hash[:]

	return token, nil
}

func (m *DBmodel) InsertToken(t *Token, u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//Delete existing token if exists for user id

	delete_stmt := `delete from tokens where user_id = $1`

	_, err := m.DB.ExecContext(ctx, delete_stmt,
		u.ID)

	if err != nil {
		return err
	}

	stmt := `INSERT INTO tokens
	(user_id, name, email, token_hash)
		VALUES($1, $2, $3, $4)`

	_, err = m.DB.ExecContext(ctx, stmt,
		u.ID,
		u.FirstName,
		u.Email,
		t.Hash)

	if err != nil {
		return err
	}

	return nil

}

func (m *DBmodel) GetUserForToken(token string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tokenHash := sha256.Sum256([]byte(token))

	var user User

	query := `select u.id, u.first_name, u.last_name, u.email
	from users u
	inner join tokens t on  (u.id = t.user_id)
	where t.token_hash = $1`

	err := m.DB.QueryRowContext(ctx, query, tokenHash[:]).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}
