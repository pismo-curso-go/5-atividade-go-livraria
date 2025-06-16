package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Livro representa a estrutura de dados para um livro.
type Livro struct {
	ID     int    `json:"id"`
	Titulo string `json:"titulo"`
	Autor  string `json:"autor"`
	Lido   bool   `json:"lido"`
}

var livros []Livro

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/livros", handleLivros)
	mux.HandleFunc("/livros/", handleSpecificLivro)

	log.Print("Servidor iniciado na porta 8080")
	http.ListenAndServe(":8080", mux)
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

func handleSpecificLivro(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/livros/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		buscarLivroPorID(w, r, id)
	case http.MethodDelete:
		deletarLivroPorId(w, r, id)
	case http.MethodPut:
		editarLivro(w, r, id)
	case http.MethodPatch:
		atualizarParcialmenteLivro(w, r, id)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}

}

// ListarLivros lida com a requisição para listar todos os livros (GET /livros).
func listarLivros(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(livros)
}

// BuscarLivroPorID lida com a requisição para buscar um livro pelo seu ID (GET /livros/{id}).
func buscarLivroPorID(w http.ResponseWriter, r *http.Request, id int) {

	for _, livro := range livros {
		if livro.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(livro)
			return
		}
	}

	http.Error(w, "Livro não encontrado", http.StatusNotFound)
}

// EditarLivro lida com a requisição para editar os dados de um livro (PUT /livros/{id}).
func editarLivro(w http.ResponseWriter, r *http.Request, id int) {

	for i, livro := range livros {
		if livro.ID == id {
			var livroAtualizado Livro
			err := json.NewDecoder(r.Body).Decode(&livroAtualizado)
			if err != nil {
				http.Error(w, "Erro ao decodificar o corpo da requisição", http.StatusBadRequest)
				return
			}

			livroAtualizado.ID = id
			livros[i] = livroAtualizado

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(livroAtualizado)
			return
		}
	}

	// Se o loop terminar, o livro não foi encontrado.
	http.Error(w, "Livro não encontrado", http.StatusNotFound)
}

func atualizarParcialmenteLivro(w http.ResponseWriter, r *http.Request, id int) {

	var alteracoes map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&alteracoes); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	// 3. Encontrar o livro e aplicar as alterações
	var livroAtualizado Livro
	livroEncontrado := false

	for i, livro := range livros {
		if livro.ID == id {
			// Verifica cada campo recebido no mapa e atualiza o livro original
			if novoTitulo, ok := alteracoes["titulo"]; ok {
				if tituloStr, ok := novoTitulo.(string); ok {
					livros[i].Titulo = tituloStr
				} else {
					http.Error(w, "Tipo inválido para o campo 'titulo'", http.StatusBadRequest)
					return
				}
			}

			if novoAutor, ok := alteracoes["autor"]; ok {
				if autorStr, ok := novoAutor.(string); ok {
					livros[i].Autor = autorStr
				} else {
					http.Error(w, "Tipo inválido para o campo 'autor'", http.StatusBadRequest)
					return
				}
			}

			if novoLido, ok := alteracoes["lido"]; ok {
				if lidoBool, ok := novoLido.(bool); ok {
					livros[i].Lido = lidoBool
				} else {
					http.Error(w, "Tipo inválido para o campo 'lido'", http.StatusBadRequest)
					return
				}
			}

			livroAtualizado = livros[i]
			livroEncontrado = true
			break
		}
	}

	if !livroEncontrado {
		http.Error(w, "Livro não encontrado", http.StatusNotFound)
		return
	}

	// 4. Retornar o livro com os dados atualizados
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(livroAtualizado)
}

func criarLivro(w http.ResponseWriter, r *http.Request) {
	var newLivro Livro

	err := json.NewDecoder(r.Body).Decode(&newLivro)
	if err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
	}
	newLivro.ID = len(livros) + 1
	livros = append(livros, newLivro)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newLivro)
}

func deletarLivroPorId(w http.ResponseWriter, _ *http.Request, id int) {
	for i, livro := range livros {
		if livro.ID == id {
			livros = append(livros[:1], livros[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Livro não encontrado", http.StatusNotFound)
}
