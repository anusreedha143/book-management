package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"readinglist.demo.io/internal/data"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	// Fprintf allows you to use %s for formatting
	// fmt.Fprintln(w, "Status: available")
	// fmt.Fprintf(w, "environment: %s\n", app.config.env)
	// fmt.Fprintf(w, "version: %s\n", version)
	// js, err := json.Marshal(data)
	// if err != nil {
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 	return
	// }

	// js = append(js, '\n')

	// w.Header().Set("Content-Type", "application/json")
	// w.Write(js)
	if err := app.writeJson(w, http.StatusOK, envelope{"config": data}, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) getCreateBooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		books, err := app.models.Books.GetAll()
		if err != nil {
			log.Printf("GET /books: failed to fetch books: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := app.writeJson(w, http.StatusOK, envelope{"books": books}, nil); err != nil {
			log.Printf("Error parsing GetBooks response: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	if r.Method == http.MethodPost {
		var input struct {
			Title     string   `json:"title"`
			Published int      `json:"published"`
			Pages     int      `json:"pages"`
			Genres    []string `json:"genres"`
			Rating    float32  `json:"rating"`
		}

		err := app.readJSON(w, r, &input)
		if err != nil {
			log.Printf("%s %s from %s: invalid JSON: %v", r.Method, r.URL.Path, r.RemoteAddr, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		book := &data.Book{
			Title:     input.Title,
			Published: input.Published,
			Pages:     input.Pages,
			Genres:    input.Genres,
			Rating:    input.Rating,
		}

		err = app.models.Books.Insert(book)
		if err != nil {
			log.Printf("POST /books: failed to insert book: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("v1/books/%d", book.ID))

		// Write the JSON response with a 201 Created status code and the Location header set.
		err = app.writeJson(w, http.StatusCreated, envelope{"book": book}, headers)
		if err != nil {
			log.Printf("POST /books: failed to insert book (title=%s): %v", book.ID, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (app *application) getUpdateDeleteBooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			app.getBook(w, r)
		}
	case http.MethodPut:
		{
			app.updateBook(w, r)
		}
	case http.MethodDelete:
		{
			app.deleteBook(w, r)
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getBook(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/books/")
	log.Printf("DEBUG: API attempting to parse ID string: '%s'", id)

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("ERROR: Could not parse ID from path %q: %v", r.URL.Path, err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	book, err := app.models.Books.Get(idInt)
	if err != nil {
		switch {
		case errors.Is(err, errors.New("record not found")):
			log.Printf("%s %s: book not found (id=%d)", r.Method, r.URL.Path, idInt)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			log.Printf("ERROR: DB Get failed for ID %d: %v", id, err)
			// Check if it's actually a "not found" vs a DB error
			http.Error(w, "Book not found or database error", 500)
			return
		}
		return
	}

	if err := app.writeJson(w, http.StatusOK, envelope{"book": book}, nil); err != nil {
		log.Printf("%s %s: failed to write JSON response (id=%d): %v", r.Method, r.URL.Path, idInt, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) updateBook(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/books/")
	log.Printf("DEBUG: API attempting to parse ID string: '%s'", id)

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("ERROR: UpdateBook failed to parse ID string '%s': %v", id, err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Fetch the current book record from the database.
	book, err := app.models.Books.Get(idInt)
	if err != nil {
		log.Printf("ERROR: UpdateBook failed to find book ID %d: %v", idInt, err)
		switch {
		case errors.Is(err, errors.New("record not found")):
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// 2. Define input struct
	var input struct {
		Title     *string  `json:"title"`
		Published *int     `json:"published"`
		Pages     *int     `json:"pages"`
		Genres    []string `json:"genres"`
		Rating    *float32 `json:"rating"`
	}

	// Read JSON
	err = app.readJSON(w, r, &input)
	if err != nil {
		log.Printf("ERROR: UpdateBook failed to read JSON for ID %d: %v", idInt, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// Update the book object only where values were provided
	if input.Title != nil {
		book.Title = *input.Title
	}

	if input.Published != nil {
		book.Published = *input.Published
	}

	if input.Pages != nil {
		book.Pages = *input.Pages
	}

	if len(input.Genres) > 0 {
		book.Genres = input.Genres
	}

	if input.Rating != nil {
		book.Rating = *input.Rating
	}
	// Commit to the database
	err = app.models.Books.Update(book)
	if err != nil {
		// LOUD LOG: This is where the Postgres "Not Null" error would show up
		log.Printf("CRITICAL: UpdateBook failed DB update for ID %d: %v", idInt, err)
		http.Error(w, fmt.Sprintf("Database Update Error: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("SUCCESS: Updated book ID %d: %s", idInt, book.Title)

	if err := app.writeJson(w, http.StatusOK, envelope{"book": book}, nil); err != nil {
		log.Printf("ERROR: UpdateBook failed to write response JSON: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {
	// 1. Extract and Log the ID
	id := strings.TrimPrefix(r.URL.Path, "/v1/books/")
	log.Printf("DEBUG: DeleteBook - Received request for ID: '%s'", id)

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("ERROR: DeleteBook - Invalid ID format '%s': %v", id, err)
		http.Error(w, "Bad Request: Invalid ID", http.StatusBadRequest)
		return
	}

	// 2. Perform the Deletion
	log.Printf("DEBUG: DeleteBook - Attempting to remove record ID: %d", idInt)
	err = app.models.Books.Delete(idInt)

	if err != nil {
		// We check the error string specifically since errors.New won't match across packages
		if err.Error() == "record not found" {
			log.Printf("INFO: DeleteBook - No record found for ID: %d", idInt)
			http.Error(w, "Record Not Found", http.StatusNotFound)
			return
		}

		// LOUD LOG: For system-level failures (DB connection, etc.)
		log.Printf("CRITICAL: DeleteBook - Database failure for ID %d: %v", idInt, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// 3. Log Success
	log.Printf("SUCCESS: DeleteBook - Record ID %d permanently deleted", idInt)

	err = app.writeJson(w, http.StatusOK, envelope{"message": "book successfully deleted"}, nil)
	if err != nil {
		log.Printf("ERROR: DeleteBook - Failed to write JSON response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
