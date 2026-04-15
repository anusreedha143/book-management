// package data

// import "database/sql"

// type Models struct {
// 	Books BookModel
// }

// func NewModels(db *sql.DB) Models {
// 	return Models{
// 		Books: BookModel{DB: db},
// 	}
// }

package data

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type Models struct {
	Books BookModel // Assuming your model is called BookModel
}

func NewModels(client *mongo.Client, logger *log.Logger) Models {
	return Models{
		// We tell the BookModel exactly which Database and Collection to use
		Books: BookModel{
			Collection: client.Database("books").Collection("readinglistdb"),
			Logger:     logger,
		},
	}
}
