package main

import (
	"log"
	"os"

	"github.com/JimmyArgote/biblia-api-go/controller"
	"github.com/JimmyArgote/biblia-api-go/database"
	"github.com/JimmyArgote/biblia-api-go/handlers"
	"github.com/JimmyArgote/biblia-api-go/repository"
	"github.com/JimmyArgote/biblia-api-go/usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
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

	//Camada de repository
	LivroRepository := repository.NewLivroRepository(database.DB)

	//Camada usecase
	LivroUseCase := usecase.NewLivroUseCase(LivroRepository)

	//Camada de controllers
	LivroController := controller.NewLivroController(LivroUseCase)

	// Configurar o roteador Gin
	router := gin.Default()

	// Configurar CORS (Cross-Origin Resource Sharing) para permitir requisições do seu frontend
	router.Use(func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// --- Definindo as rotas da API para corresponder ao projeto ASP.NET ---

	// Rotas do LivrosController.cs
	// O endpoint "faz-tudo"
	router.GET("/", handlers.LegacyIndexHandler)
	router.GET("/Livros", handlers.LegacyIndexHandler)
	router.GET("/Livros/Index", handlers.LegacyIndexHandler)

	// Rotas explícitas que também existiam no controller
	router.GET("/Livros/ListarCapitulos", handlers.ListarCapitulos)      // Espera ?livro=X
	router.GET("/Livros/ListarVersiculos", handlers.ListarVersiculos)    // Espera ?livro=X&capitulo=Y
	router.GET("/Livros/ListarVersiculo", handlers.ListarVersiculoUnico) // Espera ?Livro=X&Capitulo=Y&Versiculo=Z

	// Rota específica para `ListarVers`
	router.GET("/Livros/ListarVers", handlers.ListarVersiculos) // Espera ?livro=X&capitulo=Y

	// Rotas do SearchController.cs
	// O frontend faz um POST para essas rotas
	router.POST("/Search", handlers.Pesquisar)
	router.POST("/Search/Index", handlers.Pesquisar)

	// (Opcional) Manter as rotas RESTful mais limpas para uso futuro ou novos frontends
	api := router.Group("/api")
	{
		api.GET("/livros", handlers.ListarLivros)
		api.GET("/livros/:livro_id", handlers.ListarCapitulos)
		api.GET("/livros/:livro_id/:capitulo_id", handlers.ListarVersiculos)
		api.GET("/livros/:livro_id/:capitulo_id/:numero_versiculo", handlers.ObterVersiculoPorNumero)
	}

	apiv2 := router.Group("/api/v2")
	{
		apiv2.GET("/livros", LivroController.GetLivros)
		apiv2.GET("/livro/:livroId", LivroController.GetLivroByID)

		//apiv2.GET("/capitulos/:livroId", LivroController.GetCapitulosByLivroID)
		//apiv2.GET("/capitulo/:livroId/:capituloId", LivroController.GetCapituloByID)

		//apiv2.GET("/versiculos/:livroId/:capituloId", LivroController.GetVersiculosByCapituloID)
		//apiv2.GET("/versiculo/:livroId/:capituloId/:versiculoId", LivroController.GetVersiculoByID)
	}

	// Iniciar o servidor na porta 8081
	log.Println("Servidor iniciado em http://localhost:8081")

	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
