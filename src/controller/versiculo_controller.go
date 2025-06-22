package controller

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/JimmyArgote/biblia-api-go/models"
	"github.com/JimmyArgote/biblia-api-go/usecase"
	"github.com/gin-gonic/gin"
)

type VersiculoController struct {
	versiculoUseCase usecase.VersiculoUseCase
}

func NewVersiculoController(usecase usecase.VersiculoUseCase) VersiculoController {
	return VersiculoController{
		versiculoUseCase: usecase,
	}
}

func (vc *VersiculoController) ListByChapter(ctx *gin.Context) {

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

	response.Error = false

	livroCapVers, err := vc.versiculoUseCase.ListByChapter(id_livro, id_capitulo)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			// Se 'sql.ErrNoRows', significa que não encontrou nada.
			ctx.JSON(http.StatusNotFound, gin.H{"Error": true, "Message": "Versículos não encontrados."})

		} else {
			// Para qualquer outro erro inesperado (ex: falha de conexão com o banco)
			ctx.JSON(http.StatusInternalServerError, gin.H{"Error": true, "Message": "Ocorreu um erro interno."})
		}

		return
	}

	ctx.JSON(http.StatusOK, livroCapVers)
}

func (vc *VersiculoController) Find(ctx *gin.Context) {

	// Lendo os parâmetros de caminho
	livroId := ctx.Param("livroId")
	capituloId := ctx.Param("capituloId")
	versiculoId := ctx.Param("versiculoId")

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

	id_capitulo, err := strconv.Atoi(capituloId)
	if err != nil {
		response.Message = "Id do Capítulo precisa ser um número."
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	id_versiculo, err := strconv.Atoi(versiculoId)
	if err != nil {
		response.Message = "Id do Versículo precisa ser um número."
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if id_livro <= 0 || id_capitulo <= 0 || id_versiculo <= 0 {
		response.Message = "O Id do Livro/Capítulo/Versículo devem ser um números positivos."
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response.Error = false

	versiculo, err := vc.versiculoUseCase.Find(id_livro, id_capitulo, id_versiculo)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			// Se 'sql.ErrNoRows', significa que não encontrou nada.
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"Error": true, "Message": "Versículo não encontrado."},
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

	ctx.JSON(http.StatusOK, versiculo)
}
