package main

import (
	"log"
	"net/http"
)

type Livro struct {
	Id     int64  `json:"id"`
	Titulo string `json:"titulo"`
	Autor  string `json:"autor"`
	Lido   bool   `json:"lido"`
}

var livros []Livro

func main() {
	mux := http.NewServeMux()

	log.Print("Servidor iniciado na porta 8080")
	http.ListenAndServe(":8080", mux)
}
