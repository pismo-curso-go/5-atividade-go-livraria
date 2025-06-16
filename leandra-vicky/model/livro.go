package model

type Livro struct {
	ID     int    `json:"id"`
	Titulo string `json:"titulo"`
	Autor  string `json:"autor"`
	Lido   bool   `json:"lido"`
}

type LivroEdicao struct {
	Titulo string `json:"titulo"`
	Autor  string `json:"autor"`
}

var Livros = []Livro{
	{ID: 1, Titulo: "1984", Autor: "George Orwell", Lido: true},
	{ID: 2, Titulo: "O Senhor dos Anéis", Autor: "J.R.R. Tolkien", Lido: false},
	{ID: 3, Titulo: "Dom Casmurro", Autor: "Machado de Assis", Lido: true},
	{ID: 4, Titulo: "A Revolução dos Bichos", Autor: "George Orwell", Lido: false},
}

