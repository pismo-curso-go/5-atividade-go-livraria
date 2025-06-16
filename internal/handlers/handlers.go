package handlers

import (
	"bookstore/internal/models"
	"bookstore/internal/services"
	"bookstore/internal/utils"
	"net/http"
	"strconv"
	"strings"
)

type BookHandlers struct {
	bookService *services.BookService
}

func NewBookHandlers(bookService *services.BookService) *BookHandlers {
	return &BookHandlers{
		bookService: bookService,
	}
}

func (bh *BookHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Adiciona CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/books")

	switch {
	case path == "" || path == "/":
		bh.handleBooksCollection(w, r)
	case strings.HasSuffix(path, "/read"):
		bh.handleReadStatus(w, r, path)
	default:
		bh.handleBookByID(w, r, path)
	}
}

func (bh *BookHandlers) handleBooksCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		bh.getAllBooks(w, r)
	case "POST":
		bh.addBook(w, r)
	default:
		utils.WriteJSONError(w, utils.ErrMethodNotAllowed)
	}
}

func (bh *BookHandlers) handleBookByID(w http.ResponseWriter, r *http.Request, path string) {
	idStr := strings.Trim(path, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSONError(w, utils.ErrInvalidID)
		return
	}

	switch r.Method {
	case "GET":
		bh.getBookByID(w, r, id)
	case "PUT":
		bh.updateBook(w, r, id)
	case "DELETE":
		bh.deleteBook(w, r, id)
	default:
		utils.WriteJSONError(w, utils.ErrMethodNotAllowed)
	}
}

func (bh *BookHandlers) handleReadStatus(w http.ResponseWriter, r *http.Request, path string) {
	if r.Method != "PATCH" {
		utils.WriteJSONError(w, utils.ErrMethodNotAllowed)
		return
	}

	idStr := strings.TrimSuffix(strings.Trim(path, "/"), "/read")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSONError(w, utils.ErrInvalidID)
		return
	}

	bh.updateReadStatus(w, r, id)
}

func (bh *BookHandlers) getAllBooks(w http.ResponseWriter, r *http.Request) {
	books := bh.bookService.GetAllBooks()
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"books": books,
		"total": len(books),
	})
}

func (bh *BookHandlers) getBookByID(w http.ResponseWriter, r *http.Request, id int) {
	book, err := bh.bookService.GetBookByID(id)
	if err != nil {
		utils.WriteJSONError(w, err.(*utils.APIError))
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, book)
}

func (bh *BookHandlers) addBook(w http.ResponseWriter, r *http.Request) {
	var req models.BookRequest
	if err := utils.DecodeJSONRequest(r, &req); err != nil {
		utils.WriteJSONError(w, err.(*utils.APIError))
		return
	}

	book, err := bh.bookService.AddBook(req)
	if err != nil {
		utils.WriteJSONError(w, err.(*utils.APIError))
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, book)
}

func (bh *BookHandlers) updateBook(w http.ResponseWriter, r *http.Request, id int) {
	var req models.BookRequest
	if err := utils.DecodeJSONRequest(r, &req); err != nil {
		utils.WriteJSONError(w, err.(*utils.APIError))
		return
	}

	book, err := bh.bookService.UpdateBook(id, req)
	if err != nil {
		utils.WriteJSONError(w, err.(*utils.APIError))
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, book)
}

func (bh *BookHandlers) updateReadStatus(w http.ResponseWriter, r *http.Request, id int) {
	var req models.ReadStatusRequest
	if err := utils.DecodeJSONRequest(r, &req); err != nil {
		utils.WriteJSONError(w, err.(*utils.APIError))
		return
	}

	book, err := bh.bookService.UpdateReadStatus(id, req.Read)
	if err != nil {
		utils.WriteJSONError(w, err.(*utils.APIError))
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, book)
}

func (bh *BookHandlers) deleteBook(w http.ResponseWriter, r *http.Request, id int) {
	if err := bh.bookService.DeleteBook(id); err != nil {
		utils.WriteJSONError(w, err.(*utils.APIError))
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Book successfully removed",
	})
}
