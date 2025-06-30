package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Livro representa a estrutura de um livro.
type Livro struct {
	ID     int    `json:"id"`
	Titulo string `json:"titulo"`
	Autor  string `json:"autor"`
	Lido   bool   `json:"lido"`
}

// Livraria gerencia a coleção de livros.
type Livraria struct {
	livros    []Livro
	proximoID int
	mu        sync.Mutex
}

// NovaLivraria cria uma instância inicial da livraria.
func NovaLivraria() *Livraria {
	return &Livraria{
		livros: []Livro{
			{ID: 1, Titulo: "Dom Quixote", Autor: "Miguel de Cervantes", Lido: true},
			{ID: 2, Titulo: "1984", Autor: "George Orwell", Lido: false},
		},
		proximoID: 3,
	}
}

// Handlers (Manipuladores de Requisição)

func (l *Livraria) listarLivros(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderComErro(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}
	responderComJSON(w, http.StatusOK, l.livros)
}

func (l *Livraria) obterLivro(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderComErro(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}
	id, err := extrairID(r.URL.Path, "/livros/obter/")
	if err != nil {
		responderComErro(w, http.StatusBadRequest, "ID de livro inválido")
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	for _, livro := range l.livros {
		if livro.ID == id {
			responderComJSON(w, http.StatusOK, livro)
			return
		}
	}
	responderComErro(w, http.StatusNotFound, "Livro não encontrado")
}

func (l *Livraria) adicionarLivro(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responderComErro(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	var novoLivro Livro
	if err := json.NewDecoder(r.Body).Decode(&novoLivro); err != nil {
		responderComErro(w, http.StatusBadRequest, "Dados de requisição inválidos")
		return
	}

	if novoLivro.Titulo == "" || novoLivro.Autor == "" {
		responderComErro(w, http.StatusBadRequest, "Título e autor são obrigatórios")
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	novoLivro.ID = l.proximoID
	l.proximoID++
	l.livros = append(l.livros, novoLivro)

	responderComJSON(w, http.StatusCreated, novoLivro)
}

func (l *Livraria) excluirLivro(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		responderComErro(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}
	id, err := extrairID(r.URL.Path, "/livros/excluir/")
	if err != nil {
		responderComErro(w, http.StatusBadRequest, "ID de livro inválido")
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	for i, livro := range l.livros {
		if livro.ID == id {
			l.livros = append(l.livros[:i], l.livros[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	responderComErro(w, http.StatusNotFound, "Livro não encontrado")
}

func (l *Livraria) atualizarLivro(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		responderComErro(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}
	id, err := extrairID(r.URL.Path, "/livros/atualizar/")
	if err != nil {
		responderComErro(w, http.StatusBadRequest, "ID de livro inválido")
		return
	}

	var livroAtualizado Livro
	if err := json.NewDecoder(r.Body).Decode(&livroAtualizado); err != nil {
		responderComErro(w, http.StatusBadRequest, "Dados de requisição inválidos")
		return
	}

	if livroAtualizado.Titulo == "" || livroAtualizado.Autor == "" {
		responderComErro(w, http.StatusBadRequest, "Título e autor são obrigatórios")
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	for i, livro := range l.livros {
		if livro.ID == id {
			livroAtualizado.ID = id
			livroAtualizado.Lido = livro.Lido // Mantém o status 'lido' original
			l.livros[i] = livroAtualizado
			responderComJSON(w, http.StatusOK, livroAtualizado)
			return
		}
	}
	responderComErro(w, http.StatusNotFound, "Livro não encontrado")
}

func (l *Livraria) atualizarStatusLeitura(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		responderComErro(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}
	id, err := extrairID(r.URL.Path, "/livros/ler/")
	if err != nil {
		responderComErro(w, http.StatusBadRequest, "ID de livro inválido")
		return
	}

	var dados struct {
		Lido bool `json:"lido"`
	}
	if err := json.NewDecoder(r.Body).Decode(&dados); err != nil {
		responderComErro(w, http.StatusBadRequest, "Dados de requisição inválidos")
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	for i, livro := range l.livros {
		if livro.ID == id {
			l.livros[i].Lido = dados.Lido
			responderComJSON(w, http.StatusOK, l.livros[i])
			return
		}
	}
	responderComErro(w, http.StatusNotFound, "Livro não encontrado")
}

// Funções Auxiliares

func responderComErro(w http.ResponseWriter, codigo int, msg string) {
	responderComJSON(w, codigo, map[string]string{"erro": msg})
}

func responderComJSON(w http.ResponseWriter, codigo int, dados interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)
	json.NewEncoder(w).Encode(dados)
}

func extrairID(path, prefixo string) (int, error) {
	if !strings.HasPrefix(path, prefixo) {
		return 0, fmt.Errorf("prefixo inválido")
	}
	idStr := strings.TrimPrefix(path, prefixo)
	return strconv.Atoi(idStr)
}

// Rotas

func configurarRotas(l *Livraria) {
	http.HandleFunc("/livros", l.listarLivros)
	http.HandleFunc("/livros/obter/", l.obterLivro)
	http.HandleFunc("/livros/adicionar", l.adicionarLivro)
	http.HandleFunc("/livros/excluir/", l.excluirLivro)
	http.HandleFunc("/livros/atualizar/", l.atualizarLivro)
	http.HandleFunc("/livros/ler/", l.atualizarStatusLeitura)
}

func main() {
	livraria := NovaLivraria()
	configurarRotas(livraria)

	fmt.Println("Servidor rodando na porta 8080...")
	http.ListenAndServe(":8080", nil)
}
