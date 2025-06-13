package models

// Versiculo representa a estrutura completa de um versículo, usada principalmente para busca.
type Versiculo struct {
	ID              int    `json:"id"`
	CapituloID      int    `json:"capitulo_id"`
	LivroID         int    `json:"livro_id"`
	VersaoID        int    `json:"versao_id"`
	Numero          int    `json:"numero"`
	Formatado       string `json:"formatado"`
	LivroNome       string `json:"livro_nome,omitempty"`
	LivroSigla      string `json:"livro_sigla,omitempty"`
	Capitulo        int    `json:"capitulo,omitempty"` // No C# é `Capitulo`, que vem de `cap.id AS capitulo`
	LivroTestamento string `json:"livro_testamento,omitempty"`
	QtdVers         int    `json:"qtd_vers,omitempty"`
	QtdCaps         int    `json:"qtd_caps,omitempty"`
}

// VersiculoSlim é uma versão leve de um versículo para listas.
type VersiculoSlim struct {
	Numero    int    `json:"numero"`
	Formatado string `json:"formatado"`
}

// LivroCapVers é a estrutura de dados complexa retornada ao listar os versículos de um capítulo.
// Corresponde ao retorno do método ListarVersiculos no C#.
type LivroCapVers struct {
	Error          interface{}     `json:"error"` // Pode ser `false` ou um objeto de erro
	LivroID        int             `json:"livro_id"`
	CapituloID     int             `json:"capitulo_id"`
	VersTotal      int             `json:"vers_total"`
	CapsTotal      int             `json:"caps_total"`
	LivroNome      string          `json:"livro_nome"`
	LivroSigla     string          `json:"livro_sigla"`
	Testamento     string          `json:"testamento"`
	VersiculosList []VersiculoSlim `json:"versiculosList"`
}
