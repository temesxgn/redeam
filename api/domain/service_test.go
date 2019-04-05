package domain

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestService_FindAll(t *testing.T) {
	books := Books{
		{
			Author: "thg090020",
		},
	}
	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindAll(nil, nil).Return(books, nil)
	books, err := bookService.FindAll(nil, nil)

	if err != nil {
		t.Fatalf(`Invalid response.. Expected to get response but got error %s`, err.Error())
	}

	if len(books) != 1 || books[0].Author != "thg090020" {
		t.Fatalf(`Invalid response.. Expected %d but Got %d\n`, 1, len(books))
	}
}

func TestService_FindAll_WithError(t *testing.T) {
	findOptions := &options.FindOptions{}
	findOptions = findOptions.SetSkip(0).SetLimit(10)

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindAll(bson.M{}, findOptions).Return(nil, NewDatabaseOperationError("connection error"))
	_, err := bookService.FindAll(bson.M{}, findOptions)

	if err == nil {
		t.Fatalf(`Invalid response.. Expected to get an error but was nil`)
	}
}

func TestService_FindOne(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(testBook, nil)
	resBook, _ := bookService.FindOne(testBook.ID.Hex())

	assert.Equal(t, resBook.ID.Hex(), testBook.ID.Hex(),
		`Invalid response.. Expected %s but Got %s\n`, testBook.ID.Hex(), resBook.ID.Hex())
}

func TestService_Update_WithValidationError(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)
	bookRepo.EXPECT().Update(testBook.ID.Hex(), testBook).Return(nil)
	err := bookService.Update(testBook.ID.Hex(), testBook)

	assert.NotNil(t, err, `Invalid response.. Expected error but Got %s\n`, err.Error())
}

func TestService_Create(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)
	bookRepo.EXPECT().IsExistingEntry(testBook).Return(false)
	bookRepo.EXPECT().Save(testBook).Return(testBook.ID.Hex(), nil)
	bookId, err := bookService.Create(testBook)

	assert.Nil(t, err, `Invalid response.. Expected error to be nul but Got %s\n`, err)
	assert.Equal(t, testBook.ID.Hex(), bookId,
		`Invalid response.. Expected %s but Got %s\n`, testBook.ID.Hex(), bookId)
}

func TestService_Create_WithExisting(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)
	bookRepo.EXPECT().IsExistingEntry(testBook).Return(true)
	_, err := bookService.Create(testBook)

	assert.Equal(t, err.errorType, ExistingRecord,
		`Invalid response.. Expected existing record error but Got %s\n`, err.errorType)
}

func TestService_Create_WithValidationError(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)
	bookRepo.EXPECT().IsExistingEntry(testBook).Return(false)
	_, err := bookService.Create(testBook)

	assert.Equal(t, err.errorType, ValidationError,
		`Invalid response.. Expected validation error but Got %s\n`, err.errorType.Name())
}

func TestService_Create_WithPersistenceError(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)
	bookRepo.EXPECT().IsExistingEntry(testBook).Return(false)
	bookRepo.EXPECT().Save(testBook).Return("", NewPersistError("Error"))
	_, err := bookService.Create(testBook)

	assert.Equal(t, err.errorType, PersistError,
		`Invalid response.. Expected existing record error but Got %s\n`, err.errorType.Name())
}

func TestService_Delete(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(Book{}, nil)
	bookRepo.EXPECT().Delete(testBook.ID.Hex()).Return(nil)
	err := bookService.Delete(testBook.ID.Hex())

	assert.Nil(t, err, `Invalid response.. Expected no error but Got %s\n`, err)
}

func TestService_Update(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)
	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(Book{}, nil)
	bookRepo.EXPECT().Update(testBook.ID.Hex(), gomock.Any()).Return(nil)
	err := bookService.Update(testBook.ID.Hex(), testBook)

	assert.Nil(t, err, `Invalid response.. Expected error to be nul but Got %s\n`, err)
}

func TestService_UpdateWithNotFoundError(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)
	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(Book{}, NewNotFoundError(testBook.ID.Hex()))
	err := bookService.Update(testBook.ID.Hex(), testBook)

	assert.NotNil(t, err, `Invalid response.. Expected error to be nul but Got %s\n`, err)
}

func TestService_DeleteWithNotFoundError(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(Book{}, NewNotFoundError(testBook.ID.Hex()))
	err := bookService.Delete(testBook.ID.Hex())

	assert.NotNil(t, err, `Invalid response.. Expected an error but Got %s\n`, err.errorType.Name())
}

func TestService_DeleteWithError(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(Book{}, nil)
	bookRepo.EXPECT().Delete(testBook.ID.Hex()).Return(NewDatabaseOperationError("error deleting"))
	err := bookService.Delete(testBook.ID.Hex())

	assert.NotNil(t, err, `Invalid response.. Expected an error but Got %s\n`, err.errorType.Name())
}

func TestService_CheckOut(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
		Status: CheckedIn,
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(testBook, nil)
	bookRepo.EXPECT().Update(testBook.ID.Hex(), gomock.Any()).Return(nil)
	err := bookService.CheckOut(testBook.ID.Hex())

	assert.Nil(t, err, `Invalid response.. Expected an error but Got %s\n`, err)
}

func TestService_CheckOut_WithNotFound(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
		Status: CheckedIn,
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(testBook, NewNotFoundError(testBook.ID.Hex()))
	bookRepo.EXPECT().Update(testBook.ID.Hex(), gomock.Any()).Return(nil)
	err := bookService.CheckOut(testBook.ID.Hex())

	assert.NotNil(t, err, `Invalid response.. Expected an error but Got %s\n`, err)
}

func TestService_CheckOut_WithCheckedOutError(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
		Status: CheckedOut,
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(testBook, nil)
	err := bookService.CheckOut(testBook.ID.Hex())

	assert.NotNil(t, err, `Invalid response.. Expected an error but Got %s\n`, err)
}

func TestService_CheckOut_WithPersistenceError(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
		Status: CheckedIn,
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(testBook, nil)
	bookRepo.EXPECT().Update(testBook.ID.Hex(), gomock.Any()).Return(NewPersistError("error"))
	err := bookService.CheckOut(testBook.ID.Hex())

	assert.NotNil(t, err, `Invalid response.. Expected an error but Got %s\n`, err)
}

func TestService_CheckIn(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
		Status: CheckedOut,
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(testBook, nil)
	bookRepo.EXPECT().Update(testBook.ID.Hex(), gomock.Any()).Return(nil)
	err := bookService.CheckIn(testBook.ID.Hex())

	assert.Nil(t, err, `Invalid response.. Expected an error but Got %s\n`, err)
}

func TestService_CheckIn_WithNotFound(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
		Status: CheckedOut,
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(testBook, NewNotFoundError(testBook.ID.Hex()))
	bookRepo.EXPECT().Update(testBook.ID.Hex(), gomock.Any()).Return(nil)
	err := bookService.CheckIn(testBook.ID.Hex())

	assert.NotNil(t, err, `Invalid response.. Expected an error but Got %s\n`, err)
}

func TestService_CheckIn_WithCheckedOutError(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
		Status: CheckedIn,
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(testBook, nil)
	err := bookService.CheckIn(testBook.ID.Hex())

	assert.NotNil(t, err, `Invalid response.. Expected an error but Got %s\n`, err)
}

func TestService_CheckIn_WithPersistenceError(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
		Status: CheckedOut,
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(testBook, nil)
	bookRepo.EXPECT().Update(testBook.ID.Hex(), gomock.Any()).Return(NewPersistError("error"))
	err := bookService.CheckIn(testBook.ID.Hex())

	assert.NotNil(t, err, `Invalid response.. Expected an error but Got %s\n`, err)
}

func TestService_Rate(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(testBook, nil)
	bookRepo.EXPECT().Update(testBook.ID.Hex(), gomock.Any()).Return(nil)
	err := bookService.Rate(testBook.ID.Hex(), 2)

	assert.Nil(t, err, `Invalid response.. Expected an error but Got %s\n`, err)
}

func TestService_Rate_WithNotFound(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(testBook, NewNotFoundError(testBook.ID.Hex()))
	err := bookService.Rate(testBook.ID.Hex(), 2)

	assert.NotNil(t, err, `Invalid response.. Expected an error but Got %s\n`, err)
}

func TestService_Rate_WithValidationError(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(testBook, nil)
	err := bookService.Rate(testBook.ID.Hex(), 2)

	assert.NotNil(t, err, `Invalid response.. Expected an error but Got %s\n`, err)
}

func TestService_Rate_WithPersistenceError(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      CheckedIn,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	bookRepo := NewMockRepository(gomock.NewController(t))
	bookService := NewService(bookRepo)

	bookRepo.EXPECT().FindOne(testBook.ID.Hex()).Return(testBook, nil)
	bookRepo.EXPECT().Update(testBook.ID.Hex(), gomock.Any()).Return(NewPersistError("error"))
	err := bookService.Rate(testBook.ID.Hex(), 2)

	assert.NotNil(t, err, `Invalid response.. Expected an error but Got %s\n`, err)
}
