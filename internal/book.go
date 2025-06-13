package internal

import (
	"sync"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Read   bool   `json:"read"`
}

var (
	books  = []Book{}
	nextID = 1
	mu     sync.Mutex
)

func GetAllBooks() []Book {
	mu.Lock()
	defer mu.Unlock()
	return books
}

func AddBook(b Book) *Book {
	mu.Lock()
	defer mu.Unlock()

	b.ID = nextID
	nextID++
	books = append(books, b)
	return &b
}

func GetBookByID(id int) (*Book, error) {
	mu.Lock()
	defer mu.Unlock()

	for idx := range books {
		if books[idx].ID == id {
			return &books[idx], nil
		}
	}

	return nil, ErrBookNotFound
}

func UpdateBook(id int, updatedBook Book) (*Book, error) {
	mu.Lock()
	defer mu.Unlock()

	for i, book := range books {
		if book.ID == id {
			updatedBook.ID = id
			books[i] = updatedBook
			return &updatedBook, nil
		}
	}
	return nil, ErrBookNotUpdated
}

func UpdateReadStatus(id int, read bool) (*Book, error) {
	mu.Lock()
	defer mu.Unlock()

	for i := range books {
		if books[i].ID == id {
			books[i].Read = read
			return &books[i], nil
		}
	}
	return nil, ErrBookStatusNotUpdated
}

func DeleteBookByID(id int) error {
	mu.Lock()
	defer mu.Unlock()

	var bookIsPresent bool
	var bookIndex int

	for idx := range books {
		if books[idx].ID == id {
			bookIsPresent = true
			bookIndex = idx
			break
		}
	}

	if bookIsPresent {
		books = append(books[:bookIndex], books[bookIndex+1:]...)
		return nil
	}

	return ErrBookCannotBeDeleted
}
