package models

// Capitulo representa um capítulo de um livro.
type Capitulo struct {
	ID       int    `json:"id"`
	LivroID  int    `json:"livro_id"`
	VersaoID int    `json:"versao_id"`
	Titulo   string `json:"titulo"`
}
