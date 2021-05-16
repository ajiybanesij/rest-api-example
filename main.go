package main

import (
	"log"
	"net/http"
	"rest-api-example/controller"
	"rest-api-example/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// Init Routher
	r := mux.NewRouter()

	// Mock Data @todo - implement DB
	// books = append(books, models.Book{ID: "1", Isbn: "12345", Title: "Book One", Author: &models.Author{Firstname: "John", Lastname: "Doe"}})
	// books = append(books, models.Book{ID: "2", Isbn: "12346", Title: "Book Two", Author: &models.Author{Firstname: "John2", Lastname: "Doe2"}})

	// Route Handlers / Endpoints
	r.HandleFunc("/api/login", controller.Login).Methods("POST")
	r.Handle("/api/books", middleware.IsAuthorized(controller.GetBooks)).Methods("GET")
	r.Handle("/api/books/{id}", middleware.IsAuthorized(controller.GetBook)).Methods("GET")
	r.Handle("/api/books", middleware.IsAuthorized(controller.CreateBooks)).Methods("POST")
	r.Handle("/api/books/{id}", middleware.IsAuthorized(controller.UpdateBooks)).Methods("PUT")
	r.Handle("/api/books", middleware.IsAuthorized(controller.DeleteBooks)).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))

}
