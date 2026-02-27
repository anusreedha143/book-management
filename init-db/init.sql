-- -- Create user and grant privileges
-- DO
-- $$
-- BEGIN
--    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'readinglistdbuser') THEN
--       CREATE ROLE readinglistdbuser LOGIN PASSWORD 'vikky';
--    END IF;
-- END
-- $$;

-- GRANT ALL PRIVILEGES ON DATABASE readinglist TO readinglistdbuser;
-- GRANT ALL ON SCHEMA public TO readinglistdbuser;

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

-- 2. Ensure your DB user owns the table 
-- (This ensures your Go app can write to it without permission errors)
ALTER TABLE books OWNER TO readinglistdbuser;

-- 3. Optional: Add a sample book to verify persistence on first run
INSERT INTO books (title, published, pages, genres, rating) 
VALUES ('The Go Programming Language', 2015, 380, '{Education, Programming}', 4.9)
ON CONFLICT DO NOTHING;

