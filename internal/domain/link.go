package domain

// Link é a representação central de um link encurtado em nosso sistema.
// Ele não sabe sobre banco de dados, JSON ou HTTP. É apenas a estrutura de dados pura.
type Link struct {
	ShortCode   string
	OriginalURL string
}