package database

import (
	"database/sql"
	"fmt"
	"time"

	"data-access/internal/config"

	"github.com/go-sql-driver/mysql"
)

type Database struct {
	*sql.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := mysql.Config{
		User:                 cfg.DBUser,
		Passwd:               cfg.DBPassword,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		DBName:               cfg.DBName,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return &Database{db}, nil
}

func (db *Database) Close() error {
	return db.DB.Close()
}
