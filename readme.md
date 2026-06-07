# biblia-api-go

API REST em Go para consulta à Bíblia Sagrada. Fornece endpoints para listar livros, capítulos e versículos, além de busca textual com MySQL Full-Text Search.

## Funcionalidades

- **Livros**: listar todos os livros da Bíblia, filtrar por ID
- **Capítulos**: listar capítulos de um livro, buscar capítulo específico
- **Versículos**: listar versículos de um capítulo, buscar versículo por número
- **Busca**: pesquisa de palavras no texto dos versículos via `MATCH ... AGAINST`
- **Três conjuntos de API**: legado (compatibilidade), RESTful v1 e RESTful v2 (arquitetura limpa)

## Stack

| Camada       | Tecnologia                          |
| ------------ | ----------------------------------- |
| Linguagem    | Go 1.22                             |
| HTTP Router  | Gin v1.10                           |
| Banco        | MySQL (go-sql-driver/mysql v1.9)    |
| Config       | godotenv v1.5                       |

## Pré-requisitos

- Go 1.22+
- MySQL 5.7+ (ou MariaDB 10.3+) com suporte a Full-Text Search
- Arquivo `.env` configurado

## Configuração

Crie um arquivo `.env` na raiz do projeto:

```env
DB_CONNECTION_STRING=usuario:senha@tcp(localhost:3306)/nome_do_banco?parseTime=true
```

## Execução

```bash
# Instalar dependências
go mod tidy

# Executar
go run main.go
```

O servidor inicia em `http://localhost:8081`.

## Estrutura do Projeto

```
biblia-api-go/
├── main.go                  # Ponto de entrada, rotas e CORS
├── go.mod / go.sum          # Dependências Go
├── .gitignore
├── .env                     # Variáveis de ambiente (não versionado)
└── src/
    ├── models/              # Structs de dados
    │   ├── livro.go
    │   ├── capitulo.go
    │   ├── versiculo.go
    │   └── response.go
    ├── database/            # Conexão MySQL
    │   └── database.go
    ├── repository/          # Camada de acesso a dados (SQL)
    │   ├── livro_repository.go
    │   ├── capitulo_repository.go
    │   └── versiculo_repository.go
    ├── usecase/             # Camada de lógica de negócio
    │   ├── livro_usecase.go
    │   ├── capitulo_usecase.go
    │   └── versiculo_usecase.go
    ├── controller/          # Handlers HTTP da API v2
    │   ├── livro_controller.go
    │   ├── capitulo_controller.go
    │   └── versiculo_controller.go
    ├── handlers/            # Handlers legados e da API v1
    │   ├── legacy_handler.go
    │   ├── livro_handler.go
    │   ├── capitulo_handler.go
    │   └── versiculo_handler.go
    └── di/                  # Injeção de dependências
        └── container.go
```

## API Reference

### Rotas Legadas (compatibilidade com frontend existente)

| Método | Rota                          | Query Params                              | Descrição                       |
| ------ | ----------------------------- | ----------------------------------------- | ------------------------------- |
| GET    | `/`                           | `?Livro=&Capitulo=&Versiculo=`            | Dispatcher (roteia internamente) |
| GET    | `/Livros`                     | `?Livro=&Capitulo=&Versiculo=`            | Idem                            |
| GET    | `/Livros/Index`               | `?Livro=&Capitulo=&Versiculo=`            | Idem                            |
| GET    | `/Livros/ListarCapitulos`     | `?livro=X`                                | Lista capítulos                 |
| GET    | `/Livros/ListarVersiculos`    | `?livro=X&capitulo=Y`                     | Lista versículos                |
| GET    | `/Livros/ListarVers`          | `?livro=X&capitulo=Y`                     | Idem (alias)                    |
| GET    | `/Livros/ListarVersiculo`     | `?Livro=X&Capitulo=Y&Versiculo=Z`         | Versículo único                 |
| POST   | `/Search`                     | `{"palavra": "termo"}`                   | Busca textual                   |
| POST   | `/Search/Index`               | `{"palavra": "termo"}`                   | Idem (alias)                    |

### API v1 (RESTful simples, usa handlers legados)

| Método | Rota                                              | Descrição               |
| ------ | ------------------------------------------------- | ----------------------- |
| GET    | `/api/livros`                                     | Lista todos os livros   |
| GET    | `/api/livros/:livro_id`                           | Capítulos de um livro   |
| GET    | `/api/livros/:livro_id/:capitulo_id`              | Versículos do capítulo  |
| GET    | `/api/livros/:livro_id/:capitulo_id/:numero_versiculo` | Versículo por número |

### API v2 (arquitetura em camadas)

| Método | Rota                                       | Descrição                       |
| ------ | ------------------------------------------ | ------------------------------- |
| GET    | `/api/v2/livros`                           | Lista todos os livros           |
| GET    | `/api/v2/livro/:livroId`                   | Livro por ID                    |
| GET    | `/api/v2/capitulos/:livroId`               | Capítulos de um livro           |
| GET    | `/api/v2/capitulo/:livroId/:capituloId`    | Capítulo específico             |
| GET    | `/api/v2/versiculos/:livroId/:capituloId`  | Versículos do capítulo          |
| GET    | `/api/v2/versiculo/:livroId/:capituloId/:versiculoId` | Versículo por ID     |

---

## Exemplos de Respostas

### `GET /api/v2/livros`

```json
[
  {
    "id": 1,
    "ordem": 1,
    "nome": "Gênesis",
    "sigla": "Gn",
    "testamento": "VT"
  }
]
```

### `GET /api/v2/versiculos/1/1`

```json
{
  "error": false,
  "message": "",
  "livro_id": 1,
  "capitulo_id": 1,
  "vers_total": 31,
  "caps_total": 50,
  "livro_nome": "Gênesis",
  "livro_sigla": "Gn",
  "testamento": "VT",
  "versiculosList": [
    { "numero": 1, "formatado": "No princípio, criou Deus..." }
  ]
}
```

### `POST /Search` — Body: `{"palavra": "amor"}`

```json
{
  "lista": [
    {
      "id": 123,
      "capitulo_id": 3,
      "livro_id": 43,
      "numero": 16,
      "formatado": "Porque Deus amou o mundo...",
      "livro_nome": "João",
      "livro_sigla": "Jo",
      "capitulo": 3
    }
  ]
}
```
