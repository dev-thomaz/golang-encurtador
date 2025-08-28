package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
var random = rand.New(rand.NewSource(time.Now().UnixNano()))

type ShortenRequest struct {
	URLs []string `json:"urls"`
}
type ShortenResult struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

// ... (initDB e outras funções permanecem as mesmas) ...
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./storage.db")
	if err != nil {
		log.Fatalf("Erro ao abrir o banco de dados: %s", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS links (
		"short_code" TEXT NOT NULL PRIMARY KEY,
		"original_url" TEXT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Erro ao criar a tabela no banco de dados: %s", err)
	}

	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_original_url ON links (original_url);`
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		log.Fatalf("Erro ao criar índice no banco de dados: %s", err)
	}
}

func generateShortCode(length int) string {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(chars[random.Intn(len(chars))])
	}
	return sb.String()
}

// NOVA FUNÇÃO HANDLER PARA LIMPAR OS LINKS
func clearLinksHandler(w http.ResponseWriter, r *http.Request) {
	// 1. EXECUTANDO O COMANDO DELETE
	// `db.Exec` é perfeito para isso, pois DELETE não retorna linhas.
	res, err := db.Exec("DELETE FROM links")
	if err != nil {
		log.Printf("Erro ao limpar a tabela de links: %s", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	// 2. PEGANDO O NÚMERO DE LINHAS AFETADAS
	// O resultado de `Exec` nos informa quantas linhas foram afetadas pela operação.
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Erro ao obter o número de linhas afetadas: %s", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	// 3. ENVIANDO A RESPOSTA DE CONFIRMAÇÃO
	// Criamos um map simples para construir nossa resposta JSON dinamicamente.
	response := map[string]interface{}{
		"message":         "Todos os links foram removidos com sucesso.",
		"links_removidos": rowsAffected,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao escrever resposta JSON: %s", err)
	}
}


// ... (listLinksHandler, shortenHandler, redirectHandler permanecem os mesmos) ...
func listLinksHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT short_code, original_url FROM links")
	if err != nil {
		log.Printf("Erro ao consultar os links: %s", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	links := make([]ShortenResult, 0)
	for rows.Next() {
		var shortCode, originalURL string
		if err := rows.Scan(&shortCode, &originalURL); err != nil {
			log.Printf("Erro ao escanear a linha do banco: %s", err)
			continue
		}
		link := ShortenResult{
			OriginalURL: originalURL,
			ShortURL:    fmt.Sprintf("http://localhost:8080/%s", shortCode),
		}
		links = append(links, link)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(links); err != nil {
		log.Printf("Erro ao escrever resposta JSON: %s", err)
	}
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}
	if len(req.URLs) == 0 {
		http.Error(w, "O array 'urls' não pode ser vazio", http.StatusBadRequest)
		return
	}
	results := make([]ShortenResult, 0, len(req.URLs))
	for _, originalURL := range req.URLs {
		if originalURL == "" {
			continue
		}
		var shortCode string
		row := db.QueryRow("SELECT short_code FROM links WHERE original_url = ?", originalURL)
		err := row.Scan(&shortCode)
		if err == sql.ErrNoRows {
			shortCode = generateShortCode(6)
			_, insertErr := db.Exec("INSERT INTO links (short_code, original_url) VALUES (?, ?)", shortCode, originalURL)
			if insertErr != nil {
				log.Printf("Erro ao inserir nova URL no banco de dados: %s", insertErr)
				continue
			}
		} else if err != nil {
			log.Printf("Erro ao verificar URL no banco de dados: %s", err)
			continue
		}
		result := ShortenResult{
			OriginalURL: originalURL,
			ShortURL:    fmt.Sprintf("http://localhost:8080/%s", shortCode),
		}
		results = append(results, result)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Erro ao escrever resposta JSON: %s", err)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")
	if shortCode == "" {
		http.NotFound(w, r)
		return
	}
	var originalURL string
	row := db.QueryRow("SELECT original_url FROM links WHERE short_code = ?", shortCode)
	err := row.Scan(&originalURL)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		log.Printf("Erro ao buscar no banco de dados: %s", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, originalURL, http.StatusFound)
}


func main() {
	initDB()
	r := chi.NewRouter()

	r.Post("/shorten", shortenHandler)
	r.Get("/links", listLinksHandler)
	// ADICIONANDO A NOVA ROTA DELETE
	r.Delete("/links", clearLinksHandler)
	r.Get("/{shortCode}", redirectHandler)

	port := "8080"
	fmt.Printf("Servidor (versão final) iniciando na porta %s...\n", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Não foi possível iniciar o servidor na porta %s: %s", port, err)
	}
}