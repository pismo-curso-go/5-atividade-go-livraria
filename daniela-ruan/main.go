package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Livro struct {
	ID     int    `json:"id"`
	Titulo string `json:"titulo"`
	Autor  string `json:"autor"`
	Lido   bool   `json:"lido"`
}

var (
	livros    []Livro
	idCounter = 1
	mu        sync.Mutex
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/livros", handleLivros)
	mux.HandleFunc("/livros/", handleLivro) // Inclui /livros/{id} e /livros/{id}/lido

	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleLivros(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listarLivros(w, r)
	case http.MethodPost:
		criarLivro(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func handleLivro(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/livros/")
	// Verifica se a rota é /livros/{id}/lido
	if strings.HasSuffix(path, "/lido") {
		idStr := strings.TrimSuffix(path, "/lido")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}
		if r.Method == http.MethodPatch {
			atualizarStatusLeitura(w, r, id)
			return
		}
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Rota padrão /livros/{id}
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		buscarLivroPorID(w, r, id)
	case http.MethodDelete:
		deletarLivroPorID(w, r, id)
	case http.MethodPut:
		editarLivro(w, r, id)
	case http.MethodPatch:
		atualizarParcialmenteLivro(w, r, id)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func listarLivros(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(livros)
}

func buscarLivroPorID(w http.ResponseWriter, r *http.Request, id int) {
	mu.Lock()
	defer mu.Unlock()

	for _, livro := range livros {
		if livro.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(livro)
			return
		}
	}
	http.Error(w, "Livro não encontrado", http.StatusNotFound)
}

func criarLivro(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var newLivro Livro
	err := json.NewDecoder(r.Body).Decode(&newLivro)
	if err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Validação simples
	if newLivro.Titulo == "" || newLivro.Autor == "" {
		http.Error(w, "Título e Autor são obrigatórios", http.StatusBadRequest)
		return
	}

	newLivro.ID = idCounter
	idCounter++
	livros = append(livros, newLivro)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newLivro)
}

func deletarLivroPorID(w http.ResponseWriter, r *http.Request, id int) {
	mu.Lock()
	defer mu.Unlock()

	for i, livro := range livros {
		if livro.ID == id {
			livros = append(livros[:i], livros[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Livro não encontrado", http.StatusNotFound)
}

func editarLivro(w http.ResponseWriter, r *http.Request, id int) {
	mu.Lock()
	defer mu.Unlock()

	var livroAtualizado Livro
	err := json.NewDecoder(r.Body).Decode(&livroAtualizado)
	if err != nil {
		http.Error(w, "Erro ao decodificar o corpo da requisição", http.StatusBadRequest)
		return
	}

	if liv
