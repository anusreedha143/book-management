package data

import (
	"context"
	"errors"
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	// MongoDB uses ObjectIDs. we useprimitive.ObjectID from the mongo-driver/bson/primitive package to represent these IDs in Go.
	// The 'bson' struct tag tells the MongoDB driver how to map this field to the '_id' field in the database.
	// The 'json' struct tag allows us to control how this field is represented in JSON when we send responses to clients.
	// The 'omitempty' option in the 'bson' tag means that if the ID is empty (zero value), it will not be included in the BSON document sent to MongoDB.
	// This is important because MongoDB will automatically generate an ObjectID for new documents if the '_id' field is not provided. By using 'omitempty', we can allow MongoDB to handle ID generation without having to set it manually in our Go code before insertion.
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	Title     string             `bson:"title" json:"title"`
	Published int                `bson:"published,omitempty" json:"published,omitempty"`
	Pages     int                `bson:"pages,omitempty" json:"pages,omitempty"`
	Genres    []string           `bson:"genres,omitempty" json:"genres,omitempty"`
	Rating    float64            `bson:"rating,omitempty" json:"rating,omitempty"`
	Version   int32              `bson:"version" json:"-"`
}

type BookModel struct {
	Collection *mongo.Collection
	Logger     *log.Logger
}

// Insert a new book into the database
func (m BookModel) Insert(book *Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	book.CreatedAt = time.Now()
	book.Version = 1
	book.Rating = math.Round(book.Rating*10) / 10

	result, err := m.Collection.InsertOne(ctx, book)
	if err != nil {
		m.Logger.Printf("ERROR: Database InsertOne failed: %v", err)
		return err
	}

	// Assign the generated ID back to the book object
	book.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// Get a book by ID
func (m BookModel) Get(id string) (*Book, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		m.Logger.Printf("WARN: Invalid Hex ID provided: %s", objID)
		return nil, errors.New("record not found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var book Book
	err = m.Collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&book)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("record not found")
		}
		m.Logger.Printf("ERROR: Database FindOne failed for ID %s: %v", id, err)
		return nil, err
	}

	return &book, nil
}

func (m BookModel) Update(book *Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.M{"_id": book.ID, "version": book.Version}
	update := bson.M{
		"$set": bson.M{
			"title": book.Title, "published": book.Published,
			"pages": book.Pages, "genres": book.Genres, "rating": book.Rating,
		},
		"$inc": bson.M{"version": 1},
	}

	result, err := m.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		m.Logger.Printf("ERROR: Database UpdateOne failed for ID %v: %v", book.ID, err)
		return err
	}

	if result.MatchedCount == 0 {
		m.Logger.Printf("WARN: Update conflict or record missing for ID %v (Version %d)", book.ID, book.Version)
		return errors.New("edit conflict or record not found")
	}

	book.Version++
	return nil
}

func (m BookModel) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("record not found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.Collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		m.Logger.Printf("ERROR: Database DeleteOne failed for ID %s: %v", id, err)
		return err
	}

	if result.DeletedCount == 0 {
		m.Logger.Printf("WARN: Attempted to delete non-existent ID %s", id)
		return errors.New("record not found")
	}

	return nil
}

// GetAll retrieves all books from the database
func (m BookModel) GetAll() ([]*Book, error) {
	// 1. Create a timeout context (Prevents the API from hanging if DB is slow)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 2. Define Sorting (Like your SQL "ORDER BY id")
	// 1 = Ascending, -1 = Descending
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	// 3. Execute the Find query
	// bson.D{} is an empty filter, meaning "Find Everything"
	cursor, err := m.Collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		m.Logger.Printf("ERROR: Database Find failed: %v", err)
		return nil, err
	}
	// Always ensure the cursor is closed to prevent memory leaks
	defer cursor.Close(ctx)

	// 4. Decode all results into the slice at once
	var books []*Book
	if err = cursor.All(ctx, &books); err != nil {
		m.Logger.Printf("ERROR: Failed to decode books from cursor: %v", err)
		return nil, err
	}

	// 5. Check for cursor errors after iteration
	if err = cursor.Err(); err != nil {
		m.Logger.Printf("ERROR: Cursor iteration error: %v", err)
		return nil, err
	}

	return books, nil
}

// type Book struct {
// 	ID        int64     `json:"id"`
// 	CreatedAt time.Time `json:"-"`
// 	Title     string    `json:"title"`
// 	Published int       `json:"published,omitempty"`
// 	Pages     int       `json:"pages,omitempty"`
// 	Genres    []string  `json:"genres,omitempty"`
// 	Rating    float32   `json:"rating,omitempty"`
// 	Version   int32     `json:"-"`
// }

// type BookModel struct {
// 	DB *sql.DB
// }

// // Go treats $1 as "just data", never as "executable code". Using placeholders like $1, $2 etc(Postgres) in your queries, you are already safe from SQL injection
// func (b BookModel) Insert(book *Book) error {
// 	query := `
// 		INSERT INTO books (title, published, pages, genres, rating)
// 		VALUES ($1, $2, $3, $4, $5)
// 		RETURNING id, created_at, version`

// 	args := []interface{}{book.Title, book.Published, book.Pages, pq.Array(book.Genres), book.Rating}
// 	// return the auto generated system values to Go object
// 	return b.DB.QueryRow(query, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)
// }

// func (b BookModel) Get(id int64) (*Book, error) {
// 	if id < 1 {
// 		return nil, errors.New("record not found")
// 	}

// 	query := `
// 		SELECT id, created_at, title, published, pages, genres, rating, version
// 		FROM books
// 		WHERE id = $1`

// 	var book Book

// 	err := b.DB.QueryRow(query, id).Scan(
// 		&book.ID,
// 		&book.CreatedAt,
// 		&book.Title,
// 		&book.Published,
// 		&book.Pages,
// 		pq.Array(&book.Genres),
// 		&book.Rating,
// 		&book.Version,
// 	)

// 	if err != nil {
// 		switch {
// 		case errors.Is(err, sql.ErrNoRows):
// 			return nil, errors.New("record not found")
// 		default:
// 			return nil, err
// 		}
// 	}

// 	return &book, nil
// }

// func (b BookModel) Update(book *Book) error {
// 	query := `
// 		UPDATE books
// 		SET title = $1, published = $2, pages = $3, genres = $4, rating = $5, version = version + 1
// 		WHERE id = $6 AND version = $7
// 		RETURNING version`

// 	args := []interface{}{book.Title, book.Published, book.Pages, pq.Array(book.Genres), book.Rating, book.ID, book.Version}
// 	return b.DB.QueryRow(query, args...).Scan(&book.Version)
// }

// func (b BookModel) Delete(id int64) error {
// 	if id < 1 {
// 		return errors.New("record not found")
// 	}

// 	query := `
// 		DELETE FROM books
// 		WHERE id = $1`

// 	results, err := b.DB.Exec(query, id)
// 	if err != nil {
// 		return err
// 	}

// 	rowsAffected, err := results.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	if rowsAffected == 0 {
// 		return errors.New("record not found")
// 	}

// 	return nil
// }

// func (b BookModel) GetAll() ([]*Book, error) {
// 	query := `
// 	  SELECT id, created_at, title, published, pages, genres, rating, version
// 	  FROM public.books
// 	  ORDER BY id`

// 	rows, err := b.DB.Query(query)
// 	if err != nil {
// 		log.Printf("ERROR: Database Query failed: %v", err)
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	books := []*Book{}

// 	for rows.Next() {
// 		var book Book

// 		err := rows.Scan(
// 			&book.ID,
// 			&book.CreatedAt,
// 			&book.Title,
// 			&book.Published,
// 			&book.Pages,
// 			pq.Array(&book.Genres),
// 			&book.Rating,
// 			&book.Version,
// 		)
// 		if err != nil {
// 			log.Printf("ERROR: GetAll Scan failed: %v", err)
// 			return nil, err
// 		}

// 		books = append(books, &book)
// 	}

// 	if err = rows.Err(); err != nil {
// 		log.Printf("ERROR: GetAll Rows iteration error: %v", err)
// 		return nil, err
// 	}

// 	return books, nil
// }
