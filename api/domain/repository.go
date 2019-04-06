package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

type repo struct{}

var db *mongo.Collection

// FindAll Queries MongoDB with optional filters on custom attributes and Collection options
// It returns a list of paginated Books or an API Error Response
func (r *repo) FindAll(filters bson.M, findOptions *options.FindOptions) (Books, *BookAPIError) {
	cur, colErr := db.Find(nil, filters, findOptions)
	if colErr != nil {
		return nil, NewDatabaseOperationError(colErr.Error())
	}

	var books = Books{}
	for cur.Next(nil) {
		book := Book{}
		decodeErr := cur.Decode(&book)
		if decodeErr != nil {
			return nil, NewDatabaseOperationError(decodeErr.Error())
		}

		books = append(books, book)
	}

	if curErr := cur.Err(); curErr != nil {
		return nil, NewDatabaseOperationError(curErr.Error())
	}

	return books, nil
}

// FineOne Queries MongoDB for a specific Book
// It returns one Book or an API Error Response
func (r *repo) FindOne(id string) (Book, *BookAPIError) {
	var book = Book{}
	objectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objectID}}
	decodeErr := db.FindOne(nil, filter).Decode(&book)
	if decodeErr == mongo.ErrNoDocuments {
		return book, NewNotFoundError(id)
	}

	return book, nil
}

// Delete Hard deletes the Book with the specified ID
// It returns an API Error Response if failed
func (r *repo) Delete(id string) *BookAPIError {
	objectID, _ := primitive.ObjectIDFromHex(id)
	_, deleteError := db.DeleteOne(nil, bson.D{{"_id", objectID}})
	if deleteError != nil {
		return NewDatabaseOperationError(deleteError.Error())
	}

	return nil
}

// Update updates the Book with the specified ID
// It returns an API Error Response if failed
func (r *repo) Update(id string, updatedFields bson.D) *BookAPIError {
	objectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objectID}}
	_, updateError := db.UpdateOne(nil, filter, updatedFields)
	if updateError != nil {
		return NewDatabaseOperationError(updateError.Error())
	}

	return nil
}

// ExistingEntry Checks the DB if it has a Book with the same unique composite fields
// Author, Title, Publish_Date
// It returns a boolean
func (r *repo) IsExistingEntry(book Book) bool {
	var existingBook Book
	filter := bson.D{{"author", book.Author}, {"title", book.Title}, {"publish_date", book.PublishDate}}
	decodeErr := db.FindOne(nil, filter).Decode(&existingBook)
	if decodeErr == mongo.ErrNoDocuments {
		return false
	}

	return true
}

// Save Saves the Book Payload
// It returns the persisted Book ID or an API Error Response if failed
func (r *repo) Save(book Book) (string, *BookAPIError) {
	created, insertError := db.InsertOne(nil, book)
	if insertError != nil {
		return "", NewPersistError(insertError.Error())
	}

	return created.InsertedID.(primitive.ObjectID).Hex(), nil
}

// NewRepository Initializes a repository instance
// It returns an API Error Response if failed
func NewRepository() (Repository, *BookAPIError) {
	server, serverPresent := os.LookupEnv("mongodb_url")
	if !serverPresent {
		return nil, NewMissingEnvVariable("need to set mongodb_url environment variable")
	}

	dbName, dbPresent := os.LookupEnv("database_name")
	if !dbPresent {
		return nil, NewMissingEnvVariable("need to set database_name environment variable")
	}

	collection, collectionPresent := os.LookupEnv("collection_name")
	if !collectionPresent {
		return nil, NewMissingEnvVariable("need to set collection_name environment variable")
	}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(server))

	if err != nil {
		return nil, NewDatabaseOperationError(err.Error())
	}

	db = client.Database(dbName).Collection(collection)
	return &repo{}, nil
}
