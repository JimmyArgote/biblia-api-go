package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/JimmyArgote/biblia-api-go/models"
	"github.com/gin-gonic/gin"
)

type VersiculoRepository struct {
	connection *sql.DB
}

func NewVersiculoRepository(connection *sql.DB) VersiculoRepository {
	return VersiculoRepository{
		connection: connection,
	}
}

func (vr *VersiculoRepository) GetVersiculosByCapituloIdAndLivroId(id_capitulo int, id_livro int) (models.LivroCapVers, error) {

	var livroCapVers models.LivroCapVers
	var err error

	queryInfo, err := vr.connection.Prepare(`
		SELECT
			l.nome,
			l.sigla,
			l.testamento,
			(SELECT COUNT(1) FROM versiculo v WHERE v.livro_id = l.id AND v.capitulo_id = ?) as qtd_vers,
			(SELECT COUNT(1) FROM capitulo c WHERE c.livro_id = l.id) as qtd_caps
		FROM livro l
		WHERE l.id = ?`)

	rows, err := queryInfo.Query(queryInfo, id_capitulo, id_livro)

	if err != nil {
		log.Printf("Erro ao executar a query de versiculo: %v", err)
		return models.LivroCapVers{}, err
	}

	queryInfo.Close()

	for rows.Next() {
		err := rows.Scan(
			&livroCapVers.LivroNome,
			&livroCapVers.LivroSigla,
			&livroCapVers.Testamento,
			&livroCapVers.VersTotal,
			&livroCapVers.CapsTotal,
		)

		if err != nil {
			log.Printf("Erro ao escanear a linha do versiculo: %v", err)
			fmt.Println(err)
			return models.LivroCapVers{}, err
		}
	}

	livroID := strconv.Itoa(id_livro)
	capituloID := strconv.Itoa(id_capitulo)

	queryVersiculos, err := vr.connection.Prepare(`
		SELECT numero, formatado
		FROM versiculo
		WHERE livro_id = ? AND capitulo_id = ?
		ORDER BY numero ASC`)

	if err != nil {
		log.Printf("Erro ao executar query de versículos: %v", err)
		fmt.Println(err)
		return models.LivroCapVers{}, err
	}

	rows, err = queryVersiculos.Query(livroID, capituloID)

	if err != nil {
		log.Printf("Erro ao executar a query de versiculo: %v", err)
		return models.LivroCapVers{}, err
	}

	var versiculosSlimList []models.VersiculoSlim
	for rows.Next() {
		var vs models.VersiculoSlim
		if err := rows.Scan(&vs.Numero, &vs.Formatado); err != nil {
			log.Printf("Erro ao escanear a linha do versiculo: %v", err)
			fmt.Println(err)
			return models.LivroCapVers{}, err
		}
		versiculosSlimList = append(versiculosSlimList, vs)
	}

	if len(versiculosSlimList) == 0 {
		livroCapVers.Error = gin.H{
			"state":   true,
			"message": "Index was out of range. Must be non-negative and less than the size of the collection.",
		}
	} else {
		livroCapVers.Error = false
	}

	// Preencher o resto da estrutura de resultado
	livroCapVers.LivroID, _ = strconv.Atoi(livroID)
	livroCapVers.CapituloID, _ = strconv.Atoi(capituloID)
	livroCapVers.VersiculosList = versiculosSlimList

	defer rows.Close()
	queryVersiculos.Close()

	return livroCapVers, nil
}

func (vr *VersiculoRepository) GetVersiculoByLivroIdAndCapituloIdAndVersiculoId(id_livro int, id_capitulo int, id_versiculo int) (*models.Versiculo, error) {

	var versiculo models.Versiculo
	var err error

	if id_livro <= 0 || id_capitulo <= 0 || id_versiculo <= 0 {
		return nil, err
	}

	query, err := vr.connection.Prepare(`
		SELECT id, livro_id, capitulo_id, versao, numero, formatado 
		FROM versiculo 
		WHERE livro_id = ? AND capitulo_id = ? AND id = ?`)

	if err != nil {
		log.Printf("Erro ao preparar a query de versiculo: %v", err)
		return nil, err
	}

	err = query.QueryRow(id_versiculo).Scan(
		&versiculo.ID,
		&versiculo.LivroID,
		&versiculo.CapituloID,
		&versiculo.Numero,
		&versiculo.Formatado)

	if err != nil {
		log.Printf("Erro ao escanear a linha do versiculo: %v", err)
		return nil, err
	}

	return &versiculo, nil
}
