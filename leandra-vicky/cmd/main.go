package main

import (
	"biblioteca/handler"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	// /livros para GET (listar) e POST (adicionar)
	http.HandleFunc("/livros", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetAll(w, r)
		case http.MethodPost:
			handler.PostNewBook(w, r)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	})

	// /livros/{id} para GET, DELETE, PUT
	// /livros/{id}/lido para PATCH
	http.HandleFunc("/livros/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// PATCH /livros/{id}/lido
		if strings.HasSuffix(path, "/lido") && r.Method == http.MethodPatch {
			handler.PatchReadStatus(w, r)
			return
		}

		// /livros/{id}
		if strings.Count(path, "/") == 2 && len(path) > len("/livros/") {
			switch r.Method {
			case http.MethodGet:
				handler.GetByID(w, r)
			case http.MethodDelete:
				handler.DeleteById(w, r)
			case http.MethodPut:
				handler.UpdateById(w, r)
			default:
				http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			}
			return
		}

		http.NotFound(w, r)
	})

	fmt.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", nil)
}
