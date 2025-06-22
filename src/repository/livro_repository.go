package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/JimmyArgote/biblia-api-go/models"
)

type LivroRepository struct {
	connection *sql.DB
}

func NewLivroRepository(connection *sql.DB) LivroRepository {
	return LivroRepository{
		connection: connection,
	}
}

func (lr *LivroRepository) GetLivros() ([]models.Livro, error) {

	query := `SELECT id, ordem, nome, sigla, testamento FROM livro ORDER BY testamento ASC, ordem ASC`
	rows, err := lr.connection.Query(query)
	if err != nil {
		log.Printf("Erro ao executar a query de livros: %v", err)
		return []models.Livro{}, err
	}
	defer rows.Close()

	livrosList := []models.Livro{}
	livroObj := models.Livro{}

	for rows.Next() {
		err := rows.Scan(
			&livroObj.ID,
			&livroObj.Ordem,
			&livroObj.Nome,
			&livroObj.Sigla,
			&livroObj.Testamento)

		if err != nil {
			log.Printf("Erro ao escanear a linha do livro: %v", err)
			fmt.Println(err)
			return nil, err
		}
		livrosList = append(livrosList, livroObj)
	}

	rows.Close()

	return livrosList, nil
}

func (lr *LivroRepository) GetLivroById(id_livro int) (*models.Livro, error) {

	query, err := lr.connection.Prepare("SELECT id, ordem, nome, sigla, testamento FROM livro WHERE id = ?")
	if err != nil {
		log.Printf("Erro ao executar a query de livro: %v", err)
		return nil, err
	}

	var livro models.Livro

	err = query.QueryRow(id_livro).Scan(
		&livro.ID,
		&livro.Ordem,
		&livro.Nome,
		&livro.Sigla,
		&livro.Testamento)

	if err != nil {
		log.Printf("Erro ao escanear a linha do livro: %v", err)
		fmt.Println(err)
		return nil, err
	}

	query.Close()

	return &livro, nil
}
