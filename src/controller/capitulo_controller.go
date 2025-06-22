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

type CapituloController struct {
	capituloUseCase usecase.CapituloUseCase
}

func NewCapituloController(usecase usecase.CapituloUseCase) CapituloController {
	return CapituloController{
		capituloUseCase: usecase,
	}
}

func (cc *CapituloController) GetCapitulosByLivroId(ctx *gin.Context) {

	livroId := ctx.Param("livroId")

	response := models.Response{
		Message: "",
		Error:   true,
	}

	// Convertendo string para int
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

	capitulos, err := cc.capituloUseCase.GetCapitulosByLivroId(id_livro)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			// Se 'sql.ErrNoRows', significa que não encontrou nada.
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"Error": true, "Message": "Capítulos não encontrados."},
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

	ctx.JSON(http.StatusOK, capitulos)
}

func (cc *CapituloController) GetCapituloByLivroIdAndCapituloId(ctx *gin.Context) {

	// Lendo os parâmetros de caminho
	livroId := ctx.Param("livroId")
	capituloId := ctx.Param("capituloId")

	response := models.Response{
		Message: "",
		Error:   true,
	}

	// Convertendo string para int
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

	// Convertendo string para int
	id_capitulo, err := strconv.Atoi(capituloId)
	if err != nil {
		response.Message = "Id do Capítulo precisa ser um número."
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if id_capitulo <= 0 {
		response.Message = "O Id do Capítulo deve ser um número positivo."
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	capitulo, err := cc.capituloUseCase.GetCapituloByLivroIdAndCapituloId(id_livro, id_capitulo)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			// Se 'sql.ErrNoRows', significa que não encontrou nada.
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"Error": true, "Message": "Capítulo não encontrado."},
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

	ctx.JSON(http.StatusOK, capitulo)
}
