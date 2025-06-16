package server

import (
	"bookstore/internal/handlers"
	"bookstore/internal/services"
	"log"
	"net/http"
)

const (
	DefaultPort = ":8080"
)



func StartServer() {
	
	log.Printf("Escutando na porta %s", DefaultPort)
	log.Fatal(http.ListenAndServe(DefaultPort, CreateMux()))

}

func CreateMux() *http.ServeMux {
	bookService := services.NewBookService()
	bookHandlers := handlers.NewBookHandlers(bookService)

	mux := http.NewServeMux()

	mux.HandleFunc("/books/", bookHandlers.ServeHTTP)

	return mux
}

func ShowRoutes() {
	routes := []string{
		"GET /books - Get all books",
		"POST /books - Add a new book",
		"GET /books/{id} - Get book by ID",
		"PUT /books/{id} - Update book by ID",
		"PATCH /books/{id}/read - Update read status of a book by ID",
		"DELETE /books/{id} - Delete book by ID",
	}

	log.Println("Disponible Routes:")
	for _, route := range routes {
		log.Println(route)
	}
}