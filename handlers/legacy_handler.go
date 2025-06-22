package handlers

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LegacyIndexHandler replica o comportamento do método Index(livro, capitulo, versiculo)
// do LivrosController.cs, que direciona a requisição com base nos parâmetros de query.
func LegacyIndexHandler(c *gin.Context) {
	// O C# usa parâmetros com 'L', 'C', 'V' maiúsculos.
	// c.DefaultQuery retorna "0" se o parâmetro não for encontrado.
	livroStr := c.DefaultQuery("Livro", "0")
	capituloStr := c.DefaultQuery("Capitulo", "0")
	versiculoStr := c.DefaultQuery("Versiculo", "0")

	livroID, _ := strconv.Atoi(livroStr)
	capituloID, _ := strconv.Atoi(capituloStr)
	versiculoID, _ := strconv.Atoi(versiculoStr)

	// GET livros
	if livroID == 0 && capituloID == 0 && versiculoID == 0 {
		log.Println("Legacy: Redirecionando para ListarLivros")
		ListarLivros(c)
		return
	}

	// GET capitulos
	if capituloID == 0 && versiculoID == 0 {
		log.Printf("Legacy: Redirecionando para ListarCapitulos (Livro: %d)", livroID)
		// O handler ListarCapitulos já busca por ?livro=, então apenas passamos o contexto.
		// Para garantir, vamos usar o parâmetro em minúsculas que o handler espera.
		c.Request.URL.Query().Set("livro", livroStr)
		ListarCapitulos(c)
		return
	}

	// GET versiculos
	if versiculoID == 0 {
		log.Printf("Legacy: Redirecionando para ListarVersiculos (Livro: %d, Capitulo: %d)", livroID, capituloID)
		// O handler ListarVersiculos busca por ?livro= e ?capitulo=
		c.Request.URL.Query().Set("livro", livroStr)
		c.Request.URL.Query().Set("capitulo", capituloStr)
		ListarVersiculos(c)
		return
	}

	// GET versiculo único
	log.Printf("Legacy: Redirecionando para ListarVersiculoUnico (Livro: %d, Capitulo: %d, VersiculoID: %d)", livroID, capituloID, versiculoID)
	// O contexto `c` já contém os query params (`Livro`, `Capitulo`, `Versiculo`),
	// e o handler `ListarVersiculoUnico` foi feito para ler exatamente esses parâmetros.
	ListarVersiculoUnico(c)
}
