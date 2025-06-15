package models

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Read   bool   `json:"read"`
}

type BookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Read   bool   `json:"read"`
}

type ReadStatusRequest struct {
	Read bool `json:"read"`
}
