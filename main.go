package main

import (
	"log"

	"github.com/JimmyArgote/biblia-api-go/src/di"
	"github.com/JimmyArgote/biblia-api-go/src/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	log.Println("Iniciando a API da Bíblia...")

	// 1. Construir toda a aplicação com uma única chamada!
	container := di.NewContainer()

	// Configurar o roteador Gin
	router := gin.Default()

	// Configurar CORS (Cross-Origin Resource Sharing) para permitir requisições do seu frontend
	router.Use(func(ctx *gin.Context) {

		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	})

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
		apiv2.GET("/livros", container.LivroController.GetLivros)
		apiv2.GET("/livro/:livroId", container.LivroController.GetLivroByID)

		apiv2.GET("/capitulos/:livroId", container.CapituloController.GetCapitulosByLivroId)
		apiv2.GET("/capitulo/:livroId/:capituloId", container.CapituloController.GetCapituloByLivroIdAndCapituloId)

		apiv2.GET("/versiculos/:livroId/:capituloId", container.VersiculoController.ListByChapter)
		apiv2.GET("/versiculo/:livroId/:capituloId/:versiculoId", container.VersiculoController.Find)
	}

	// Iniciar o servidor na porta 8081
	log.Println("Servidor iniciado em http://localhost:8081")

	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
