package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/JimmyArgote/biblia-api-go/models"
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

	query, err := cr.connection.Prepare("SELECT id, livro_id, versao_id, titulo FROM capitulo WHERE livro_id = ?")

	if err != nil {
		log.Printf("Erro ao executar a query de capitulo: %v", err)
		return []models.Capitulo{}, err
	}

	rows, err := query.Query(id_livro)

	if err != nil {
		log.Printf("Erro ao executar a query de capitulo: %v", err)
		return []models.Capitulo{}, err
	}

	capitulosList := []models.Capitulo{}
	capituloObj := models.Capitulo{}

	for rows.Next() {
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

	rows.Close()
	query.Close()

	return capitulosList, nil
}

func (cr *CapituloRepository) GetCapituloByLivroIdAndCapituloId(id_livro int, id_capitulo int) (*models.Capitulo, error) {

	query, err := cr.connection.Prepare("SELECT id, livro_id, versao_id, titulo FROM capitulo WHERE livro_id = ? AND id = ?")

	if err != nil {
		log.Printf("Erro ao executar a query de capitulo: %v", err)
		return nil, err
	}

	var capitulo models.Capitulo

	err = query.QueryRow(id_capitulo).Scan(
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

	return &capitulo, nil
}
