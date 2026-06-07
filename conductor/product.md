# Product: biblia-api-go

## Visão Geral

**biblia-api-go** é uma API REST que serve o conteúdo da Bíblia Sagrada. O projeto é uma reescrita em Go de uma API originalmente desenvolvida em C# (ASP.NET), mantendo compatibilidade com o frontend existente enquanto introduz uma nova arquitetura interna mais limpa.

## Propósito

Fornecer acesso programático ao texto bíblico estruturado em livros, capítulos e versículos. A API suporta tanto consultas diretas (por livro/capítulo/versículo) quanto busca textual por palavras, servindo como backend para aplicações web e mobile de leitura e estudo da Bíblia.

## Público-Alvo

- Desenvolvedores de frontends bíblicos (web, mobile)
- Aplicações de estudo bíblico
- Integrações com sistemas eclesiásticos

## Funcionalidades Principais

1. **Catálogo de Livros**: listagem completa dos 66 livros (VT e NT) com ordem canônica
2. **Navegação Hierárquica**: Livro → Capítulos → Versículos
3. **Busca Textual**: pesquisa Full-Text nos versículos por palavra-chave
4. **Compatibilidade Legada**: rotas no formato do controller C# original (`/Livros/*`)
5. **API Versionada**: endpoints `/api/v2` com separação clara de responsabilidades

## Decisões de Produto

- **Porta 8081**: escolhida para não conflitar com outras APIs locais
- **Português**: nomes de rotas, mensagens de erro e código são em português por ser o idioma do público-alvo
- **CORS aberto** (`*`): configuração permissiva para desenvolvimento — deve ser restringido em produção
- **Três conjuntos de rotas simultâneos**: legado (compatibilidade), v1 (RESTful simples), v2 (arquitetura limpa). A v2 é o caminho recomendado para novos consumidores.

## Status Atual

O projeto está funcional e em produção. As rotas v2 são o padrão moderno; as rotas legadas existem apenas para compatibilidade com frontends que dependem do formato antigo.
