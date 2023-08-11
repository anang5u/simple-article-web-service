package service

import (
	"database/sql"
	"fmt"
	"log"
	"simple-ddd-cqrs/config"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

// GetDBConnection
func GetDBConnection() *sql.DB {
	if db != nil {
		return db
	}

	return initDBConnection()
}

// initDBConnection
func initDBConnection() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Get("DB_HOST"),
		config.Get("DB_PORT"),
		config.Get("DB_USER"),
		config.Get("DB_PASSWORD"),
		config.Get("DB_NAME"),
	)

	// Initialise a new connection
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	// connection pool configuration
	maxOpenConns, _ := strconv.Atoi(config.Get("DB_POOL_MAX_OPEN_CONNS"))
	maxIdleConns, _ := strconv.Atoi(config.Get("DB_POOL_MAX_IDLE_CONNS"))
	maxLifetime, _ := strconv.Atoi(config.Get("DB_POOL_MAX_LIFETIME"))

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Minute)

	return db
}
