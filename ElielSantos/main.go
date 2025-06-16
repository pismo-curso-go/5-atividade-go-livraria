package main

import (
    "net/http"
    "strconv"

    "github.com/labstack/echo/v4"
)

type Livro struct {
    ID     int    `json:"id"`
    Titulo string `json:"titulo"`
    Autor  string `json:"autor"`
    Lido   bool   `json:"lido"`
}

var livros []Livro
var nextID = 1

func listarLivros(c echo.Context) error {
    return c.JSON(http.StatusOK, livros)
}

func buscarLivro(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"erro": "ID inválido"})
    }
    for _, l := range livros {
        if l.ID == id {
            return c.JSON(http.StatusOK, l)
        }
    }
    return c.JSON(http.StatusNotFound, echo.Map{"erro": "Livro não encontrado"})
}

func adicionarLivro(c echo.Context) error {
    var livro Livro
    if err := c.Bind(&livro); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"erro": "JSON inválido"})
    }
    livro.ID = nextID
    nextID++
    livros = append(livros, livro)
    return c.JSON(http.StatusCreated, livro)
}

func removerLivro(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"erro": "ID inválido"})
    }
    for i, l := range livros {
        if l.ID == id {
            livros = append(livros[:i], livros[i+1:]...)
            return c.NoContent(http.StatusNoContent)
        }
    }
    return c.JSON(http.StatusNotFound, echo.Map{"erro": "Livro não encontrado"})
}

func editarLivro(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"erro": "ID inválido"})
    }
    var dados Livro
    if err := c.Bind(&dados); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"erro": "JSON inválido"})
    }
    for i, l := range livros {
        if l.ID == id {
            l.Titulo = dados.Titulo
            l.Autor = dados.Autor
            l.Lido = dados.Lido
            livros[i] = l
            return c.JSON(http.StatusOK, l)
        }
    }
    return c.JSON(http.StatusNotFound, echo.Map{"erro": "Livro não encontrado"})
}

func editarStatusLeitura(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"erro": "ID inválido"})
    }
    var body struct {
        Lido bool `json:"lido"`
    }
    if err := c.Bind(&body); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"erro": "JSON inválido"})
    }
    for i, l := range livros {
        if l.ID == id {
            l.Lido = body.Lido
            livros[i] = l
            return c.JSON(http.StatusOK, l)
        }
    }
    return c.JSON(http.StatusNotFound, echo.Map{"erro": "Livro não encontrado"})
}

func main() {
    e := echo.New()

    e.GET("/livros", listarLivros)
    e.GET("/livros/:id", buscarLivro)
    e.POST("/livros", adicionarLivro)
    e.DELETE("/livros/:id", removerLivro)
    e.PUT("/livros/:id", editarLivro)
    e.PATCH("/livros/:id/lido", editarStatusLeitura)

    e.Logger.Fatal(e.Start(":8080"))
}
