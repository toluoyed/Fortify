-- Create enums if they don't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role') THEN
        CREATE TYPE role AS ENUM ('USER', 'SUPERUSER');
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN
        CREATE TYPE status AS ENUM ('COMPLETE', 'INCOMPLETE');
    END IF;
END$$;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role role NOT NULL DEFAULT 'USER',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create members table
CREATE TABLE IF NOT EXISTS members (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    phone_number TEXT,
    cohort TEXT NOT NULL,
    year INTEGER NOT NULL,
    session1 BOOLEAN DEFAULT FALSE,
    session2 BOOLEAN DEFAULT FALSE,
    session3 BOOLEAN DEFAULT FALSE,
    session4 BOOLEAN DEFAULT FALSE,
    session1_completion_time TIMESTAMPTZ,
    session2_completion_time TIMESTAMPTZ,
    session3_completion_time TIMESTAMPTZ,
    session4_completion_time TIMESTAMPTZ,
    status status NOT NULL DEFAULT 'INCOMPLETE',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE usersADD COLUMN first_name TEXT NOT NULL,ADD COLUMN last_name TEXT NOT NULL;
