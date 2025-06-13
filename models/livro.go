package models

// Livro representa um livro da Bíblia.
type Livro struct {
	ID         int    `json:"id"`
	Ordem      int    `json:"ordem"`
	Nome       string `json:"nome"`
	Sigla      string `json:"sigla"`
	Testamento string `json:"testamento"`
}
