#!/bin/bash
set -e

# This "EOQ" block allows us to write SQL inside a bash script
# We can inject variables like $POSTGRES_USER directly here
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    
    -- 1. Create the table
    CREATE TABLE IF NOT EXISTS books (
        id bigserial PRIMARY KEY,
        created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
        title text NOT NULL,
        published integer NOT NULL,
        pages integer NOT NULL,
        genres text[] NOT NULL,
        rating real NOT NULL,
        version integer NOT NULL DEFAULT 1
    );

    -- 2. Alter Owner using the Environment Variable
    -- Since we are logging in AS this user to run the script, 
    -- they usually own it by default, but this enforces it explicitly.
    ALTER TABLE books OWNER TO "$POSTGRES_USER";

    -- 3. Insert Data
    INSERT INTO books (title, published, pages, genres, rating) 
    VALUES ('The Go Programming Language', 2015, 380, '{Education, Programming}', 4.9)
    ON CONFLICT DO NOTHING;

EOSQL