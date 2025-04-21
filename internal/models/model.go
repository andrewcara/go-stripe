package models

import "database/sql"

type DBmodel struct {
	DB *sql.DB
}

type Models struct {
	DB        DBmodel
	artist_id string `json:"artist_id"`
}
