package main

import (
	"log"
	"os"

	"github.com/JimmyArgote/biblia-api-go/database"
	"github.com/JimmyArgote/biblia-api-go/handlers"
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

	// Configurar o roteador Gin
	r := gin.Default()

	// Configurar CORS (Cross-Origin Resource Sharing) para permitir requisições do seu frontend
	r.Use(func(c *gin.Context) {
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
	r.GET("/", handlers.LegacyIndexHandler)
	r.GET("/Livros", handlers.LegacyIndexHandler)
	r.GET("/Livros/Index", handlers.LegacyIndexHandler)

	// Rotas explícitas que também existiam no controller
	r.GET("/Livros/ListarCapitulos", handlers.ListarCapitulos)      // Espera ?livro=X
	r.GET("/Livros/ListarVersiculos", handlers.ListarVersiculos)    // Espera ?livro=X&capitulo=Y
	r.GET("/Livros/ListarVersiculo", handlers.ListarVersiculoUnico) // Espera ?Livro=X&Capitulo=Y&Versiculo=Z

	// Rota específica para `ListarVers`
	r.GET("/Livros/ListarVers", handlers.ListarVersiculos) // Espera ?livro=X&capitulo=Y

	// Rotas do SearchController.cs
	// O frontend faz um POST para essas rotas
	r.POST("/Search", handlers.Pesquisar)
	r.POST("/Search/Index", handlers.Pesquisar)

	// (Opcional) Manter as rotas RESTful mais limpas para uso futuro ou novos frontends
	api := r.Group("/api")
	{
		api.GET("/livros", handlers.ListarLivros)
		api.GET("/livros/:livro_id/capitulos", handlers.ListarCapitulos)
		api.GET("/livros/:livro_id/capitulos/:capitulo_id/versiculos", handlers.ListarVersiculos)
	}

	// Iniciar o servidor na porta 8080
	log.Println("Servidor iniciado em http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
