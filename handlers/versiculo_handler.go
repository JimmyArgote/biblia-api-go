package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/JimmyArgote/biblia-api-go/database"
	"github.com/JimmyArgote/biblia-api-go/models"
	"github.com/gin-gonic/gin"
)

// ListarVersiculos busca os versículos de um capítulo e livro específicos.
// Refatorado para MySQL para ser mais eficiente e limpo.
func ListarVersiculos(c *gin.Context) {
	livroID := c.Param("livro_id")
	if livroID == "" {
		livroID = c.Query("livro")
	}
	capituloID := c.Param("capitulo_id")
	if capituloID == "" {
		capituloID = c.Query("capitulo")
	}

	if livroID == "" || capituloID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDs do livro e capítulo são obrigatórios"})
		return
	}

	var result models.LivroCapVers
	var err error

	// 1. Obter contagens e informações do livro em uma única consulta otimizada.
	queryInfo := `
		SELECT
			l.nome,
			l.sigla,
			l.testamento,
			(SELECT COUNT(1) FROM versiculo v WHERE v.livro_id = l.id AND v.capitulo_id = ?) as qtd_vers,
			(SELECT COUNT(1) FROM capitulo c WHERE c.livro_id = l.id) as qtd_caps
		FROM livro l
		WHERE l.id = ?`

	err = database.DB.QueryRow(queryInfo, capituloID, livroID).Scan(
		&result.LivroNome, &result.LivroSigla, &result.Testamento, &result.VersTotal, &result.CapsTotal,
	)
	if err != nil {
		log.Printf("Erro ao buscar informações do livro/capítulo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar informações do livro"})
		return
	}

	// 2. Obter a lista de versículos.
	queryVersiculos := `
		SELECT numero, formatado
		FROM versiculo
		WHERE livro_id = ? AND capitulo_id = ?
		ORDER BY numero ASC`

	rows, err := database.DB.Query(queryVersiculos, livroID, capituloID)
	if err != nil {
		log.Printf("Erro ao executar query de versículos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar versículos"})
		return
	}
	defer rows.Close()

	var versiculosSlim []models.VersiculoSlim
	for rows.Next() {
		var vs models.VersiculoSlim
		if err := rows.Scan(&vs.Numero, &vs.Formatado); err != nil {
			log.Printf("Erro ao escanear linha de versículo: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar dados dos versículos"})
			return
		}
		versiculosSlim = append(versiculosSlim, vs)
	}

	// Preencher o resto da estrutura de resultado
	result.LivroID, _ = strconv.Atoi(livroID)
	result.CapituloID, _ = strconv.Atoi(capituloID)
	result.VersiculosList = versiculosSlim

	if len(versiculosSlim) == 0 {
		result.Error = gin.H{
			"state":   true,
			"message": "Index was out of range. Must be non-negative and less than the size of the collection.",
		}
	} else {
		result.Error = false
	}

	c.JSON(http.StatusOK, result)
}

// ListarVersiculoUnico busca um único versículo.
func ListarVersiculoUnico(c *gin.Context) {
	livroID := c.Query("Livro")
	capituloID := c.Query("Capitulo")
	versiculoID := c.Query("Versiculo")

	if livroID == "" || capituloID == "" || versiculoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros 'Livro', 'Capitulo' e 'Versiculo' são obrigatórios"})
		return
	}

	// Sintaxe SQL ajustada e placeholders mudados para ?
	query := `
		SELECT id, livro_id, capitulo_id, versao_id, numero, formatado
		FROM versiculo
		WHERE livro_id = ? AND capitulo_id = ? AND id = ?`

	rows, err := database.DB.Query(query, livroID, capituloID, versiculoID)
	if err != nil {
		log.Printf("Erro na query de versículo único: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar versículo"})
		return
	}
	defer rows.Close()

	var versiculos []models.Versiculo
	for rows.Next() {
		var v models.Versiculo
		if err := rows.Scan(&v.ID, &v.LivroID, &v.CapituloID, &v.VersaoID, &v.Numero, &v.Formatado); err != nil {
			log.Printf("Erro ao escanear versículo único: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar dados do versículo"})
			return
		}
		versiculos = append(versiculos, v)
	}

	c.JSON(http.StatusOK, versiculos)
}

// ObterVersiculoPorNumero busca um único versículo pelo número, não pelo ID.
func ObterVersiculoPorNumero(c *gin.Context) {
	// Lendo os PARÂMETROS DE CAMINHO com c.Param()
	livroID := c.Param("livro_id")
	capituloID := c.Param("capitulo_id")
	numeroVersiculo := c.Param("numero_versiculo") // O número do versículo

	if livroID == "" || capituloID == "" || numeroVersiculo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDs do livro, capítulo e o número do versículo são obrigatórios no caminho da URL."})
		return
	}

	// Query corrigida para buscar pelo NÚMERO do versículo, não pelo ID.
	query := `
		SELECT id, livro_id, capitulo_id, versao_id, numero, formatado
		FROM versiculo
		WHERE livro_id = ? AND capitulo_id = ? AND numero = ?` // <- MUDANÇA IMPORTANTE: numero = ?

	// Usamos QueryRow pois esperamos no máximo 1 resultado.
	row := database.DB.QueryRow(query, livroID, capituloID, numeroVersiculo)

	var v models.Versiculo
	// Usamos Scan diretamente no 'row'.
	err := row.Scan(&v.ID, &v.LivroID, &v.CapituloID, &v.VersaoID, &v.Numero, &v.Formatado)
	if err != nil {
		// Se 'sql.ErrNoRows', significa que não encontrou nada.
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Versículo não encontrado"})
			return
		}
		// Outros erros são erros de servidor.
		log.Printf("Erro ao escanear versículo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar dados do versículo"})
		return
	}

	// Retorna um único objeto JSON, não um array.
	c.JSON(http.StatusOK, v)
}

// Pesquisar busca por uma palavra no texto dos versículos via POST.
func Pesquisar(c *gin.Context) {
	// 1. Crie uma variável do tipo da sua nova struct.
	var requestBody models.PesquisaRequest

	// 2. Use ShouldBindJSON para preencher a struct com os dados do corpo da requisição.
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Corpo da requisição inválido", "details": err.Error()})
		return
	}

	// 3. Agora, use o valor da struct.
	palavra := requestBody.Palavra

	if palavra == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "insira uma palavra chave para fazer sua busca!"})
		return
	}

	// Sintaxe SQL ajustada, usando CONCAT para o LIKE no MySQL.
	query := `
		SELECT DISTINCT vers.id, vers.capitulo_id, vers.livro_id, vers.numero, vers.formatado,
		livro.nome AS livro_nome, livro.sigla AS livro_sigla, cap.id AS capitulo
		FROM versiculo vers
		INNER JOIN livro ON vers.livro_id = livro.id
		INNER JOIN capitulo cap ON vers.capitulo_id = cap.id
		WHERE MATCH(vers.texto) AGAINST(? IN NATURAL LANGUAGE MODE)
		ORDER BY vers.livro_id ASC, vers.capitulo_id ASC, vers.numero ASC`

	rows, err := database.DB.Query(query, palavra)
	if err != nil {
		log.Printf("Erro na query de pesquisa: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao realizar pesquisa"})
		return
	}
	defer rows.Close()

	var resultados []models.Versiculo
	for rows.Next() {
		var v models.Versiculo
		if err := rows.Scan(&v.ID, &v.CapituloID, &v.LivroID, &v.Numero, &v.Formatado, &v.LivroNome, &v.LivroSigla, &v.Capitulo); err != nil {
			log.Printf("Erro ao escanear resultado da pesquisa: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar resultado da pesquisa"})
			return
		}
		resultados = append(resultados, v)
	}

	c.JSON(http.StatusOK, gin.H{
		"lista": resultados,
	})
}
