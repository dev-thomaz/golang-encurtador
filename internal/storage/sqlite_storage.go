package storage

import (
	"database/sql"
	"log"

	"github.com/dev-thomaz/golang-encurtador/internal/domain"
	_ "github.com/mattn/go-sqlite3"
)

// LinkStore define o "contrato" que nossa camada de storage deve seguir.
// Qualquer tipo que implementar todos estes métodos satisfaz a interface LinkStore.
type LinkStore interface {
	GetByShortCode(shortCode string) (*domain.Link, error)
	GetByOriginalURL(originalURL string) (*domain.Link, error)
	GetAll() ([]*domain.Link, error)
	Save(link *domain.Link) error
	ClearAll() (int64, error)
}

// SQLiteStorage é a implementação concreta da nossa interface LinkStore usando SQLite.
type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage cria e retorna uma nova instância de SQLiteStorage.
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS links (
		"short_code" TEXT NOT NULL PRIMARY KEY,
		"original_url" TEXT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_original_url ON links (original_url);`
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		return nil, err
	}

	return &SQLiteStorage{db: db}, nil
}

// Agora, transformamos as funções de banco de dados em métodos do SQLiteStorage.

func (s *SQLiteStorage) GetByShortCode(shortCode string) (*domain.Link, error) {
	var link domain.Link
	link.ShortCode = shortCode

	row := s.db.QueryRow("SELECT original_url FROM links WHERE short_code = ?", shortCode)
	err := row.Scan(&link.OriginalURL)
	if err == sql.ErrNoRows {
		return nil, nil // Retornamos nil, nil para indicar "não encontrado", mas sem erro.
	} else if err != nil {
		return nil, err
	}

	return &link, nil
}

func (s *SQLiteStorage) GetByOriginalURL(originalURL string) (*domain.Link, error) {
	var link domain.Link
	link.OriginalURL = originalURL

	row := s.db.QueryRow("SELECT short_code FROM links WHERE original_url = ?", originalURL)
	err := row.Scan(&link.ShortCode)
	if err == sql.ErrNoRows {
		return nil, nil // Não encontrado.
	} else if err != nil {
		return nil, err
	}
	return &link, nil
}

func (s *SQLiteStorage) GetAll() ([]*domain.Link, error) {
	rows, err := s.db.Query("SELECT short_code, original_url FROM links")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*domain.Link
	for rows.Next() {
		var link domain.Link
		if err := rows.Scan(&link.ShortCode, &link.OriginalURL); err != nil {
			log.Printf("Erro ao escanear a linha do banco: %s", err)
			continue
		}
		links = append(links, &link)
	}
	return links, nil
}

func (s *SQLiteStorage) Save(link *domain.Link) error {
	_, err := s.db.Exec("INSERT INTO links (short_code, original_url) VALUES (?, ?)", link.ShortCode, link.OriginalURL)
	return err
}

func (s *SQLiteStorage) ClearAll() (int64, error) {
	res, err := s.db.Exec("DELETE FROM links")
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}