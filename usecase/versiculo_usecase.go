package usecase

import (
	"github.com/JimmyArgote/biblia-api-go/models"
	"github.com/JimmyArgote/biblia-api-go/repository"
)

type VersiculoUseCaseInterface interface {
	ListByChapter(id_livro int, id_capitulo int) (models.LivroCapVers, error)
	Find(id_livro int, id_capitulo int, id_versiculo int) (*models.Versiculo, error)
}

type VersiculoUseCase struct {
	repository repository.VersiculoRepository
}

func NewVersiculoUseCase(repo repository.VersiculoRepository) VersiculoUseCase {
	return VersiculoUseCase{
		repository: repo,
	}
}

func (vu *VersiculoUseCase) ListByChapter(id_livro int, id_capitulo int) (models.LivroCapVers, error) {

	capitulo, err := vu.repository.ListByChapter(id_livro, id_capitulo)

	if err != nil {
		return models.LivroCapVers{}, err
	}

	return capitulo, nil
}

func (vu *VersiculoUseCase) Find(id_livro int, id_capitulo int, id_versiculo int) (*models.Versiculo, error) {

	versiculo, err := vu.repository.Find(id_livro, id_capitulo, id_versiculo)

	if err != nil {
		return nil, err
	}

	return versiculo, nil
}
