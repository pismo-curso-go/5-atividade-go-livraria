package handler

import (
	"biblioteca/model"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	proximoID int
	mu        sync.Mutex
)

// GetAll retorna todos os livros cadastrados
func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(model.Livros); err != nil {
		http.Error(w, "Erro ao codificar os livros", http.StatusInternalServerError)
		return
	}
}

// GetByID retorna o livro correspondete com o ID fornecido na URL
func GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := ExtrairID(r.URL.Path)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	for _, livro := range model.Livros {
		if livro.ID == id {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(livro); err != nil {
				http.Error(w, "Erro ao codificar o livro", http.StatusInternalServerError)
			}
			return
		}
	}

	http.Error(w, "Livro não encontrado", http.StatusNotFound)
}

// PostNewBook adiciona um novo livro à lista de livros
func PostNewBook(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // Limita o tamanho do body para 1MB

	var novoLivro model.Livro
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&novoLivro); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(novoLivro.Titulo) == "" || strings.TrimSpace(novoLivro.Autor) == "" {
		http.Error(w, "Título e autor são obrigatórios e não podem ser vazios", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	novoLivro.ID = proximoID
	proximoID++
	model.Livros = append(model.Livros, novoLivro)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(novoLivro); err != nil {
		http.Error(w, "Erro ao codificar o livro", http.StatusInternalServerError)
	}
}

// DeleteByID deleta o livro correspondente ao ID fornecido na URL
func DeleteById(w http.ResponseWriter, r *http.Request) {
	id, err := ExtrairID(r.URL.Path)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, livro := range model.Livros {
		if livro.ID == id {
			model.Livros = append(model.Livros[:i], model.Livros[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	http.Error(w, "Livro não encontrado", http.StatusNotFound)
}

// UpdateById atualiza o livro correspondente ao ID fornecido na URL
func UpdateById(w http.ResponseWriter, r *http.Request) {
	id, err := ExtrairID(r.URL.Path)
	if err != nil {
		http.Error(w, "ID Invalido", http.StatusBadRequest)
		return
	}

	var novo model.LivroEdicao
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&novo); err != nil {
		http.Error(w, "Json invalido", http.StatusBadRequest)
		return
	}

	if VerificarVazio([]string{novo.Titulo, novo.Autor}) {
		http.Error(w, "Os campos sao obrigatorios e nao podem ser vazios", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for x, livro := range model.Livros {
		if livro.ID == id {
			model.Livros[x].Titulo = novo.Titulo
			model.Livros[x].Autor = novo.Autor
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "Livro não encontrado", http.StatusNotFound)
}

func PatchReadStatus(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/livros/"), "/lido")
	idStr = strings.Trim(idStr, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var status struct {
		Lido bool `json:"lido"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&status); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, livro := range model.Livros {
		if livro.ID == id {
			model.Livros[i].Lido = status.Lido
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(model.Livros[i])
			return
		}
	}
	http.Error(w, "Livro não encontrado", http.StatusNotFound)
}

// ExtrairID extrai o ID do livro da URL
func ExtrairID(path string) (int, error) {
	partes := strings.Split(path, "/")
	idStr := partes[len(partes)-1]
	return strconv.Atoi(idStr)
}

// VerificarVazio verifica se algum dos campos fornecidos está vazio
func VerificarVazio(campos []string) bool {
	for _, str := range campos {
		if strings.TrimSpace(str) == "" {
			return true
		}
	}
	return false
}
