package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dev-thomaz/golang-encurtador/internal/handler"
	"github.com/dev-thomaz/golang-encurtador/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// 1. CONFIGURAÇÃO
	const port = "8080"
	const dbPath = "./storage.db"

	// 2. INICIALIZAÇÃO DAS CAMADAS (Injeção de Dependência)
	// Primeiro, criamos a camada de storage.
	store, err := storage.NewSQLiteStorage(dbPath)
	if err != nil {
		log.Fatalf("Erro ao inicializar o storage: %v", err)
	}

	// Em seguida, criamos a camada de handler, injetando o storage nela.
	linkHandler := handler.NewLinkHandler(store)

	// 3. CONFIGURAÇÃO DO ROTEADOR E REGISTRO DAS ROTAS
	r := chi.NewRouter()
	r.Use(middleware.Logger) // Um middleware útil do Chi para logar as requisições
	r.Use(middleware.Recoverer) // Middleware para "panics"

	// As rotas agora chamam os métodos do nosso handler.
	r.Post("/shorten", linkHandler.ShortenLinks)
	r.Get("/links", linkHandler.ListLinks)
	r.Delete("/links", linkHandler.ClearLinks)
	r.Get("/{shortCode}", linkHandler.RedirectLink)

	// 4. INICIALIZAÇÃO DO SERVIDOR
	fmt.Printf("Servidor V2 iniciando na porta %s...\n", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Não foi possível iniciar o servidor na porta %s: %s", port, err)
	}
}