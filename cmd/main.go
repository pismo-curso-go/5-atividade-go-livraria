package main

import (
	"bookstore/internal/server"
	"flag"
	"log"
)

func main() {
	// Define flags de linha de comando
	var port string
	flag.StringVar(&port, "port", ":8080", "Porta do servidor (ex: :8080)")
	flag.Parse()

	// Cria e inicializa o servidor
	srv := server.NewServer(port)
	srv.Initialize()

	// Inicia o servidor
	if err := srv.Start(); err != nil {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}
