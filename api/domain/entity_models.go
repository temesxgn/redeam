package domain

import (
	"github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/** ======== Book Model ========*/
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
func (b Book) Validate() *BookApiError {
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
	} else {
		return nil
	}
}

/** ======== Status Model ========*/
type Status int

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

/** ======== Query Criterion Model ========*/
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
	FindAll(filters bson.M, findOptions *options.FindOptions) (Books, *BookApiError)
	FindOne(id string) (Book, *BookApiError)
	Update(id string, fields bson.D) *BookApiError
	Delete(id string) *BookApiError
	IsExistingEntry(book Book) bool
	Save(book Book) (string, *BookApiError)
}

/** ======== Service Interface ========*/
type Service interface {
	FindAll(filters bson.M, findOptions *options.FindOptions) (Books, *BookApiError)
	FindOne(id string) (Book, *BookApiError)
	Update(id string, book Book) *BookApiError
	Delete(id string) *BookApiError
	CheckOut(id string) *BookApiError
	CheckIn(id string) *BookApiError
	Create(book Book) (string, *BookApiError)
	Rate(id string, rate int) *BookApiError
}
