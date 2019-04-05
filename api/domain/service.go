package domain

import (
	"github.com/fatih/structs"
	"github.com/temesxgn/redeam/api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type service struct {
	repository Repository
}

func (s *service) FindAll(filters bson.M, options *options.FindOptions) (Books, *BookApiError) {
	blogs, err := s.repository.FindAll(filters, options)
	if err != nil {
		return nil, err
	}

	return blogs, nil
}

func (s *service) FindOne(id string) (Book, *BookApiError) {
	return s.repository.FindOne(id)
}

func (s *service) Create(book Book) (string, *BookApiError) {
	doesExist := s.repository.IsExistingEntry(book)

	if doesExist {
		return "", NewAlreadyExistsError()
	}

	if validationError := book.Validate(); validationError != nil {
		return "", NewValidationError(validationError.Error())
	}

	id, persistError := s.repository.Save(book)
	if persistError != nil {
		return "", NewPersistError(persistError.Error())
	}

	return id, nil
}

func (s *service) Update(id string, book Book) *BookApiError {

	if err := book.Validate(); err != nil {
		return NewValidationError(err.Error())
	}

	if _, findError := s.repository.FindOne(id); findError != nil {
		return NewNotFoundError(id)
	}

	mapper := utils.ModelMapper{}
	fields := mapper.ToMongoDocument(structs.Fields(book))
	return s.repository.Update(id, fields)
}

func (s *service) Delete(id string) *BookApiError {

	if _, err := s.repository.FindOne(id); err != nil {
		return NewNotFoundError(id)
	}

	return s.repository.Delete(id)
}

func (s *service) CheckOut(id string) *BookApiError {
	book, err := s.repository.FindOne(id)
	if err != nil {
		return NewNotFoundError(id)
	}

	if book.Status == CheckedOut {
		return NewAlreadyCheckedOutError(id)
	}

	fields := bson.D{
		{"$set", bson.D{{"status", CheckedOut}}},
	}

	if updateError := s.repository.Update(id, fields); updateError != nil {
		return NewUpdateError(updateError.Error())
	}

	return nil
}

func (s *service) CheckIn(id string) *BookApiError {
	book, err := s.repository.FindOne(id)
	if err != nil {
		return NewNotFoundError(id)
	}

	if book.Status == CheckedIn {
		return NewAlreadyCheckedInError(id)
	}

	fields := bson.D{
		{"$set", bson.D{{"status", CheckedIn}}},
	}

	if updateError := s.repository.Update(id, fields); updateError != nil {
		return NewUpdateError(updateError.Error())
	}

	return nil
}

func (s *service) Rate(id string, rate int) *BookApiError {
	book, err := s.repository.FindOne(id)
	if err != nil {
		return NewNotFoundError(id)
	}

	book.Rating = rate
	if validationError := book.Validate(); validationError != nil {
		return NewValidationError(validationError.Error())
	}

	mapper := utils.ModelMapper{}
	fields := mapper.ToMongoDocument(structs.Fields(book))
	if updateError := s.repository.Update(id, fields); updateError != nil {
		return NewUpdateError(updateError.Error())
	}

	return nil
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}
