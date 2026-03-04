package main

import (
	"bytes"
	"encoding/json"
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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", books)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", 500)
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
		http.NotFound(w, r)
		return
	}

	book, err := app.readinglist.Get(int64(id))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", book)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
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

// func (app *application) bookCreateForm(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "<html><head><title>Create Book</title></head>"+
// 		"<body><h1>Create Book</h1><form action=\"/book/create\" method=\"post\">"+
// 		"<label for=\"title\">Title</label><input type=\"text\" name=\"title\" id=\"title\">"+
// 		"<label for=\"pages\">Pages</label><input type=\"number\" name=\"pages\" id=\"pages\">"+
// 		"<label for=\"published\">Published</label><input type=\"number\" name=\"published\" id=\"published\">"+
// 		"<label for=\"genres\">Genres</label><input type=\"text\" name=\"genres\" id=\"genres\">"+
// 		"<label for=\"rating\">Rating</label><input type=\"number\" step=\"0.1\" name=\"rating\" id=\"rating\">"+
// 		"<button type=\"submit\">Create</button></form></body></html>")
// }

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

	genres := strings.Split(r.PostFormValue("genres"), " ")

	ratingFloat, err := strconv.ParseFloat(r.PostFormValue("rating"), 32)
	if err != nil {
		http.Error(w, "Rating is empty", http.StatusBadRequest)
		return
	}
	rating := float32(ratingFloat)

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
		Rating:    rating,
	}

	data, err := json.Marshal(book)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", app.readinglist.Endpoint, bytes.NewBuffer(data))
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

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Unexpected status: %s", resp.Status)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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

	// 4. Pass the fetched book to the template
	err = ts.ExecuteTemplate(w, "base", book)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
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
	// This handles "Fantasy, Sci-Fi" and "Fantasy,Sci-Fi" correctly
	genres := strings.Split(rawGenres, ",")

	for i := range genres {
		genres[i] = strings.TrimSpace(genres[i]) // Remove leading/trailing spaces
	}
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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
