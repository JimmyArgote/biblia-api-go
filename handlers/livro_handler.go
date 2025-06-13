package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/biblia-api-go/database"
	"github.com/seu-usuario/biblia-api-go/models"
)

// ListarLivros busca e retorna todos os livros da Bíblia.
func ListarLivros(c *gin.Context) {
	// Sintaxe SQL ajustada para MySQL (sem colchetes e sem [dbo]).
	sql := `
		SELECT id, ordem, nome, sigla, testamento
		FROM livro
		ORDER BY testamento ASC, ordem ASC`

	rows, err := database.DB.Query(sql)
	if err != nil {
		log.Printf("Erro ao executar a query de livros: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar livros"})
		return
	}
	defer rows.Close()

	livros := []models.Livro{}
	for rows.Next() {
		var livro models.Livro
		if err := rows.Scan(&livro.ID, &livro.Ordem, &livro.Nome, &livro.Sigla, &livro.Testamento); err != nil {
			log.Printf("Erro ao escanear a linha do livro: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar dados dos livros"})
			return
		}
		livros = append(livros, livro)
	}

	c.JSON(http.StatusOK, livros)
}
