package driver

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func OpenDB(dsn string) (*sql.DB, error) {
	connStr := dsn
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return db, err

}
