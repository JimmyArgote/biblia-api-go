# Technical Documentation: biblia-api-go

## Arquitetura

O projeto adota uma **arquitetura em camadas** com injeção de dependências manual:

```
HTTP Request
    │
    ▼
┌──────────────┐     ┌──────────────┐
│  handlers/   │     │ controller/  │   ← Camada de apresentação
│  (legado)    │     │  (v2)        │
└──────┬───────┘     └──────┬───────┘
       │                    │
       │   (handlers legados acessam database.DB diretamente)
       │                    │
       │              ┌─────▼─────┐
       │              │  usecase/ │   ← Camada de negócio (thin pass-through)
       │              └─────┬─────┘
       │                    │
       │              ┌─────▼──────┐
       └──────────────│ repository/│   ← Camada de dados (SQL)
                      └─────┬──────┘
                            │
                      ┌─────▼──────┐
                      │  database/ │   ← Conexão MySQL (variável global)
                      └────────────┘
```

### Principais Decisões Arquiteturais

1. **Duas camadas de apresentação coexistem**: `handlers/` (legado, acesso direto ao banco) e `controller/` (v2, usa usecase → repository). Isso é intencional para manter compatibilidade enquanto migra gradualmente.

2. **UseCases são thin pass-through**: atualmente apenas delegam para o repository. Existem como pontos de extensão futura para validação, transformação ou regras de negócio.

3. **Injeção manual via Container**: `di/container.go` monta todo o grafo de dependências na inicialização. Não usa frameworks de DI.

4. **Variável global `database.DB`**: os handlers legados dependem dela; os repositories recebem a conexão por parâmetro no construtor.

## Stack Técnica

| Componente       | Biblioteca                  | Versão  |
| ---------------- | --------------------------- | ------- |
| HTTP Router      | `gin-gonic/gin`             | v1.10.1 |
| Driver MySQL     | `go-sql-driver/mysql`       | v1.9.3  |
| Config           | `joho/godotenv`             | v1.5.1  |
| Go Runtime       | stdlib                      | 1.22.3  |

### Por que Gin?
- Roteamento com path parameters (`:id`)
- Middleware de CORS nativo
- JSON binding (`ShouldBindJSON`)
- Performance (router baseado em radix tree)

### Por que database/sql direto (sem ORM)?
- Controle total sobre as queries SQL
- Mapeamento explícito (sem magic)
- O projeto original C# usava queries SQL brutas — manter a mesma abordagem facilita a migração

## Estrutura de Diretórios

```
src/
├── models/          # Structs Go que representam entidades e DTOs
├── database/        # Inicialização da conexão MySQL (InitDB)
├── repository/      # Queries SQL. Um repository por entidade
├── usecase/         # Orquestração. Um usecase por entidade
├── controller/      # Handlers HTTP da API v2. Um controller por entidade
├── handlers/        # Handlers HTTP legados e da API v1
└── di/              # Container de injeção de dependências
```

## Banco de Dados

### Tabelas

| Tabela      | Colunas principais                              |
| ----------- | ----------------------------------------------- |
| `livro`     | `id`, `ordem`, `nome`, `sigla`, `testamento`    |
| `capitulo`  | `id`, `livro_id`, `versao_id`, `titulo`         |
| `versiculo` | `id`, `livro_id`, `capitulo_id`, `versao_id`, `numero`, `formatado`, `texto` |

### Índices

- A busca textual (`/Search`) usa `MATCH(versiculo.texto) AGAINST(? IN NATURAL LANGUAGE MODE)`, o que requer um índice **FULLTEXT** na coluna `versiculo.texto`.

### Models de busca

- `PesquisaRequest`: `palavra` (string, obrigatório), `limite` (int, default 100 se ≤ 0), `offset` (int, default 0 se < 0)
- `PesquisaResponse`: `lista` ([]Versiculo, nunca null — retorna `[]` vazio), `total` (int, contagem total de resultados), `limite` (int), `offset` (int)

### Conexão

- Configurada via variável de ambiente `DB_CONNECTION_STRING`
- Formato: `usuario:senha@tcp(host:porta)/nome_banco?parseTime=true`
- Carregada do arquivo `.env` (não versionado)

## Convenções de Código

### Nomenclatura
- **Português** para nomes de domínio: `Livro`, `Capitulo`, `Versiculo`, `ListarCapitulos`
- **Português** para mensagens de erro e logs
- **Inglês** para termos técnicos Go: `Repository`, `UseCase`, `Controller`, `Handler`
- Structs: PascalCase (`LivroRepository`)
- Métodos: PascalCase (`GetLivros`, `ListByChapter`)
- Arquivos: snake_case com sufixo da camada (`livro_repository.go`)

### Padrões

1. **Um arquivo por entidade por camada**: `livro_repository.go`, `livro_usecase.go`, `livro_controller.go`
2. **Construtores `New*`**: toda struct exportada tem um construtor que recebe dependências
3. **Retorno de erros, não panic**: erros de banco são retornados e tratados nos controllers
4. **Validação de parâmetros no controller**: path/query params são validados e convertidos na camada de apresentação
5. **Respostas de erro padronizadas**: `{"Error": true, "Message": "..."}` ou `gin.H{"error": "..."}` (inconsistente entre handlers e controllers — a ser padronizado)

### Anti-padrões identificados

- **Duplicação handlers ↔ controllers**: `handlers/versiculo_handler.go` e `repository/versiculo_repository.go` têm lógica SQL duplicada
- **Variável global `database.DB`**: usada pelos handlers legados, deveria ser injetada
- **Formato de erro inconsistente**: handlers usam `{"error": "..."}`, controllers usam `{"Error": true, "Message": "..."}`
- **Repository `ListByChapter` faz mais que dados**: inclui validação de capítulo inexistente (deveria estar no usecase)

## Rotas e seus Handlers

### Rotas Legadas (`main.go` + `handlers/`)

| Rota                              | Handler                   | Arquivo                      |
| --------------------------------- | ------------------------- | ---------------------------- |
| `GET /`                           | `LegacyIndexHandler`      | `handlers/legacy_handler.go` |
| `GET /Livros`                     | `LegacyIndexHandler`      | `handlers/legacy_handler.go` |
| `GET /Livros/Index`               | `LegacyIndexHandler`      | `handlers/legacy_handler.go` |
| `GET /Livros/ListarCapitulos`     | `ListarCapitulos`         | `handlers/capitulo_handler.go`|
| `GET /Livros/ListarVersiculos`    | `ListarVersiculos`        | `handlers/versiculo_handler.go`|
| `GET /Livros/ListarVersiculo`     | `ListarVersiculoUnico`    | `handlers/versiculo_handler.go`|
| `GET /Livros/ListarVers`          | `ListarVersiculos`        | `handlers/versiculo_handler.go`|
| `POST /Search`                    | `Pesquisar`               | `handlers/versiculo_handler.go` — corpo: `{"palavra", "limite"?, "offset"?}`; resposta: `{"lista", "total", "limite", "offset"}` |
| `POST /Search/Index`              | `Pesquisar`               | (idêntico ao acima) |

### API v1 (`main.go` + `handlers/`)

| Rota                                                   | Handler                   |
| ------------------------------------------------------ | ------------------------- |
| `GET /api/livros`                                      | `ListarLivros`            |
| `GET /api/livros/:livro_id`                            | `ListarCapitulos`         |
| `GET /api/livros/:livro_id/:capitulo_id`               | `ListarVersiculos`        |
| `GET /api/livros/:livro_id/:capitulo_id/:numero_versiculo` | `ObterVersiculoPorNumero` |

### API v2 (`main.go` + `controller/`)

| Rota                                                   | Controller Method                          |
| ------------------------------------------------------ | ------------------------------------------ |
| `GET /api/v2/livros`                                   | `LivroController.GetLivros`                |
| `GET /api/v2/livro/:livroId`                           | `LivroController.GetLivroByID`             |
| `GET /api/v2/capitulos/:livroId`                       | `CapituloController.GetCapitulosByLivroId` |
| `GET /api/v2/capitulo/:livroId/:capituloId`            | `CapituloController.GetCapituloByLivroIdAndCapituloId` |
| `GET /api/v2/versiculos/:livroId/:capituloId`          | `VersiculoController.ListByChapter`        |
| `GET /api/v2/versiculo/:livroId/:capituloId/:versiculoId` | `VersiculoController.Find`              |

## Fluxo de Inicialização

1. `main.go` chama `di.NewContainer()`
2. `NewContainer` carrega `.env`, lê `DB_CONNECTION_STRING`
3. `database.InitDB(connStr)` abre o pool de conexões MySQL
4. Repository → UseCase → Controller são instanciados em cadeia
5. Rotas Gin são registradas com os handlers e controllers
6. `router.Run(":8081")` inicia o servidor

## Deploy & Infraestrutura

### Docker

- **Dockerfile**: build multi-stage — `golang:1.22-alpine` (compilação) → `alpine:3.20` (runtime). Binário compilado com `CGO_ENABLED=0` e `-ldflags="-s -w"` para imagem mínima. Porta exposta: 8081.
- **.dockerignore**: exclui `.git`, `conductor/`, `.env`, `*.md`, `Dockerfile`, `docker-compose.yml` do contexto de build.
- **docker-compose.yml**: dois serviços — `app` (Go API, 256MB) e `db` (MySQL 8.0, 512MB, volume `mysql_data`). Rede `internal` isolada para o banco; rede `proxy-bible-api` para roteamento externo. Healthcheck no MySQL (`mysqladmin ping`) garante que o app só inicia com o banco pronto.
- **Conexão no Docker**: `DB_CONNECTION_STRING` no compose usa `db` como host (nome do serviço MySQL), construída a partir de `${MYSQL_USER}` e `${MYSQL_PASSWORD}`. `GIN_MODE=release` é setado automaticamente.

### CI/CD

- **`.github/workflows/deploy.yml`**: dispara em push na branch `main`. Builda imagem Docker, pusha para GHCR (`ghcr.io/<repo>:latest`), conecta via SSH ao servidor de produção e executa `docker compose up -d --force-recreate` + `docker image prune -f`.
- **Secrets necessários no GitHub**: `GHCR_TOKEN`, `SSH_HOST`, `SSH_PORT`, `SSH_USER`, `SSH_KEY`.

## Considerações para Desenvolvimento

### Para adicionar uma nova entidade
1. Criar model em `src/models/`
2. Criar repository com queries SQL
3. Criar usecase (mesmo que thin pass-through)
4. Criar controller com validação de parâmetros e tratamento de erros
5. Registrar no `Container` e adicionar rota em `main.go`

### Para refatorar um handler legado para v2
1. Mover a lógica SQL para um repository
2. Criar usecase que chama o repository
3. Criar controller que chama o usecase
4. Adicionar ao `Container`
5. Adicionar rota `/api/v2/...` em `main.go`
6. Não remover o handler legado até que todos os consumidores migrem
