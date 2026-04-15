package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Book struct {
	// 1. Changed to string to accept the MongoDB Hex ID
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Published int      `json:"published"`
	Pages     int      `json:"pages"`
	Genres    []string `json:"genres"`
	// 2. Changed to float64 to match MongoDB's 'double' precision
	Rating float64 `json:"rating"`
	// 3. Usually, the Web UI doesn't need the Version,
	// but keep it if you are displaying it.
	Version int32 `json:"version"`
}

type BookResponse struct { //test 123
	Book *Book `json:"book"`
}

type BooksResponse struct {
	Books *[]Book `json:"books"`
}

type ReadinglistModel struct {
	Endpoint string
}

func (m *ReadinglistModel) GetAll() (*[]Book, error) {
	resp, err := http.Get(m.Endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var booksResp BooksResponse
	err = json.Unmarshal(data, &booksResp)
	if err != nil {
		return nil, err
	}

	return booksResp.Books, nil
}

func (m *ReadinglistModel) Get(id string) (*Book, error) {
	url := fmt.Sprintf("%s/%s", m.Endpoint, id)
	log.Printf("DEBUG: ReadinglistModel.Get - Fetching book with ID: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var bookResp BookResponse
	err = json.Unmarshal(data, &bookResp)
	if err != nil {
		return nil, err
	}

	log.Printf("DEBUG: Web App fetched book ID %s with Version %d", bookResp.Book.ID, bookResp.Book.Version)

	return bookResp.Book, nil
}
