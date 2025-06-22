package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/JimmyArgote/biblia-api-go/src/models"
)

type CapituloRepository struct {
	connection *sql.DB
}

func NewCapituloRepository(connection *sql.DB) CapituloRepository {
	return CapituloRepository{
		connection: connection,
	}
}

func (cr *CapituloRepository) GetCapitulosByLivroId(id_livro int) ([]models.Capitulo, error) {

	rows, err := cr.connection.Query("SELECT id, livro_id, versao_id, titulo FROM capitulo WHERE livro_id = ? ORDER BY id", id_livro)

	if err != nil {
		log.Printf("Erro ao executar a query de capitulo: %v", err)
		return []models.Capitulo{}, err
	}
	defer rows.Close()

	capitulosList := []models.Capitulo{}

	for rows.Next() {
		var capituloObj models.Capitulo
		err := rows.Scan(
			&capituloObj.ID,
			&capituloObj.LivroID,
			&capituloObj.VersaoID,
			&capituloObj.Titulo)

		if err != nil {
			log.Printf("Erro ao escanear a linha do capitulo: %v", err)
			fmt.Println(err)
			return nil, err
		}
		capitulosList = append(capitulosList, capituloObj)
	}

	// Verifique por erros durante a iteração
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(capitulosList) == 0 {
		return nil, sql.ErrNoRows
	}

	return capitulosList, nil
}

func (cr *CapituloRepository) GetCapituloByLivroIdAndCapituloId(id_livro int, id_capitulo int) (*models.Capitulo, error) {

	query, err := cr.connection.Prepare("SELECT id, livro_id, versao_id, titulo FROM capitulo WHERE livro_id = ? AND id = ?")

	if err != nil {
		log.Printf("Erro ao executar a query de capitulo: %v", err)
		return nil, err
	}

	var capitulo models.Capitulo

	err = query.QueryRow(id_livro, id_capitulo).Scan(
		&capitulo.ID,
		&capitulo.LivroID,
		&capitulo.VersaoID,
		&capitulo.Titulo)

	if err != nil {
		log.Printf("Erro ao escanear a linha do capitulo: %v", err)
		fmt.Println(err)
		return nil, err
	}

	query.Close()

	return &capitulo, err
}
