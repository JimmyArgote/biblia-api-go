package controller

import (
	"net/http"
	"strconv"

	"github.com/JimmyArgote/biblia-api-go/models"
	"github.com/JimmyArgote/biblia-api-go/usecase"
	"github.com/gin-gonic/gin"
)

type livroController struct {
	livroUseCase usecase.LivroUseCase
}

func NewLivroController(usecase usecase.LivroUseCase) livroController {
	return livroController{
		livroUseCase: usecase,
	}
}

func (lc *livroController) GetLivros(c *gin.Context) {
	livros, err := lc.livroUseCase.GetLivros()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar livros"})
		return
	}

	c.JSON(http.StatusOK, livros)
}

func (lc *livroController) GetLivroByID(c *gin.Context) {

	id := c.Param("livroId")
	if id == "" {
		response := models.Response{
			Message: "Id do Livro nao pode ser nulo",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	livroId, err := strconv.Atoi(id)
	if err != nil {
		response := models.Response{
			Message: "Id do Livro precisa ser um numero",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	livro, err := lc.livroUseCase.GetLivroById(livroId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if livro == nil {
		response := models.Response{
			Message: "Livro nao foi encontrado na base de dados",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	c.JSON(http.StatusOK, livro)

}
