package database

import (
	"database/sql"
	"fmt"
)

type DBConfig struct {
	Server   string
	Port     int
	User     string
	Password string
	Database string
}

func NewConnection(config DBConfig) (*sql.DB, error) {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		config.Server, config.User, config.Password, config.Port, config.Database)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return db, nil
}
