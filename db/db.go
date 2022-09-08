package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
	"tinderutf/utils/env"
)

var Db *sql.DB

func StartDB() {
	host := env.GetEnv("DB_HOST", "localhost")
	port := env.GetEnv("DB_PORT", "15432")
	nameDB := env.GetEnv("DB_NAME", "tinder")
	user := env.GetEnv("DB_USER", "tinder")
	passwd := env.GetEnv("DB_PASS", "123")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, nameDB, passwd)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error on open connection: ", err)
	}

	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxIdleConns(10)
	if err := db.Ping(); err != nil {
		log.Fatalf("Could not connect to the DB: %v", err)
	}

	Db = db
}

func Close() {
	if Db != nil {
		Db.Close()
	}
}
