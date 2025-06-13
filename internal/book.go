package internal

import (
	"errors"
	"sync"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Read   bool   `json:"read"`
}

var (
	books  []Book
	nextID = 1
	mu     sync.Mutex
)

func GetAllBooks() []Book {
	mu.Lock()
	defer mu.Unlock()
	return books
}

func AddBook(b Book) Book {
	mu.Lock()
	defer mu.Unlock()

	b.ID = nextID
	nextID++
	books = append(books, b)
	return b
}

func GetBookByID(id int) (*Book, error) {
	mu.Lock()
	defer mu.Unlock()

	for _, book := range books {
		if book.ID == id {
			return &book, nil
		}
	}

	return nil, errors.New("book not found")
}

func UpdateBook(id int, updated Book) (Book, error) {
	mu.Lock()
	defer mu.Unlock()

	for i, book := range books {
		if book.ID == id {
			updated.ID = id
			books[i] = updated
			return updated, nil
		}
	}
	return Book{}, errors.New("cant update book")
}

func UpdateReadStatus(id int, read bool) (Book, error) {
	mu.Lock()
	defer mu.Unlock()

	for i := range books {
		if books[i].ID == id {
			books[i].Read = read
			return books[i], nil
		}
	}
	return Book{}, errors.New("cant update read status")
}

func DeleteBookByID(id int) error {
	mu.Lock()
	defer mu.Unlock()

	var bookIsPresent bool
	var bookIndex int

	for idx, book := range books {
		if book.ID == id {
			bookIsPresent = true
			bookIndex = idx
			break
		}
	}

	if bookIsPresent {
		books = append(books[:bookIndex], books[bookIndex+1:]...)
	}
	return nil
}
