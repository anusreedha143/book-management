package data

import (
	"testing"

	"readinglist.demo.io/internal/data"
)

func TestGetBook(t *testing.T) { // AAA style - Arrange, Act and Assert
	// Arrange
	expect := Book{
		ID:    1,
		Title: "clean code",
	}
	got, error := Get(expect.ID)
}
