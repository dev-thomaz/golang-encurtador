package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5" // 1. IMPORTANDO O CHI
)

// Nossas structs e variáveis globais continuam as mesmas.
var urls = make(map[string]string)
var mutex = &sync.Mutex{}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
var random = rand.New(rand.NewSource(time.Now().UnixNano()))

type ShortenRequest struct {
	URLs []string `json:"urls"`
}
type ShortenResult struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}
type ShortenResponse struct {
	Results []ShortenResult `json:"results"`
}

func generateShortCode(length int) string {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(chars[random.Intn(len(chars))])
	}
	return sb.String()
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	// Com o chi, não precisamos mais verificar o método HTTP aqui,
	// pois a rota só será acionada por um POST.
	// if r.Method != http.MethodPost { ... } -> REMOVIDO!

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

		shortCode := generateShortCode(6)
		mutex.Lock()
		urls[shortCode] = originalURL
		mutex.Unlock()

		result := ShortenResult{
			OriginalURL: originalURL,
			ShortURL:    fmt.Sprintf("http://localhost:8080/%s", shortCode),
		}
		results = append(results, result)
	}

	response := ShortenResponse{Results: results}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao escrever resposta JSON: %s", err)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Também não precisamos mais verificar o método GET.

	// 2. PEGANDO O PARÂMETRO DA URL COM O CHI
	// Em vez de 'r.URL.Path[1:]', usamos chi.URLParam. É mais limpo e seguro.
	shortCode := chi.URLParam(r, "shortCode")
	if shortCode == "" {
		http.NotFound(w, r)
		return
	}

	mutex.Lock()
	originalURL, found := urls[shortCode]
	mutex.Unlock()

	if found {
		http.Redirect(w, r, originalURL, http.StatusFound)
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	// 3. CONFIGURANDO AS ROTAS COM O CHI
	r := chi.NewRouter() // Criamos uma nova instância do roteador chi.

	// Definimos nossa rota POST para encurtar. A sintaxe é muito mais clara.
	r.Post("/shorten", shortenHandler)

	// Definimos nossa rota GET para redirecionar.
	// O trecho `{shortCode}` define um parâmetro de URL que podemos capturar.
	r.Get("/{shortCode}", redirectHandler)
	
	port := "8080"
	fmt.Printf("Servidor (com Chi) iniciando na porta %s...\n", port)
	
	// 4. USANDO O ROTEADOR DO CHI NO SERVIDOR
	// Em vez de 'nil', passamos nosso roteador 'r' para o ListenAndServe.
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Não foi possível iniciar o servidor na porta %s: %s", port, err)
	}
}