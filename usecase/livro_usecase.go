package usecase

import (
	"github.com/JimmyArgote/biblia-api-go/models"
	"github.com/JimmyArgote/biblia-api-go/repository"
)

type LivroUseCaseInterface interface {
	GetLivros() ([]models.Livro, error)
	GetLivroById(id_livro int) (*models.Livro, error)
}

type LivroUseCase struct {
	repository repository.LivroRepository
}

func NewLivroUseCase(repo repository.LivroRepository) LivroUseCase {
	return LivroUseCase{
		repository: repo,
	}
}

func (lu *LivroUseCase) GetLivros() ([]models.Livro, error) {

	livros, err := lu.repository.GetLivros()

	if err != nil {
		return nil, err
	}

	return livros, nil
}

func (lu *LivroUseCase) GetLivroById(id_livro int) (*models.Livro, error) {

	livro, err := lu.repository.GetLivroById(id_livro)

	if err != nil || livro == nil {
		return nil, err
	}

	return livro, nil
}
