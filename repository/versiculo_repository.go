package repository

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/JimmyArgote/biblia-api-go/models"
)

type VersiculoRepository struct {
	connection *sql.DB
}

func NewVersiculoRepository(connection *sql.DB) VersiculoRepository {
	return VersiculoRepository{
		connection: connection,
	}
}

func (vr *VersiculoRepository) ListByChapter(id_livro, id_capitulo int) (models.LivroCapVers, error) {

	var livroCapVers models.LivroCapVers
	var err error

	livroCapVers.Error = false
	livroCapVers.Message = ""

	// DEBUG: Adicione esta linha!
	log.Printf("DEBUG: Recebido em ListByChapter -> id_capitulo: %d, id_livro: %d", id_capitulo, id_livro)

	// 1. Obter contagens e informações do livro em uma única consulta otimizada.
	query1 := `
		SELECT
			l.nome,
			l.sigla,
			l.testamento,
			(SELECT COUNT(1) FROM versiculo v WHERE v.livro_id = l.id AND v.capitulo_id = ?) as qtd_vers,
			(SELECT COUNT(1) FROM capitulo c WHERE c.livro_id = l.id) as qtd_caps
		FROM livro l
		WHERE l.id = ?`

	err = vr.connection.QueryRow(query1, id_capitulo, id_livro).Scan(
		&livroCapVers.LivroNome, &livroCapVers.LivroSigla, &livroCapVers.Testamento, &livroCapVers.VersTotal, &livroCapVers.CapsTotal,
	)

	if err != nil {
		// 3. Se o livro não for encontrado, simplesmente retorne o erro.
		// O handler irá traduzir `sql.ErrNoRows` para um 404 Not Found.
		// Não retorne um modelo parcialmente preenchido.
		return models.LivroCapVers{
			Error:   true,
			Message: "Livro inexistente",
		}, err
	}

	livroID := strconv.Itoa(id_livro)
	capituloID := strconv.Itoa(id_capitulo)

	if id_capitulo > livroCapVers.CapsTotal {
		return models.LivroCapVers{
			Error:   true,
			Message: "Capítulo inexistente",
		}, err
	}

	// 2. Obter a lista de versículos.
	query2 := `
		SELECT numero, formatado
		FROM versiculo
		WHERE livro_id = ? AND capitulo_id = ?
		ORDER BY numero ASC`

	rows, err := vr.connection.Query(query2, livroID, capituloID)

	if err != nil {
		log.Printf("Erro ao executar query de versículos: %v", err)
		livroCapVers.Error = true
		livroCapVers.Message = "Erro ao executar query de versículos"
		return livroCapVers, err
	}

	//defer rows.Close()

	var versiculosSlimList []models.VersiculoSlim

	for rows.Next() {

		var vs models.VersiculoSlim

		if err := rows.Scan(&vs.Numero, &vs.Formatado); err != nil {
			log.Printf("Erro ao escanear linha de versículo: %v", err)
			return models.LivroCapVers{}, err
		}

		versiculosSlimList = append(versiculosSlimList, vs)
	}

	if err = rows.Err(); err != nil {
		return models.LivroCapVers{}, err
	}

	if versiculosSlimList == nil {

		return models.LivroCapVers{
			Error:   true,
			Message: "Versículos não encontrados ou inexistentes",
		}, err

	}

	livroCapVers.Error = false

	// Preencher o resto da estrutura de resultado
	livroCapVers.LivroID, _ = strconv.Atoi(livroID)
	livroCapVers.CapituloID, _ = strconv.Atoi(capituloID)
	livroCapVers.VersiculosList = versiculosSlimList

	defer rows.Close()

	return livroCapVers, err
}

func (vr *VersiculoRepository) Find(id_livro int, id_capitulo int, id_versiculo int) (*models.Versiculo, error) {

	var versiculo models.Versiculo
	var err error

	query, err := vr.connection.Prepare(`
		SELECT id, livro_id, capitulo_id, versao_id, numero, formatado 
		FROM versiculo 
		WHERE livro_id = ? AND capitulo_id = ? AND id = ?`)

	if err != nil {
		log.Printf("Erro ao preparar a query de versiculo: %v", err)
		return nil, err
	}

	err = query.QueryRow(id_livro, id_capitulo, id_versiculo).Scan(
		&versiculo.ID,
		&versiculo.LivroID,
		&versiculo.CapituloID,
		&versiculo.VersaoID,
		&versiculo.Numero,
		&versiculo.Formatado)

	if err != nil {
		log.Printf("Erro ao escanear a linha do versiculo: %v", err)
		return nil, err
	}

	return &versiculo, err
}
