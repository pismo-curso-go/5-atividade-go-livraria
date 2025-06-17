package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// type Emprestavel interface {
// 	Emprestar() string
// 	Devolver() string
// }

type Livro struct {
	ID         int
	Titulo     string
	Autor      string
	Lido       bool
	Disponivel bool
}

var livros []Livro

func listrarLivros(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(livros)
}

func buscarLivro(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	idConvertido, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	for _, livro := range livros {
		if livro.ID == idConvertido {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(livro)
			return
		}
	}
	http.Error(w, "Livro não encontrado", http.StatusNotFound)
}

// func adicionarLivro(w, r) {}

// func removerLivro(w, r) {}

// func editarLivro(w, r) {}

// func marcarComoLido(w, r) {}

func registrarRotas() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/livros", listrarLivros).Methods("GET")
	router.HandleFunc("/livros/{id}", buscarLivro).Methods("GET")
	// router.HandleFunc("/livros", adicionarLivro).Methods("POST")
	// router.HandleFunc("/livros/{id}", removerLivro).Methods("DELETE")
	// router.HandleFunc("/livros/{id}", editarLivro).Methods("PUT")
	// router.HandleFunc("/livros/{id}/lido", marcarComoLido).Methods("PATCH")
	return router
}
func main() {
	livros = append(livros, Livro{ID: 1, Titulo: "Teste 1", Autor: "Esse livro é teste 1", Lido: false, Disponivel: true})
	livros = append(livros, Livro{ID: 2, Titulo: "Teste 2", Autor: "Esse livro é teste 2", Lido: false, Disponivel: false})
	livros = append(livros, Livro{ID: 3, Titulo: "Teste 3", Autor: "Esse livro é teste 3", Lido: false, Disponivel: false})

	router := registrarRotas()
	fmt.Println("Servidor rodando na porta 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
