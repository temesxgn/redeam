package domain

import (
	"github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Book slice of books
type Books []Book

// Book model for Book schema
type Book struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Author      string             `bson:"author" json:"author"`
	Title       string             `bson:"title" json:"title"`
	Publisher   string             `bson:"publisher" json:"publisher"`
	Status      Status             `bson:"status" json:"status"`
	Rating      int                `bson:"rating" json:"rating"`
	PublishDate string             `bson:"publish_date" json:"publish_date"`
}

// Validate validates the Book fields.
func (b Book) Validate() *BookAPIError {
	errors := validation.ValidateStruct(&b,
		validation.Field(&b.Author, validation.Required, validation.Length(1, 30)),
		validation.Field(&b.Title, validation.Required, validation.Length(1, 50)),
		validation.Field(&b.Publisher, validation.Required, validation.Length(1, 20)),
		validation.Field(&b.Status, validation.Required, validation.In(CheckedIn, CheckedOut)),
		validation.Field(&b.Rating, validation.In(0, 1, 2, 3)),
		validation.Field(&b.PublishDate, validation.Required, validation.Date("2006")),
	)

	if errors != nil {
		return NewValidationError(errors.Error())
	}

	return nil
}

// Status
type Status int

// Status options
const (
	Unknown Status = iota
	CheckedIn
	CheckedOut
)

func (status Status) String() string {
	names := [...]string{
		"CheckedIn",
		"CheckedOut",
	}

	// prevent panicking in case of
	// `status` is out of range
	if status < CheckedIn || status > CheckedOut {
		return "Unknown"
	}

	return names[status]
}

type SortOrder int

const (
	ASC  SortOrder = 1
	DESC           = -1
)

// QueryOperator possible query operations
type QueryOperator string

const (
	Equals             QueryOperator = "="
	LessThan                         = "<"
	GreaterThan                      = ">"
	GreaterThanOrEqual               = GreaterThan + Equals
	LessThanOrEqual                  = LessThan + Equals
	DoesNotEqual                     = "!="
)

type Query struct {
	Field    string
	Operator QueryOperator
	Value    string
}

/** ======== Repository Interface ========*/
type Repository interface {
	FindAll(filters bson.M, findOptions *options.FindOptions) (Books, *BookAPIError)
	FindOne(id string) (Book, *BookAPIError)
	Update(id string, fields bson.D) *BookAPIError
	Delete(id string) *BookAPIError
	IsExistingEntry(book Book) bool
	Save(book Book) (string, *BookAPIError)
}

/** ======== Service Interface ========*/
type Service interface {
	FindAll(filters bson.M, findOptions *options.FindOptions) (Books, *BookAPIError)
	FindOne(id string) (Book, *BookAPIError)
	Update(id string, book Book) *BookAPIError
	Delete(id string) *BookAPIError
	CheckOut(id string) *BookAPIError
	CheckIn(id string) *BookAPIError
	Create(book Book) (string, *BookAPIError)
	Rate(id string, rate int) *BookAPIError
}
