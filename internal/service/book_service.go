package services

import (
	"bookstore/internal/models"
	"bookstore/internal/utils"
	"sync"
)

type BookService struct {
	books  []models.Book
	nextID int
	mu     sync.RWMutex
}

func (bs *BookService) GetAllBooks() []models.Book {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	// Returns a copy to avoid external modifications
	booksCopy := make([]models.Book, len(bs.books))
	copy(booksCopy, bs.books)
	return booksCopy
}

func (bs *BookService) GetBookByID(id int) (*models.Book, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	for _, book := range bs.books {
		if book.ID == id {
			return &book, nil
		}
	}
	return nil, utils.ErrBookNotFound
}

func (bs *BookService) AddBook(req models.BookRequest) (*models.Book, error) {
	if req.Title == "" || req.Author == "" {
		return nil, utils.ErrInvalidBookData
	}

	bs.mu.Lock()
	defer bs.mu.Unlock()

	book := models.Book{
		ID:     bs.nextID,
		Title:  req.Title,
		Author: req.Author,
		Read:   req.Read,
	}

	bs.books = append(bs.books, book)
	bs.nextID++

	return &book, nil
}

func (bs *BookService) UpdateBook(id int, req models.BookRequest) (*models.Book, error) {
	if req.Title == "" || req.Author == "" {
		return nil, utils.ErrInvalidBookData
	}

	bs.mu.Lock()
	defer bs.mu.Unlock()

	for i, book := range bs.books {
		if book.ID == id {
			bs.books[i].Title = req.Title
			bs.books[i].Author = req.Author
			bs.books[i].Read = req.Read
			return &bs.books[i], nil
		}
	}

	return nil, utils.ErrBookNotFound
}

func (bs *BookService) UpdateReadStatus(id int, lido bool) (*models.Book, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	for i, book := range bs.books {
		if book.ID == id {
			bs.books[i].Read = lido
			return &bs.books[i], nil
		}
	}

	return nil, utils.ErrBookNotFound
}

func (bs *BookService) DeleteBook(id int) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	for i, book := range bs.books {
		if book.ID == id {
			bs.books = append(bs.books[:i], bs.books[i+1:]...)
			return nil
		}
	}

	return utils.ErrBookNotFound
}
