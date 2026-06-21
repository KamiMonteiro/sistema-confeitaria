# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Sistema de gerenciamento para a confeitaria **Rosas & Tortas**. Backend em Go com banco SQLite, frontend em HTML/CSS/JS puro sem framework. Desenvolvido para uso interno: cadastro de clientes, produtos, pedidos, formas de pagamento e usuários do sistema.

---

## Commands

```bash
# Rodar o servidor (a partir da raiz do projeto)
go run ./go/main.go
# Servidor sobe em http://localhost:8080

# Compilar binário
go build -o confeitaria ./go/main.go

# Atualizar dependências
go mod tidy
```

**Desenvolvimento com Live Server:** os HTMLs ficam em `html/` e são abertos pelo VS Code Live Server na porta 5500. O Live Server serve os arquivos estáticos; as chamadas de API vão para `http://localhost:8080`. O Go server também pode servir o frontend diretamente em `localhost:8080`, mas no desenvolvimento o Live Server é mais prático pelo hot reload.

---

## Arquitetura

```
go/main.go              → registra todas as rotas HTTP e inicializa o banco
  handler/*.go          → valida request, chama repositório, escreve response
    repository/*.go     → executa SQL contra o banco
      database/db.go    → abre conexão SQLite
      database/migrate.go → cria as tabelas se não existirem
model/modelagem.go      → structs Go com tags JSON (mapeamento API ↔ banco)
```

**CORS duplo em `main.go`:**
- `corsHandler` — envolve `http.HandlerFunc` (rotas de API) → libera `127.0.0.1:5500`
- `corsMiddleware` — envolve `http.Handler` (file server) → libera `localhost:5500`
Os dois são necessários porque o file server e os handlers de API usam tipos diferentes.

**Banco:** `./go/confeitaria.db` (SQLite). Criado automaticamente ao subir o servidor via `RunMigrations`. Nunca apagar o arquivo em produção — é o banco real.

---

## Database Schema

Tabelas criadas em ordem de dependência de FK:

```sql
USUARIO (
  id_usuario   INTEGER PK,
  nome_usuario VARCHAR(100),
  cpf          CHAR(14),        -- armazenado só com dígitos (11 chars)
  email_usuario VARCHAR(100),
  senha        VARCHAR(100)     -- texto plano, sem hash
)

FORMA_PAGAMENTO (
  id_forma_pagamento INTEGER PK,
  descricao          VARCHAR(100),
  ativo              CHAR(2) CHECK (ativo IN ('S','N'))
)

CLIENTE (
  id_cliente  INTEGER PK,
  nome        VARCHAR(100),
  cpf         CHAR(14),         -- armazenado só com dígitos (11 chars)
  telefone    VARCHAR(11),      -- armazenado só com dígitos (10 ou 11 chars)
  email       VARCHAR(100),
  endereco    VARCHAR(150),
  localizador VARCHAR(500)      -- link do Google Maps ou texto livre
)

PRODUTO (
  id_produto     INTEGER PK,
  nome           VARCHAR(100),
  preco_unitario NUMERIC(7,2),
  descricao      VARCHAR(100),
  categoria      VARCHAR(100),
  ativo          CHAR(2) CHECK (ativo IN ('S','N'))
)

PEDIDO (
  id_pedido          INTEGER PK,
  status_entrega     VARCHAR(100),  -- "Pendente", "Preparo", "Entregue", "Cancelado"
  dt_pedido          DATETIME,      -- formato YYYY-MM-DD (comparado com date('now') no dashboard)
  vlr_total_pedido   NUMERIC(7,2),
  dt_entrega         DATETIME,
  observacao_entrega VARCHAR(100),
  id_cliente         INTEGER FK → CLIENTE,
  id_forma_pagamento INTEGER FK → FORMA_PAGAMENTO,
  id_usuario         INTEGER FK → USUARIO
)

ITEM_PEDIDO (
  id_item        INTEGER PK,
  qtd_item       INTEGER,
  vlr_total_item NUMERIC(7,2),   -- preco_unitario × qtd_item (calculado no front)
  id_produto     INTEGER FK → PRODUTO,
  id_pedido      INTEGER FK → PEDIDO
)
```

---

## API Endpoints

Base URL: `http://localhost:8080/api/`. Todas as respostas são JSON. IDs em rotas com barra final (`/listar/`, `/excluir/`) — o Go usa prefix matching, então `/api/pedido/listar/5` cai no handler registrado em `/api/pedido/listar/`.

### Auth
| Método | Rota | Body / Query | Resposta |
|--------|------|-------------|----------|
| POST | `/api/auth/login` | `{"email":"","senha":""}` | `{"token":"dummy-token-{id}","user":{...}}` |

Após login o front salva `usuario_nome` e `usuario_email` no `localStorage` para exibir no sidebar.

### Usuários (`handler/auth.go`)
| Método | Rota | Body / Query | Observação |
|--------|------|-------------|------------|
| POST | `/api/novo/usuario` | `{"nome_usuario","cpf","email_usuario","senha"}` | CPF validado (11 dígitos) |
| GET | `/api/todos/usuario` | — | Lista completa |
| GET | `/api/usuarios/buscar?nome=` | query string | Busca por nome ou e-mail |
| GET | `/api/usuarios/listar/{id}` | — | Um usuário pelo ID |
| PUT | `/api/atualizar/usuarios` | `{"id_usuario","nome_usuario","cpf","email_usuario","senha"}` | Senha em branco = não altera |
| DELETE | `/api/usuarios/excluir/{id}` | — | — |

```bash
# Exemplo: criar usuário
curl -X POST http://localhost:8080/api/novo/usuario \
  -H "Content-Type: application/json" \
  -d '{"nome_usuario":"Maria","cpf":"12345678901","email_usuario":"maria@email.com","senha":"123"}'

# Buscar por nome
curl "http://localhost:8080/api/usuarios/buscar?nome=maria"
```

### Formas de Pagamento (`handler/auth.go`)
| Método | Rota | Body / Query | Observação |
|--------|------|-------------|------------|
| POST | `/api/novo/pagamento` | `{"descricao","ativo"}` | ativo: SIM/NAO/S/N → normalizado para S/N |
| GET | `/api/todos/pagamento` | — | Ativas e inativas |
| GET | `/api/pagamento/buscar?descricao=` | query string | Busca parcial |
| GET | `/api/pagamento/listar/{id}` | — | Uma forma pelo ID |
| PUT | `/api/atualizar/pagamento` | `{"id_forma_pagamento","descricao","ativo"}` | — |
| DELETE | `/api/pagamento/excluir/{id}` | — | 204 sem corpo |

```bash
# Criar forma de pagamento
curl -X POST http://localhost:8080/api/novo/pagamento \
  -H "Content-Type: application/json" \
  -d '{"descricao":"Pix","ativo":"SIM"}'

# Buscar por descrição
curl "http://localhost:8080/api/pagamento/buscar?descricao=pix"
```

### Clientes (`handler/cliente.go`)
| Método | Rota | Body / Query | Observação |
|--------|------|-------------|------------|
| POST | `/api/novo/cliente` | `{"nome","cpf","telefone","email","endereco","localizador"}` | Só nome é obrigatório |
| GET | `/api/todos/cliente` | — | Ordem alfabética |
| GET | `/api/cliente/buscar?nome=` | query string | Busca por nome ou e-mail |
| GET | `/api/cliente/listar/{id}` | — | Um cliente pelo ID |
| PUT | `/api/atualizar/cliente` | `{"id_cliente","nome",...}` | — |
| DELETE | `/api/cliente/excluir/{id}` | — | Falha se cliente tiver pedidos (FK) |

CPF e telefone são salvos **sem máscara** (só dígitos). A formatação `000.000.000-00` e `(XX) XXXXX-XXXX` é aplicada no front (form com máscara no `oninput`, lista com funções `formatarCPF()`/`formatarTelefone()`).

O campo `localizador` armazena o link de compartilhamento do Google Maps (ex: `https://maps.app.goo.gl/...`).

```bash
# Criar cliente
curl -X POST http://localhost:8080/api/novo/cliente \
  -H "Content-Type: application/json" \
  -d '{"nome":"Ana Silva","cpf":"12345678901","telefone":"11999998888","email":"ana@email.com","endereco":"Rua das Flores, 10","localizador":"https://maps.app.goo.gl/xyz"}'

# Buscar por nome ou e-mail
curl "http://localhost:8080/api/cliente/buscar?nome=ana"
```

### Produtos (`handler/produto.go`)
| Método | Rota | Body / Query | Observação |
|--------|------|-------------|------------|
| POST | `/api/novo/produto` | `{"nome","preco_unitario","descricao","categoria","ativo"}` | Preço > 0 obrigatório |
| GET | `/api/todos/produto` | — | Todos (ativos + inativos) |
| GET | `/api/produto/buscar?nome=` | query string | Busca por nome ou categoria |
| GET | `/api/produto/listar/{id}` | — | Um produto pelo ID |
| PUT | `/api/atualizar/produto` | `{"id_produto","nome","preco_unitario",...}` | — |
| DELETE | `/api/produto/excluir/{id}` | — | Falha se produto estiver em pedidos (FK) |

O campo `ativo` aceita SIM/NAO/S/N no request — o handler normaliza para S/N antes de salvar. Na listagem, o filtro de status é feito no front (a API retorna todos).

```bash
# Criar produto
curl -X POST http://localhost:8080/api/novo/produto \
  -H "Content-Type: application/json" \
  -d '{"nome":"Bolo de Chocolate","preco_unitario":85.00,"descricao":"Recheio brigadeiro","categoria":"Bolos","ativo":"SIM"}'

# Buscar por nome ou categoria
curl "http://localhost:8080/api/produto/buscar?nome=bolo"
```

### Pedidos (`handler/pedido.go`)
| Método | Rota | Body / Query | Observação |
|--------|------|-------------|------------|
| POST | `/api/novo/pedido` | `{"id_cliente","id_forma_pagamento","id_usuario","status_entrega","dt_pedido","dt_entrega","vlr_total_pedido","observacao_entrega","itens":[...]}` | Retorna `{"id_pedido":N}` |
| GET | `/api/todos/pedido` | — | Com nome do cliente e forma de pagamento (JOIN) |
| GET | `/api/pedido/buscar?status=` | query string | Filtro por status (busca parcial) |
| GET | `/api/pedido/listar/{id}` | — | Cabeçalho + itens com nome/preço do produto |
| PUT | `/api/atualizar/pedido` | mesmo formato do POST + `"id_pedido"` | Substitui todos os itens (delete + insert) |
| DELETE | `/api/pedido/excluir/{id}` | — | Remove itens primeiro (FK), depois o pedido |

Criar e atualizar pedido usam transação SQL — se algum item falhar, o pedido inteiro é desfeito.

```bash
# Criar pedido com dois itens
curl -X POST http://localhost:8080/api/novo/pedido \
  -H "Content-Type: application/json" \
  -d '{
    "id_cliente": 1,
    "id_forma_pagamento": 2,
    "id_usuario": 1,
    "status_entrega": "Pendente",
    "dt_pedido": "2026-06-21",
    "dt_entrega": "2026-06-25",
    "vlr_total_pedido": 170.00,
    "observacao_entrega": "Entregar pela manhã",
    "itens": [
      {"id_produto": 1, "qtd_item": 2, "vlr_total_item": 170.00}
    ]
  }'

# Filtrar por status
curl "http://localhost:8080/api/pedido/buscar?status=Pendente"
```

### Dashboard (`handler/dashboard.go`)
| Método | Rota | Body / Query | Resposta |
|--------|------|-------------|----------|
| GET | `/api/dashboard` | — | KPIs do dia + últimos 5 pedidos |

```json
{
  "pedidos_hoje": 3,
  "faturamento_hoje": 420.50,
  "produtos_ativos": 12,
  "total_clientes": 48,
  "ultimos_pedidos": [ { "id_pedido": 10, "nome_cliente": "Ana", ... } ]
}
```

```bash
curl http://localhost:8080/api/dashboard
```

---

## Telas e o que cada uma faz

| Arquivo HTML | Tela | O que exibe / faz |
|---|---|---|
| `login.html` | Login | Autentica via API; salva `usuario_nome` e `usuario_email` no `localStorage` |
| `dashboard.html` | Dashboard | Chama `/api/dashboard` e exibe 4 KPIs (pedidos hoje, faturamento, produtos ativos, clientes) + tabela dos últimos 5 pedidos |
| `usuarios-lista.html` | Lista de usuários | Tabela paginada; busca por nome/e-mail via `/api/usuarios/buscar?nome=`; modal de confirmação para exclusão |
| `usuario-form.html` | Cadastro/edição de usuário | Cria via POST ou edita via PUT; detecta modo pela query string `?modo=editar&id=N` |
| `pagamentos-lista.html` | Lista de formas de pagamento | Busca por descrição via API + filtro local por status (Ativo/Inativo) na mesma lista |
| `pagamento-form.html` | Cadastro/edição de pagamento | Select de ativo com valores SIM/NAO; no carregamento mapeia S→SIM/N→NAO para o select |
| `clientes-lista.html` | Lista de clientes | CPF e telefone exibidos com máscara (`formatarCPF`/`formatarTelefone`); avatar com iniciais em verde mint |
| `cliente-form.html` | Cadastro/edição de cliente | Máscara no `oninput` para CPF e telefone; salva só dígitos; campo localizador com botão "Abrir" que abre o link no Google Maps |
| `produtos-lista.html` | Lista de produtos | Busca por nome/categoria via API + filtros locais de categoria e status ativo/inativo |
| `produto-form.html` | Cadastro/edição de produto | Preço validado (>0); ativo mapeado S→SIM no carregamento |
| `pedidos-lista.html` | Lista de pedidos | Filtro por status via API (`?status=`); badges coloridos por status; data formatada de ISO para DD/MM/AAAA |
| `pedido-form.html` | Cadastro/edição de pedido | Layout dois colunas: formulário à esquerda, resumo sticky à direita; carrega clientes, pagamentos e produtos simultaneamente com `Promise.all`; só produtos ativos aparecem na seleção; `adicionarItem()` acumula qtd se produto já estiver na lista |

---

## Convenções

**Booleanos no banco:** `'S'`/`'N'` (nunca `true`/`false`). Campos de formulário usam `SIM`/`NAO`. O handler converte com `normalizeAtivo()` (pagamentos) e `normalizarAtivoProduto()` (produtos) antes de salvar.

**Usuário logado no sidebar:** após login, `login.html` salva `usuario_nome` e `usuario_email` no `localStorage`. O `app.js` lê esses valores no `DOMContentLoaded` e preenche `#sidebar-avatar`, `#sidebar-usuario` e `#sidebar-email` em todas as páginas automaticamente. Se o login for antigo (sem esses campos), cai no fallback `login-email`.

**CPF e telefone de clientes:** salvos só com dígitos no banco. A máscara é aplicada no `oninput` do formulário e nas funções `formatarCPF()`/`formatarTelefone()` da listagem.

**Pedido usa transação:** `CriarPedido` e `AtualizarPedido` em `repository/pedido.go` abrem `BEGIN/COMMIT/ROLLBACK`. Na atualização, a estratégia é deletar todos os `ITEM_PEDIDO` do pedido e inserir os novos — mais simples que fazer diff.

**Exclusão com FK:** excluir cliente falha se tiver pedidos; excluir produto falha se estiver em pedidos; excluir pedido precisa apagar os itens primeiro (feito na transação de `ExcluirPedido`).

**Senha:** armazenada e comparada em texto plano — sem hash implementado.

**Datas:** `dt_pedido` e `dt_entrega` chegam do front como `YYYY-MM-DD`. O dashboard compara com `date('now')` do SQLite. A listagem de pedidos formata para `DD/MM/AAAA` no front.

**Paginação:** implementada no front (JavaScript), não na API. A API sempre retorna a lista completa.

---

## Estrutura de arquivos relevantes

```
go/main.go                  → único ponto de entrada; registra 32 rotas
handler/
  auth.go                   → Login, CRUD Usuário, CRUD FormaPagamento, normalizeAtivo()
  cliente.go                → CRUD Cliente
  produto.go                → CRUD Produto, normalizarAtivoProduto()
  pedido.go                 → CRUD Pedido (POST retorna id_pedido no body)
  dashboard.go              → GET /api/dashboard
repository/
  usuario.go                → SQL de usuário; AtualizarUsuario não altera senha se vier vazia
  pagamento.go              → SQL de forma de pagamento
  cliente.go                → BuscarClientesPorNome busca em nome E email
  produto.go                → BuscarProdutosPorNome busca em nome E categoria
  pedido.go                 → transações; buscarItensDoPedido (não exportada) faz JOIN com PRODUTO
  dashboard.go              → DashboardStats struct + 4 queries + query dos últimos 5 pedidos
model/modelagem.go          → 6 structs com tags json; omitempty em campos de JOIN
database/
  db.go                     → abre SQLite em ./go/confeitaria.db
  migrate.go                → cria 6 tabelas em ordem de FK
html/
  _sidebar.html             → template de referência do sidebar (não incluído por JS/server-side)
js/app.js                   → showToast, abrirModal/fecharModal, formatarMoeda, formatarDataHoje,
                               DOMContentLoaded: marca item ativo + preenche usuário no sidebar
css/styles.css              → classes de componente: btn-primary, btn-mint, badge, card, tbl, toast, etc.
```
