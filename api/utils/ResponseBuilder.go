package utils

import (
	"fmt"
	"log"
	"net/http"
)

// ResponseBuilder - helper methods to generate HTTP responses
type ResponseBuilder struct{}

func (r *ResponseBuilder) InternalServerError(w http.ResponseWriter, msg string) {
	log.Println("Internal error:", msg)
	w.WriteHeader(http.StatusInternalServerError)
	r.build(w, []byte(fmt.Sprintf("Internal Error: %s", msg)))
}

func (r *ResponseBuilder) OK(w http.ResponseWriter, data []byte) {
	w.WriteHeader(http.StatusOK)
	r.build(w, data)
}

func (r *ResponseBuilder) NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	r.build(w, []byte("Not Found"))
}

func (r *ResponseBuilder) BadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	r.build(w, []byte(msg))
}

func (r *ResponseBuilder) build(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}
