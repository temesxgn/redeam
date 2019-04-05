package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/temesxgn/redeam/api"
	"log"
	"net/http"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		middleware.Logger,    // Log API request calls
		middleware.Recoverer, // Recover from panics without crashing server
		middleware.SetHeader("Content-Type", "application/json"), // Set content-Type headers as application/json
	)

	apiRoutes, err := api.Routes()
	if err != nil {
		log.Println("Error initializing route:", err.Error())
	}

	router.Route("/", func(r chi.Router) {
		r.Mount("/books", apiRoutes)
	})

	return router
}

func main() {
	log.Println("Starting Apps...")
	router := Routes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("Walking %s %s\n", method, route) // Walk and print out all routes
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	// Note, the port is usually gotten from the environment.
	_ = http.ListenAndServe(":8080", router)
}
