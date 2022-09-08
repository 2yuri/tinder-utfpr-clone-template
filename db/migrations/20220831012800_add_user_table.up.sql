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
