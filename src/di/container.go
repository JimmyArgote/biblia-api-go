package di

import (
	"log"
	"os"

	"github.com/JimmyArgote/biblia-api-go/src/controller"
	"github.com/JimmyArgote/biblia-api-go/src/database"
	"github.com/JimmyArgote/biblia-api-go/src/repository"
	"github.com/JimmyArgote/biblia-api-go/src/usecase"
	"github.com/joho/godotenv"
)

// Container detém todas as dependências da aplicação.
// É uma maneira centralizada de construir e acessar componentes como controllers.
type Container struct {
	// Controllers - são os pontos de entrada para as rotas
	LivroController     controller.LivroController
	CapituloController  controller.CapituloController
	VersiculoController controller.VersiculoController
	// Você pode adicionar outras dependências aqui se necessário (ex: serviços de log, etc.)
}

// NewContainer é a "fábrica" que constrói toda a aplicação.
// Ele inicializa o banco de dados e monta todas as camadas em ordem.
func NewContainer() *Container {
	// 1. Inicializar a dependência fundamental: o banco de dados
	// Idealmente, a string de conexão viria de variáveis de ambiente
	// Carregar variáveis de ambiente do arquivo .env
	err := godotenv.Load()

	if err != nil {
		log.Println("Aviso: Arquivo .env não encontrado, usando variáveis de ambiente do sistema.")
	}

	// Obter a string de conexão
	connStr := os.Getenv("DB_CONNECTION_STRING")

	if connStr == "" {
		log.Fatal("A variável de ambiente DB_CONNECTION_STRING não está definida.")
	}

	// Inicializar a conexão com o banco de dados
	database.InitDB(connStr)

	// 2. Construir as dependências para a entidade "Livro"
	livroRepository := repository.NewLivroRepository(database.DB)
	livroUseCase := usecase.NewLivroUseCase(livroRepository)
	livroController := controller.NewLivroController(livroUseCase)

	// 3. Construir as dependências para a entidade "Capitulo"
	capituloRepository := repository.NewCapituloRepository(database.DB)
	capituloUseCase := usecase.NewCapituloUseCase(capituloRepository)
	capituloController := controller.NewCapituloController(capituloUseCase)

	// 4. Construir as dependências para a entidade "Versiculo"
	versiculoRepository := repository.NewVersiculoRepository(database.DB)
	versiculoUseCase := usecase.NewVersiculoUseCase(versiculoRepository)
	versiculoController := controller.NewVersiculoController(versiculoUseCase)

	// 5. Retornar o container com todos os componentes prontos para uso
	return &Container{
		LivroController:     livroController,
		CapituloController:  capituloController,
		VersiculoController: versiculoController,
	}
}
