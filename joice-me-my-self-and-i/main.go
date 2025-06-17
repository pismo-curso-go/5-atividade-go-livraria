package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Livro struct {
	ID     int    `json:"id"`
	Titulo string `json:"titulo"`
	Autor  string `json:"autor"`
	Lido   bool   `json:"lido"`
}

var livros = []Livro{
	{ID: 1, Titulo: "Ensaio sobre a cegueira", Autor: "José Saramago", Lido: true},
	{ID: 2, Titulo: "A metamorfose", Autor: "Franz Kafka", Lido: false},
}

func gerarNovoID() int {
	maiorID := 0
	for _, livro := range livros {
		if livro.ID > maiorID {
			maiorID = livro.ID
		}
	}
	return maiorID + 1
}

func handleLivros(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(livros)

	case http.MethodPost:
		var novoLivro Livro
		err := json.NewDecoder(r.Body).Decode(&novoLivro)
		if err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		novoLivro.ID = gerarNovoID()
		livros = append(livros, novoLivro)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(novoLivro)

	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func handleLivroPorID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/livros/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	for _, livro := range livros {
		if livro.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(livro)
			return
		}
	}

	http.Error(w, "Livro não encontrado", http.StatusNotFound)
}

func handleRemoverLivro(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/livros/delete/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	for i, livro := range livros {
		if livro.ID == id {
			livros = append(livros[:i], livros[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Livro não encontrado", http.StatusNotFound)
}

func handleEditarLivro(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/livros/edit/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var livroEditado Livro
	err = json.NewDecoder(r.Body).Decode(&livroEditado)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	for i, livro := range livros {
		if livro.ID == id {
			livroEditado.ID = id
			livros[i] = livroEditado
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(livroEditado)
			return
		}
	}

	http.Error(w, "Livro não encontrado", http.StatusNotFound)
}

func handleAtualizarStatusLido(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/livros/lido/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	type LeituraStatus struct {
		Lido bool `json:"lido"`
	}

	var status LeituraStatus
	err = json.NewDecoder(r.Body).Decode(&status)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	for i := range livros {
		if livros[i].ID == id {
			livros[i].Lido = status.Lido
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(livros[i])
			return
		}
	}

	http.Error(w, "Livro não encontrado", http.StatusNotFound)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/livros/delete/", handleRemoverLivro)      // DELETE
	mux.HandleFunc("/livros/edit/", handleEditarLivro)         // PUT
	mux.HandleFunc("/livros/lido/", handleAtualizarStatusLido) // PATCH
	mux.HandleFunc("/livros/", handleLivroPorID)               // GET por ID
	mux.HandleFunc("/livros", handleLivros)                    // GET e POST

	log.Println("Servidor rodando na porta 8080 🚀")
	http.ListenAndServe(":8080", mux)
}

// Agora vai em nome e jesus..
