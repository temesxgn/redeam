package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const (
	ContentType = "application/json"
)

func TestController_GetAll(t *testing.T) {
	books := Books{
		{
			Author: "thg090020",
		},
	}

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return(books, nil)
	bookController := NewController(bookService)

	wr := httptest.NewRecorder()
	testURL, _ := url.Parse("http://localhost:8080/books")
	r := &http.Request{URL: testURL}
	bookController.GetAll(wr, r)

	assert.Equal(t, http.StatusOK, wr.Code, `Invalid response... Expected 200 but got %d`, wr.Code)
	assert.True(t, bytes.Contains(wr.Body.Bytes(), []byte("thg090020")))
}

func TestController_GetAll_WithError(t *testing.T) {
	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return(nil, NewDatabaseOperationError("error"))
	bookController := NewController(bookService)

	wr := httptest.NewRecorder()
	testURL, _ := url.Parse("http://localhost:8080/books")
	r := &http.Request{URL: testURL}
	bookController.GetAll(wr, r)

	assert.Equal(t, http.StatusInternalServerError, wr.Code, `Invalid response... Expected 200 but got %d`, wr.Code)
	assert.True(t, bytes.Contains(wr.Body.Bytes(), []byte("error")))
}

func TestController_GetByID(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
	}

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().FindOne(testBook.ID.Hex()).Return(testBook, nil)
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Get("/books/{id}", bookController.GetByID)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Get(fmt.Sprintf("%s/books/%s", server.URL, testBook.ID.Hex()))
	defer closeBody(res.Body)

	body, _ := ioutil.ReadAll(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode, `Invalid response... Expected 200 but got %d`, res.StatusCode)
	assert.True(t, bytes.Contains(body, []byte("thg090020")))
}

func TestController_GetByIDWithNotFound(t *testing.T) {
	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().FindOne(gomock.Any()).Return(Book{}, NewNotFoundError("Not found"))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Get("/books/{id}", bookController.GetByID)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Get(fmt.Sprintf("%s/books/%s", server.URL, "1234"))
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusNotFound, res.StatusCode, `Invalid response... Expected 404 but got %d`, res.StatusCode)
}

func TestController_GetByIDWithInternalError(t *testing.T) {
	testBook := Book{
		ID:     primitive.NewObjectID(),
		Author: "thg090020",
	}

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().FindOne(gomock.Any()).Return(Book{}, NewDatabaseOperationError("Internal error"))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Get("/books/{id}", bookController.GetByID)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Get(fmt.Sprintf("%s/books/%s", server.URL, testBook.ID.Hex()))
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode, `Invalid response... Expected 500 but got %d`, res.StatusCode)
}

func TestController_Create(t *testing.T) {
	testBook := Book{
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().Create(gomock.Any()).Return(primitive.NewObjectID().Hex(), nil)
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books", bookController.Create)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books", server.URL), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode, `Invalid response... Expected 200 but got %d`, res.StatusCode)
}

func TestController_Create_WithMalformedBody(t *testing.T) {
	testBook := struct {
		Name string
	}{}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().Create(gomock.Any()).Return(primitive.NewObjectID().Hex(), nil)
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books", bookController.Create)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books", server.URL), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode, `Invalid response... Expected 400 but got %d`, res.StatusCode)
}

func TestController_Create_WithExistingRecord(t *testing.T) {
	testBook := Book{
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().Create(gomock.Any()).Return("", NewAlreadyExistsError())
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books", bookController.Create)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books", server.URL), ContentType, requestReader)
	defer closeBody(res.Body)

	body, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, `Invalid response... Expected 400 but got %d`, res.StatusCode)
	assert.True(t, bytes.Contains(body, []byte(NewAlreadyExistsError().msg)))
}

func TestController_Create_WithInternalError(t *testing.T) {
	testBook := Book{
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().Create(gomock.Any()).Return("", NewDatabaseOperationError("internal error"))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books", bookController.Create)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books", server.URL), ContentType, requestReader)
	defer closeBody(res.Body)

	body, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode, `Invalid response... Expected 500 but got %d`, res.StatusCode)
	assert.True(t, bytes.Contains(body, []byte("internal error")))
}

func TestController_Update(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().Update(testBook.ID.Hex(), testBook).Return(nil)
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}", bookController.Update)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s", server.URL, testBook.ID.Hex()), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode, `Invalid response... Expected 200 but got %d`, res.StatusCode)
}

func TestController_Update_WithMalformedBody(t *testing.T) {
	testBook := struct {
		Name string
	}{}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}", bookController.Update)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s", server.URL, gomock.Any()), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode, `Invalid response... Expected 400 but got %d`, res.StatusCode)
}

func TestController_Update_WithNotFound(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(NewNotFoundError("not found"))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}", bookController.Update)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s", server.URL, gomock.Any()), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusNotFound, res.StatusCode, `Invalid response... Expected 404 but got %d`, res.StatusCode)
}

func TestController_Update_WithInternalError(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(NewDatabaseOperationError("internal error"))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}", bookController.Update)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s", server.URL, gomock.Any()), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode, `Invalid response... Expected 500 but got %d`, res.StatusCode)
}

func TestController_Delete(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().Delete(testBook.ID.Hex()).Return(nil)
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}", bookController.Delete)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s", server.URL, testBook.ID.Hex()), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode, `Invalid response... Expected 200 but got %d`, res.StatusCode)
}

func TestController_Delete_WithNotFound(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().Delete(gomock.Any()).Return(NewNotFoundError("not found"))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}", bookController.Delete)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s", server.URL, gomock.Any()), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusNotFound, res.StatusCode, `Invalid response... Expected 404 but got %d`, res.StatusCode)
}

func TestController_Delete_WithInternalError(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().Delete(gomock.Any()).Return(NewDatabaseOperationError("internal error"))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}", bookController.Delete)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s", server.URL, gomock.Any()), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode, `Invalid response... Expected 500 but got %d`, res.StatusCode)
}

func TestController_CheckOut(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().CheckOut(testBook.ID.Hex()).Return(nil)
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}", bookController.CheckOut)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s", server.URL, testBook.ID.Hex()), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode, `Invalid response... Expected 200 but got %d`, res.StatusCode)
}

func TestController_CheckOut_WithNotFound(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().CheckOut(gomock.Any()).Return(NewNotFoundError("not found"))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}", bookController.CheckOut)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s", server.URL, gomock.Any()), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusNotFound, res.StatusCode, `Invalid response... Expected 404 but got %d`, res.StatusCode)
}

func TestController_CheckOut_WithAlreadyCheckedOut(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().CheckOut(gomock.Any()).Return(NewAlreadyCheckedOutError("internal error"))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}", bookController.CheckOut)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s", server.URL, gomock.Any()), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode, `Invalid response... Expected 400 but got %d`, res.StatusCode)
}

func TestController_CheckOut_WithInternalError(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().CheckOut(gomock.Any()).Return(NewDatabaseOperationError("internal error"))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}", bookController.CheckOut)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s", server.URL, gomock.Any()), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode, `Invalid response... Expected 500 but got %d`, res.StatusCode)
}

func TestController_CheckIn(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().CheckIn(testBook.ID.Hex()).Return(nil)
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}", bookController.CheckIn)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s", server.URL, testBook.ID.Hex()), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode, `Invalid response... Expected 200 but got %d`, res.StatusCode)
}

func TestController_CheckIn_WithNotFound(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().CheckIn(gomock.Any()).Return(NewNotFoundError("not found"))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}", bookController.CheckIn)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s", server.URL, gomock.Any()), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusNotFound, res.StatusCode, `Invalid response... Expected 404 but got %d`, res.StatusCode)
}

func TestController_CheckIn_WithAlreadyCheckedOut(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().CheckIn(gomock.Any()).Return(NewAlreadyCheckedInError("internal error"))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}", bookController.CheckIn)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s", server.URL, gomock.Any()), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode, `Invalid response... Expected 400 but got %d`, res.StatusCode)
}

func TestController_CheckIn_WithInternalError(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().CheckIn(gomock.Any()).Return(NewDatabaseOperationError("internal error"))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}", bookController.CheckIn)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s", server.URL, gomock.Any()), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode, `Invalid response... Expected 500 but got %d`, res.StatusCode)
}

func closeBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.Println("Error closing connection!")
	}
}

func TestController_Rate(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().Rate(testBook.ID.Hex(), gomock.Any()).Return(nil)
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}/rate/{rate}", bookController.Rate)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s/rate/%d", server.URL, testBook.ID.Hex(), 2), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode, `Invalid response... Expected 200 but got %d`, res.StatusCode)
}

func TestController_Rate_WithNotFound(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().Rate(testBook.ID.Hex(), gomock.Any()).Return(NewNotFoundError("not found"))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}/rate/{rate}", bookController.Rate)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s/rate/%d", server.URL, testBook.ID.Hex(), 2), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusNotFound, res.StatusCode, `Invalid response... Expected 404 but got %d`, res.StatusCode)
}

func TestController_Rate_WithInternalError(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().Rate(testBook.ID.Hex(), gomock.Any()).Return(NewDatabaseOperationError("not found"))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}/rate/{rate}", bookController.Rate)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s/rate/%d", server.URL, testBook.ID.Hex(), 2), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode, `Invalid response... Expected 500 but got %d`, res.StatusCode)
}

func TestController_Rate_WithInvalidRate(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookService.EXPECT().Rate(testBook.ID.Hex(), gomock.Any()).Return(NewValidationError("not found"))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}/rate/{rate}", bookController.Rate)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s/rate/%d", server.URL, testBook.ID.Hex(), 300), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode, `Invalid response... Expected 404 but got %d`, res.StatusCode)
}

func TestController_Rate_WithInvalidRateNumber(t *testing.T) {
	testBook := Book{
		ID:          primitive.NewObjectID(),
		Author:      "thg090020",
		Title:       "Test Title",
		Status:      2,
		Rating:      1,
		Publisher:   "Pub",
		PublishDate: "2019",
	}

	requestBytes, _ := json.Marshal(testBook)
	requestReader := bytes.NewReader(requestBytes)

	bookService := NewMockService(gomock.NewController(t))
	bookController := NewController(bookService)
	router := chi.NewRouter()
	router.Post("/books/{id}/rate/{rate}", bookController.Rate)

	server := httptest.NewServer(router)
	defer server.Close()

	res, _ := http.Post(fmt.Sprintf("%s/books/%s/rate/%s", server.URL, testBook.ID.Hex(), "abc"), ContentType, requestReader)
	defer closeBody(res.Body)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode, `Invalid response... Expected 400 but got %d`, res.StatusCode)
}
