package main

import (
	"log"
	"tinderutf/db"
)
func main(){
	db.StartDB()
	_, err := db.Db.Exec(`
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE sex AS ENUM (
    'MALE', 'FEMALE', 'ALL'
    );

CREATE TABLE IF NOT EXISTS users (
    id SERIAL,
    name TEXT NOT NULL,
    birth_date TIMESTAMP NOT NULL,
    email TEXT UNIQUE NOT NULL,
    instagram TEXT DEFAULT '',
    password TEXT NOT NULL,
    about TEXT DEFAULT '',
    is_active BOOLEAN DEFAULT true,
    sex sex NOT NULL,
    sex_preference sex NOT NULL DEFAULT 'ALL',
    find_distance INTEGER NOT NULL DEFAULT 10,
    geolocation geometry,
    external_id UUID NOT NULL DEFAULT uuid_generate_v1(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS interactions (
    id SERIAL,
    user_id INTEGER NOT NULL,
    target_user_id INTEGER NOT NULL,
    liked bool NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,
    FOREIGN KEY (target_user_id) REFERENCES users (id) ON DELETE RESTRICT,
    UNIQUE(user_id, target_user_id)
);

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON interactions
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();`)
if err != nil {
	log.Fatal("cannot run sql")
}

log.Println("migrations created!")
}