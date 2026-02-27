# Reading List Web Application

This project is a containerized Go web application with an API and a PostgreSQL database backend. It allows users to manage a reading list of books.

## Features
- Add, view books
- Web frontend and REST API
- PostgreSQL database
- Docker Compose setup for easy deployment

## Prerequisites
- Docker Desktop (switched to linux containers)

## Getting Started

### 1. Clone the repository
> If you want to clone the repository, you need Git installed. Alternatively, you can download the ZIP from GitHub.
```sh
git clone <your-repo-url>
cd pluralsight-webservices-and-applications
```

### 2. Create your `.env` file
Copy the contents of `sample.env` below and create a `.env` file in the project root. Update the values as needed:

```
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=readinglist
DB_HOST=postgres_db
DB_PORT=5432
READINGLIST_API_ENDPOINT=http://api:4000/v1/books
READINGLIST_DB_DSN=postgres://your_db_user:your_db_password@postgres_db:5432/readinglist?sslmode=disable
```

> **Note:** You can use any username and password for the database. The default setup creates the user and database automatically. If you change the username/password, update both the `.env` file and the `init-db/init.sql` script accordingly.

### 3. Build and run the application

To build and start all services:
```sh
docker-compose up --build
```

To run in detached mode:
```sh
docker-compose up -d --build
```

To stop and remove containers and volumes (for a clean reset):
```sh
docker-compose down -v
```

### 4. Access the application
- Web frontend: [http://localhost:8080](http://localhost:8080)
- API: [http://localhost:4000/v1/books](http://localhost:4000/v1/books)

## Customizing Database Credentials
If you want to use a different database username or password:
- Update `DB_USER` and `DB_PASSWORD` in your `.env` file
- Update the user/password in `init-db/init.sql` to match

Example in `init-db/init.sql`:
```
CREATE ROLE your_db_user LOGIN PASSWORD 'your_db_password';
```

## Troubleshooting
- If you see errors about missing tables or columns, run `docker-compose down -v` to reset the database.
- Check logs with `docker-compose logs books_management_api` or `docker-compose logs books_management_web_app`.

## License
MIT
