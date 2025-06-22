package handlers

import (
	"log"
	"net/http"

	"github.com/JimmyArgote/biblia-api-go/src/database"
	"github.com/JimmyArgote/biblia-api-go/src/models"
	"github.com/gin-gonic/gin"
)

// ListarCapitulos busca os capítulos de um livro específico.
func ListarCapitulos(c *gin.Context) {
	livroID := c.Param("livro_id")
	if livroID == "" {
		livroID = c.Query("livro")
	}

	if livroID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do livro é obrigatório"})
		return
	}

	// Sintaxe SQL ajustada e placeholder mudado de @p1 para ?
	sql := `
		SELECT id, livro_id, versao_id, titulo
		FROM capitulo
		WHERE livro_id = ?
		ORDER BY id`

	rows, err := database.DB.Query(sql, livroID)
	if err != nil {
		log.Printf("Erro ao executar a query de capítulos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar capítulos"})
		return
	}
	defer rows.Close()

	capitulos := []models.Capitulo{}
	for rows.Next() {
		var cap models.Capitulo
		if err := rows.Scan(&cap.ID, &cap.LivroID, &cap.VersaoID, &cap.Titulo); err != nil {
			log.Printf("Erro ao escanear a linha do capítulo: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar dados dos capítulos"})
			return
		}
		capitulos = append(capitulos, cap)
	}

	c.JSON(http.StatusOK, capitulos)
}
