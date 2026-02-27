-- Create user and grant privileges
DO
$$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'readinglistdbuser') THEN
      CREATE ROLE readinglistdbuser LOGIN PASSWORD 'vikky';
   END IF;
END
$$;

GRANT ALL PRIVILEGES ON DATABASE readinglist TO readinglistdbuser;
GRANT ALL ON SCHEMA public TO readinglistdbuser;

-- Create books table if not exists
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

ALTER TABLE books OWNER TO readinglistdbuser;

-- Create readinglist table if not exists
-- CREATE TABLE IF NOT EXISTS readinglist (
--     id SERIAL PRIMARY KEY,
--     user_id INTEGER,
--     book_id INTEGER REFERENCES books(id),
--     added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- ALTER TABLE readinglist OWNER TO readinglistdbuser;
