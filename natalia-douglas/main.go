package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// type Emprestavel interface {
// 	Emprestar() string
// 	Devolver() string
// }

type Livro struct {
	ID         int
	Disponivel bool
	Lido       bool
	Titulo     string
	Autor      string
}

var livros []Livro

func listrarLivros(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(livros)
}

// func buscarLivro(w, r) {}

// func adicionarLivro(w, r) {}

// func removerLivro(w, r) {}

// func editarLivro(w, r) {}

// func marcarComoLido(w, r) {}

func registrarRotas() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/livros", listrarLivros).Methods("GET")
	// router.HandleFunc("/livros/{id}", buscarLivro).Methods("GET")
	// router.HandleFunc("/livros", adicionarLivro).Methods("POST")
	// router.HandleFunc("/livros/{id}", removerLivro).Methods("DELETE")
	// router.HandleFunc("/livros/{id}", editarLivro).Methods("PUT")
	// router.HandleFunc("/livros/{id}/lido", marcarComoLido).Methods("PATCH")
	return router
}
func main() {
	livros = append(livros, Livro{ID: 1, Titulo: "Dom Casmurro", Autor: "Machado de Assis", Lido: false})

	router := registrarRotas()
	fmt.Println("Servidor rodando na porta 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
