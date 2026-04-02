package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	books, err := app.readinglist.GetAll()
	if err != nil {
		// FIX: Tell us WHICH part failed
		http.Error(w, fmt.Sprintf("API Error: %v", err), 500)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		// FIX: Tell us if the files are missing
		log.Printf("Template Error: %v", err)
		http.Error(w, fmt.Sprintf("Template Parsing Error: %v", err), 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", books)
	if err != nil {
		log.Printf("Template Execution Error: %v", err)
		http.Error(w, fmt.Sprintf("Template Execution Error: %v", err), 500)
		return
	}
	// fmt.Fprintf(w, "<html><head><title>Reading List</title></head><body><h1>Reading List</h1><ul>")
	// for _, book := range *books {
	// 	fmt.Fprintf(w, "<li>%s (%d)</li>", book.Title, book.Pages)
	// }
	// fmt.Fprintf(w, "</ul></body></html>")
}

func (app *application) bookView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		log.Printf("ERROR: Invalid ID provided: %s", r.URL.Query().Get("id"))
		http.NotFound(w, r)
		return
	}

	book, err := app.readinglist.Get(int64(id))
	if err != nil {
		// LOUD LOG: Tell us if the DB connection died or the ID doesn't exist
		log.Printf("ERROR: Database Get(%d) failed: %v", id, err)
		http.Error(w, fmt.Sprintf("Database Error: %v", err), 500)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/view.html",
	}

	// Used to convert comma-separated genres to a slice within the template.
	funcs := template.FuncMap{"join": strings.Join}

	ts, err := template.New("showBook").Funcs(funcs).ParseFiles(files...)
	if err != nil {
		log.Printf("ERROR: Template Parsing failed for ID %d: %v", id, err)
		http.Error(w, fmt.Sprintf("Template File Missing: %v", err), 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", book)
	if err != nil {
		log.Printf("ERROR: Template Execution failed for ID %d: %v", id, err)
		http.Error(w, fmt.Sprintf("Template Execution Error: %v", err), 500)
		return
	}
	//fmt.Fprintf(w, "%s (%d)\n", book.Title, book.Pages)
}

func (app *application) bookCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.bookCreateForm(w, r)
	case http.MethodPost:
		app.bookCreateProcess(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) bookCreateForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/create.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *application) bookCreateProcess(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	if title == "" {
		http.Error(w, "Title is Empty", http.StatusBadRequest)
		return
	}

	pages, err := strconv.Atoi(r.PostFormValue("pages"))
	if err != nil || pages < 1 {
		http.Error(w, "The value of the pages is empty", http.StatusBadRequest)
		return
	}

	published, err := strconv.Atoi(r.PostFormValue("published"))
	if err != nil || published < 1 {
		http.Error(w, "The value of the pages is empty", http.StatusBadRequest)
		return
	}

	// Handles multiple spaces or commas in genres better
	genreInput := r.PostFormValue("genres")
	genreInput = strings.ReplaceAll(genreInput, ",", " ")
	genres := strings.Fields(genreInput)

	ratingFloat, err := strconv.ParseFloat(r.PostFormValue("rating"), 32)
	if err != nil {
		http.Error(w, "Invalid Rating", http.StatusBadRequest)
		return
	}

	book := struct {
		Title     string   `json:"title"`
		Pages     int      `json:"pages"`
		Published int      `json:"published"`
		Genres    []string `json:"genres"`
		Rating    float32  `json:"rating"`
	}{
		Title:     r.PostFormValue("title"),
		Pages:     pages,
		Published: published,
		Genres:    genres,
		Rating:    float32(ratingFloat),
	}

	data, err := json.Marshal(book)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", app.readinglist.Endpoint, bytes.NewBuffer(data))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create request: %v", err), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("API Connection Error: %v", err), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	//Capture the API's actual error message
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("API ERROR (%s): %s", resp.Status, string(body))

		// Show the actual API error in the browser for easier debugging
		http.Error(w, fmt.Sprintf("API rejected create (%s): %s", resp.Status, string(body)), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) bookEdit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.bookEditForm(w, r)
	case http.MethodPost:
		app.bookEditProcess(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) bookEditForm(w http.ResponseWriter, r *http.Request) {
	// 1. Get ID from URL
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	book, err := app.readinglist.Get(int64(id))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// 3. Render the Template
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/edit.html",
	}

	funcs := template.FuncMap{"join": strings.Join}

	ts, err := template.New("edit.html").Funcs(funcs).ParseFiles(files...)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// 4. Create a buffer to hold the rendered template
	buf := new(bytes.Buffer)

	err = ts.ExecuteTemplate(buf, "base", book)
	if err != nil {
		log.Printf("TEMPLATE ERROR: %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// If we got here, rendering was successful. Now send the buffer to the browser.
	buf.WriteTo(w)
}

func (app *application) bookEditProcess(w http.ResponseWriter, r *http.Request) {
	// 1. Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// 2. Validate and extract the form values
	idStr := r.Form.Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	title := r.Form.Get("title")
	pages, err := strconv.Atoi(r.Form.Get("pages"))
	if err != nil || pages < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	published, err := strconv.Atoi(r.Form.Get("published"))
	if err != nil || published < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Split by comma instead of space
	rawGenres := r.Form.Get("genres")
	genres := strings.Split(rawGenres, ",")

	for i := range genres {
		genres[i] = strings.TrimSpace(genres[i]) // Remove leading/trailing spaces
	}

	var filtered []string
	for _, g := range genres {
		if g != "" {
			filtered = append(filtered, g)
		}
	}
	genres = filtered

	rating, err := strconv.ParseFloat(r.Form.Get("rating"), 32)
	if err != nil || rating < 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// 3. Create the book struct
	book := struct {
		Title     string   `json:"title"`
		Pages     int      `json:"pages"`
		Published int      `json:"published"`
		Genres    []string `json:"genres"`
		Rating    float32  `json:"rating"`
	}{
		Title:     title,
		Pages:     pages,
		Published: published,
		Genres:    genres,
		Rating:    float32(rating),
	}

	// 4. Send the updated book data to the API
	data, err := json.Marshal(book)
	if err != nil {
		http.Error(w, "Problem in JSON marshaling", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("PUT", app.readinglist.Endpoint+"/"+idStr, bytes.NewBuffer(data))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status: %s", resp.Status)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// 5. Redirect back to the book details page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) bookDelete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		app.bookDeleteProcess(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) bookDeleteProcess(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	apiURL := app.readinglist.Endpoint + "/" + idStr

	// Send DELETE request to API
	req, err := http.NewRequest(http.MethodDelete, apiURL, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Printf("Unexpected status: %s", resp.Status)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// 4. Redirect back to home after successful deletion
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
