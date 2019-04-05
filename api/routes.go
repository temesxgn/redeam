package api

import (
	"github.com/go-chi/chi"
	"github.com/temesxgn/redeam/api/domain"
)

// Enabled Routes for /books path
func Routes() (*chi.Mux, *domain.BookApiError) {
	router := chi.NewRouter()
	repo, err := domain.NewRepository()
	service := domain.NewService(repo)
	ctrl := domain.NewController(service)
	router.Get("/", ctrl.GetAll)
	router.Post("/", ctrl.Create)

	router.Get("/{id}", ctrl.GetById)
	router.Put("/{id}", ctrl.Update)
	router.Delete("/{id}", ctrl.Delete)

	router.Put("/checkout/{id}", ctrl.CheckOut)
	router.Put("/checkin/{id}", ctrl.CheckIn)

	router.Put("/{id}/rate/{rate}", ctrl.Rate)

	return router, err
}
