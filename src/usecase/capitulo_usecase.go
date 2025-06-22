package usecase

import (
	"github.com/JimmyArgote/biblia-api-go/models"
	"github.com/JimmyArgote/biblia-api-go/repository"
)

type CapituloUseCase struct {
	repository repository.CapituloRepository
}

func NewCapituloUseCase(repo repository.CapituloRepository) CapituloUseCase {
	return CapituloUseCase{
		repository: repo,
	}
}

func (cu *CapituloUseCase) GetCapitulosByLivroId(id_livro int) ([]models.Capitulo, error) {

	capitulos, err := cu.repository.GetCapitulosByLivroId(id_livro)

	if err != nil {
		return capitulos, err
	}

	return capitulos, nil
}

func (cu *CapituloUseCase) GetCapituloByLivroIdAndCapituloId(id_livro int, id_capitulo int) (*models.Capitulo, error) {

	capitulo, err := cu.repository.GetCapituloByLivroIdAndCapituloId(id_livro, id_capitulo)

	if err != nil {
		return nil, err
	}

	return capitulo, nil

}
