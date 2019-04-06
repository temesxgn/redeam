package domain

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/temesxgn/redeam/api/utils"
	"net/http"
	"strconv"
)

// A Controller - action handler for Book API
type Controller struct {
	service Service
}

// GetAll handles REST API Get '/' Endpoint
func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	responseBuilder := utils.ResponseBuilder{}
	queryBuilder := utils.QueryBuilder{}
	filters, queries := queryBuilder.GetQueryParams(r)
	blogs, err := c.service.FindAll(filters, queries)
	if err != nil {
		responseBuilder.InternalServerError(w, err.Error())
		return
	}

	data, _ := json.Marshal(blogs)
	responseBuilder.OK(w, data)
	return
}

// GetByID handles REST API Get '/{id}' Endpoint
func (c *Controller) GetByID(w http.ResponseWriter, r *http.Request) {
	responseBuilder := utils.ResponseBuilder{}
	id := chi.URLParam(r, "id")

	blog, err := c.service.FindOne(id)
	if err != nil {
		switch err.errorType {
		case NotFoundError:
			responseBuilder.NotFound(w)
			return
		default:
			responseBuilder.InternalServerError(w, err.Error())
			return
		}
	}

	data, _ := json.Marshal(blog)
	responseBuilder.OK(w, data)
}

// Create handles REST API POST '/' Endpoint
func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	responseBuilder := utils.ResponseBuilder{}

	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	if err := book.Validate(); err != nil {
		responseBuilder.BadRequest(w, err.Error())
		return
	}

	id, createError := c.service.Create(book)
	if createError != nil {
		switch createError.errorType {
		case ExistingRecord:
			responseBuilder.BadRequest(w, createError.Error())
			return
		default:
			responseBuilder.InternalServerError(w, createError.Error())
			return
		}
	}

	responseBuilder.OK(w, []byte(id))
}

// Update handles REST API PUT '/{id}' Endpoint
func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	responseBuilder := utils.ResponseBuilder{}
	id := chi.URLParam(r, "id")

	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	if err := book.Validate(); err != nil {
		responseBuilder.BadRequest(w, err.Error())
		return
	}

	if updateError := c.service.Update(id, book); updateError != nil {
		switch updateError.errorType {
		case NotFoundError:
			responseBuilder.NotFound(w)
			return
		default:
			responseBuilder.InternalServerError(w, updateError.Error())
			return
		}
	}

	responseBuilder.OK(w, []byte(""))
}

// Delete handles REST API DELETE '/{id}' Endpoint
func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	responseBuilder := utils.ResponseBuilder{}
	id := chi.URLParam(r, "id")

	if len(id) == 0 {
		responseBuilder.BadRequest(w, "Missing ID!")
		return
	}

	if updateError := c.service.Delete(id); updateError != nil {
		switch updateError.errorType {
		case NotFoundError:
			responseBuilder.NotFound(w)
			return
		default:
			responseBuilder.InternalServerError(w, updateError.Error())
			return
		}
	}

	responseBuilder.OK(w, []byte(""))
}

// CheckOut handles REST API PUT '/checkout/{id}' Endpoint
func (c *Controller) CheckOut(w http.ResponseWriter, r *http.Request) {
	responseBuilder := utils.ResponseBuilder{}
	id := chi.URLParam(r, "id")

	if len(id) == 0 {
		responseBuilder.BadRequest(w, "Missing ID!")
		return
	}

	if err := c.service.CheckOut(id); err != nil {
		switch err.errorType {
		case NotFoundError:
			responseBuilder.NotFound(w)
			return
		case AlreadyCheckedOut:
			responseBuilder.BadRequest(w, err.Error())
			return
		default:
			responseBuilder.InternalServerError(w, err.Error())
			return
		}
	}

	responseBuilder.OK(w, []byte(nil))
}

// CheckIn handles REST API PUT '/checkin/{id}' Endpoint
func (c *Controller) CheckIn(w http.ResponseWriter, r *http.Request) {
	responseBuilder := utils.ResponseBuilder{}
	id := chi.URLParam(r, "id")

	if len(id) == 0 {
		responseBuilder.BadRequest(w, "Missing ID!")
		return
	}

	if err := c.service.CheckIn(id); err != nil {
		switch err.errorType {
		case NotFoundError:
			responseBuilder.NotFound(w)
			return
		case AlreadyCheckedIn:
			responseBuilder.BadRequest(w, err.Error())
			return
		default:
			responseBuilder.InternalServerError(w, err.Error())
			return
		}
	}

	responseBuilder.OK(w, []byte(nil))
}

// Rate handles REST API PUT '/{id}/rate/{id}' Endpoint
func (c *Controller) Rate(w http.ResponseWriter, r *http.Request) {
	responseBuilder := utils.ResponseBuilder{}
	id := chi.URLParam(r, "id")
	rateParam := chi.URLParam(r, "rate")

	rate, rateError := strconv.Atoi(rateParam)
	if rateError != nil {
		responseBuilder.BadRequest(w, "Rate must be a number!")
		return
	}

	if err := c.service.Rate(id, rate); err != nil {
		switch err.errorType {
		case ValidationError:
			responseBuilder.BadRequest(w, err.Error())
			return
		case NotFoundError:
			responseBuilder.NotFound(w)
			return
		default:
			responseBuilder.InternalServerError(w, err.Error())
			return
		}
	}

	responseBuilder.OK(w, []byte(nil))
}

// NewController Creates Controller instance
func NewController(service Service) *Controller {
	return &Controller{service: service}
}
