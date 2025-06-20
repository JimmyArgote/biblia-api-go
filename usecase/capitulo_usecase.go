package usecase

import (
	"github.com/JimmyArgote/biblia-api-go/models"
	"github.com/JimmyArgote/biblia-api-go/repository"
)

type CapituloUseCaseInterface interface {
	GetCapitulosByLivroId(id_livro int) ([]models.Capitulo, error)
	GetCapituloById(id_capitulo int) (*models.Capitulo, error)
}

type CapituloUseCase struct {
	repository repository.CapituloRepository
}

func NewCapituloUseCase(repo repository.CapituloRepository) CapituloUseCase {
	return CapituloUseCase{
		repository: repo,
	}
}

func (cu *CapituloUseCase) GetCapitulosByLivroId(id_livro int) ([]models.Capitulo, error) {
	return cu.repository.GetCapitulosByLivroId(id_livro)
}

func (cu *CapituloUseCase) GetCapituloByLivroIdAndCapituloId(id_livro int, id_capitulo int) (*models.Capitulo, error) {
	return cu.repository.GetCapituloByLivroIdAndCapituloId(id_livro, id_capitulo)
}
