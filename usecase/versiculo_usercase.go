package usecase

import (
	"github.com/JimmyArgote/biblia-api-go/models"
	"github.com/JimmyArgote/biblia-api-go/repository"
)

type VersiculoUseCase struct {
	repository repository.VersiculoRepository
}

func NewVersiculoUseCase(repo repository.VersiculoRepository) VersiculoUseCase {
	return VersiculoUseCase{
		repository: repo,
	}
}

func (vu *VersiculoUseCase) GetVersiculosByCapituloIdAndLivroId(id_capitulo int, id_livro int) (models.LivroCapVers, error) {
	return vu.repository.GetVersiculosByCapituloIdAndLivroId(id_capitulo, id_livro)
}

func (vu *VersiculoUseCase) GetVersiculoByLivroIdAndCapituloIdAndVersiculoId(id_livro int, id_capitulo int, id_versiculo int) (*models.Versiculo, error) {
	return vu.repository.GetVersiculoByLivroIdAndCapituloIdAndVersiculoId(id_livro, id_capitulo, id_versiculo)
}
