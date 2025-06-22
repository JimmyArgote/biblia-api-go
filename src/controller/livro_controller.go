package controller

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/JimmyArgote/biblia-api-go/src/models"
	"github.com/JimmyArgote/biblia-api-go/src/usecase"
	"github.com/gin-gonic/gin"
)

type LivroController struct {
	livroUseCase usecase.LivroUseCase
}

func NewLivroController(usecase usecase.LivroUseCase) LivroController {
	return LivroController{
		livroUseCase: usecase,
	}
}

func (lc *LivroController) GetLivros(ctx *gin.Context) {

	livros, err := lc.livroUseCase.GetLivros()

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			// Se 'sql.ErrNoRows', significa que não encontrou nada.
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"Error": true, "Message": "Livros não encontrados."},
			)

		} else {
			// Para qualquer outro erro inesperado (ex: falha de conexão com o banco)
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"Error": true, "Message": "Ocorreu um erro interno."},
			)
		}

		return
	}

	ctx.JSON(http.StatusOK, livros)
}

func (lc *LivroController) GetLivroByID(ctx *gin.Context) {

	livroId := ctx.Param("livroId")

	response := models.Response{
		Message: "",
		Error:   true,
	}

	id_livro, err := strconv.Atoi(livroId)
	if err != nil {
		response.Message = "Id do Livro precisa ser um número."
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if id_livro <= 0 {
		response.Message = "O Id do Livro deve ser um número positivo."
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	livro, err := lc.livroUseCase.GetLivroById(id_livro)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			// Se 'sql.ErrNoRows', significa que não encontrou nada.
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"Error": true, "Message": "Livro não encontrado."},
			)

		} else {
			// Para qualquer outro erro inesperado (ex: falha de conexão com o banco)
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"Error": true, "Message": "Ocorreu um erro interno."},
			)
		}

		return
	}

	ctx.JSON(http.StatusOK, livro)
}
