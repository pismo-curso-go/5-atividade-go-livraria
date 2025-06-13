package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"api/internal"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/books", booksHandler)
	mux.HandleFunc("/books/", bookByIDHandler)

	log.Println("server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetBooks(w, r)
	case http.MethodPost:
		handleAddBook(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetBooks(w http.ResponseWriter, _ *http.Request) {
	books := internal.GetAllBooks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func handleAddBook(w http.ResponseWriter, r *http.Request) {
	var book internal.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	book = internal.AddBook(book)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func bookByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	idStr = strings.TrimSuffix(idStr, "/read")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetBookByID(w, r, id)
	case http.MethodPut:
		handleUpdateBook(w, r, id)
	case http.MethodPatch:
		handleUpdateReadStatus(w, r, id)
	case http.MethodDelete:
		handleDeleteBook(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetBookByID(w http.ResponseWriter, _ *http.Request, id int) {
	book, err := internal.GetBookByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func handleUpdateBook(w http.ResponseWriter, r *http.Request, id int) {
	var updated internal.Book
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	book, err := internal.UpdateBook(id, updated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func handleUpdateReadStatus(w http.ResponseWriter, r *http.Request, id int) {
	var data struct {
		Read bool `json:"read"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	book, err := internal.UpdateReadStatus(id, data.Read)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func handleDeleteBook(w http.ResponseWriter, _ *http.Request, id int) {
	err := internal.DeleteBookByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}
