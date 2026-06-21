package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sistema-confeitaria/model"
	"sistema-confeitaria/repository"
	"strconv"
	"strings"
)

// CriarProduto valida e salva um novo produto no catálogo.
// Nome e preço são obrigatórios. O campo ativo vem como SIM/NAO e é normalizado pra S/N.
// Rota: POST /api/novo/produto
func CriarProduto(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}
		var p model.Produto
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(p.Nome) == "" {
			http.Error(w, "Nome obrigatório", http.StatusBadRequest)
			return
		}
		if p.PrecoUnitario <= 0 {
			http.Error(w, "Preço deve ser maior que zero", http.StatusBadRequest)
			return
		}
		// Normaliza o ativo antes de salvar — aceita SIM/NAO/S/N
		var err error
		p.Ativo, err = normalizarAtivoProduto(p.Ativo)
		if err != nil {
			http.Error(w, "Status inválido", http.StatusBadRequest)
			return
		}
		if err := repository.CriarProduto(db, &p); err != nil {
			http.Error(w, "Erro ao salvar produto", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated) // 201
	}
}

// AtualizarProduto edita os dados de um produto existente.
// ID e nome são obrigatórios. Ativo também passa pela normalização.
// Rota: PUT /api/atualizar/produto
func AtualizarProduto(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}
		var p model.Produto
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		if p.ID == 0 || strings.TrimSpace(p.Nome) == "" {
			http.Error(w, "ID e nome obrigatórios", http.StatusBadRequest)
			return
		}
		var err error
		p.Ativo, err = normalizarAtivoProduto(p.Ativo)
		if err != nil {
			http.Error(w, "Status inválido", http.StatusBadRequest)
			return
		}
		if err := repository.AtualizarProduto(db, &p); err != nil {
			http.Error(w, "Erro ao atualizar produto", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// ProdutoPorID retorna um produto pelo ID na URL.
// Usado na tela de edição pra preencher o formulário.
// Rota: GET /api/produto/listar/{id}
func ProdutoPorID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extrai o ID da URL — ex: /api/produto/listar/7 → id = 7
		idStr := strings.TrimPrefix(r.URL.Path, "/api/produto/listar/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}
		p, err := repository.BuscarProdutoPorID(db, id)
		if err != nil {
			http.Error(w, "Produto não encontrado", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}

// BuscarTodosProdutos lista todos os produtos em ordem alfabética (ativos e inativos).
// O filtro de status é feito no front — a API devolve tudo.
// Também usado pelo form de pedido pra popular o select de produtos.
// Rota: GET /api/todos/produto
func BuscarTodosProdutos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		produtos, err := repository.BuscarTodosProdutos(db)
		if err != nil {
			http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
			return
		}
		// Garante [] em vez de null quando não há registros
		if produtos == nil {
			produtos = []model.Produto{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(produtos)
	}
}

// BuscarProdutosComFiltro pesquisa por nome ou categoria — recebe ?nome= na query.
// Rota: GET /api/produto/buscar?nome=
func BuscarProdutosComFiltro(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filtro := r.URL.Query().Get("nome")
		produtos, err := repository.BuscarProdutosPorNome(db, filtro)
		if err != nil {
			http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
			return
		}
		if produtos == nil {
			produtos = []model.Produto{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(produtos)
	}
}

// ExcluirProduto remove o produto pelo ID na URL.
// Vai falhar se o produto estiver em algum pedido — o banco barra pela FK em ITEM_PEDIDO.
// Rota: DELETE /api/produto/excluir/{id}
func ExcluirProduto(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}
		idStr := strings.TrimPrefix(r.URL.Path, "/api/produto/excluir/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}
		if err := repository.ExcluirProduto(db, id); err != nil {
			http.Error(w, "Erro ao excluir produto", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent) // 204
	}
}

// normalizarAtivoProduto faz a mesma coisa que normalizeAtivo em auth.go.
// Existe separada aqui pra não criar dependência entre os handlers.
// Converte SIM/NAO/S/N → 'S' ou 'N' pro banco.
func normalizarAtivoProduto(value string) (string, error) {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "SIM", "S":
		return "S", nil
	case "NAO", "N", "NÃO":
		return "N", nil
	default:
		return "", errors.New("valor inválido")
	}
}
