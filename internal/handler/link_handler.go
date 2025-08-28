package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/dev-thomaz/golang-encurtador/internal/domain"
	"github.com/dev-thomaz/golang-encurtador/internal/storage" // Importamos o storage!
	"github.com/go-chi/chi/v5"
)

// As structs de requisição e resposta JSON pertencem à camada de handler,
// pois definem o "contrato" da API com o mundo exterior.
type ShortenRequest struct {
	URLs []string `json:"urls"`
}

type ShortenResult struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

// LinkHandler agora é uma struct que contém suas dependências,
// neste caso, qualquer coisa que satisfaça a interface LinkStore.
type LinkHandler struct {
	store storage.LinkStore
}

// NewLinkHandler é o nosso "construtor" para o handler.
// Ele recebe a dependência (o storage) e retorna um LinkHandler pronto para uso.
func NewLinkHandler(s storage.LinkStore) *LinkHandler {
	return &LinkHandler{
		store: s,
	}
}

// As funções de handler agora são MÉTODOS de LinkHandler.
// Elas usam `h.store` para interagir com a camada de dados.

func (h *LinkHandler) ShortenLinks(w http.ResponseWriter, r *http.Request) {
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

		// A lógica agora é mais limpa: apenas chama os métodos do storage.
		link, err := h.store.GetByOriginalURL(originalURL)
		if err != nil {
			log.Printf("Erro ao verificar URL no storage: %s", err)
			continue
		}

		if link == nil { // Se não encontrou, cria um novo.
			newLink := &domain.Link{
				ShortCode:   generateShortCode(6),
				OriginalURL: originalURL,
			}
			if err := h.store.Save(newLink); err != nil {
				log.Printf("Erro ao salvar novo link no storage: %s", err)
				continue
			}
			link = newLink
		}

		result := ShortenResult{
			OriginalURL: link.OriginalURL,
			ShortURL:    fmt.Sprintf("http://%s/%s", r.Host, link.ShortCode),
		}
		results = append(results, result)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Erro ao escrever resposta JSON: %s", err)
	}
}

func (h *LinkHandler) ListLinks(w http.ResponseWriter, r *http.Request) {
	links, err := h.store.GetAll()
	if err != nil {
		log.Printf("Erro ao consultar os links: %s", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	results := make([]ShortenResult, len(links))
	for i, link := range links {
		results[i] = ShortenResult{
			OriginalURL: link.OriginalURL,
			ShortURL:    fmt.Sprintf("http://%s/%s", r.Host, link.ShortCode),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Erro ao escrever resposta JSON: %s", err)
	}
}

func (h *LinkHandler) RedirectLink(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")
	link, err := h.store.GetByShortCode(shortCode)

	if err != nil {
		log.Printf("Erro ao buscar no storage: %s", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}
	if link == nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}

func (h *LinkHandler) ClearLinks(w http.ResponseWriter, r *http.Request) {
	rowsAffected, err := h.store.ClearAll()
	if err != nil {
		log.Printf("Erro ao limpar o storage: %s", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":         "Todos os links foram removidos com sucesso.",
		"links_removidos": rowsAffected,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao escrever resposta JSON: %s", err)
	}
}

// generateShortCode é uma função auxiliar que não precisa ser um método.
func generateShortCode(length int) string {
	// ... (código da função continua o mesmo)
	var sb strings.Builder
	// Para evitar a necessidade de uma semente global, podemos fazer isso.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	for i := 0; i < length; i++ {
		sb.WriteByte(chars[r.Intn(len(chars))])
	}
	return sb.String()
}