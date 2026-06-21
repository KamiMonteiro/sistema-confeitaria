package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sistema-confeitaria/model"
	"sistema-confeitaria/repository"
	"strconv"
	"strings"
)

// CriarPedido valida e salva um pedido com todos os seus itens numa transação.
// Exige pelo menos um item — pedido vazio não faz sentido.
// Retorna o ID do pedido criado pra que o front possa redirecionar p/ edição se precisar.
// Rota: POST /api/novo/pedido
func CriarPedido(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}
		var p model.Pedido
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		// Valida os vínculos obrigatórios — sem cliente e sem forma de pagamento não salva
		if p.IDCliente == 0 || p.IDFormaPagamento == 0 {
			http.Error(w, "Cliente e forma de pagamento obrigatórios", http.StatusBadRequest)
			return
		}
		if len(p.Itens) == 0 {
			http.Error(w, "Pedido deve ter ao menos um item", http.StatusBadRequest)
			return
		}
		id, err := repository.CriarPedido(db, &p)
		if err != nil {
			http.Error(w, "Erro ao salvar pedido", http.StatusInternalServerError)
			return
		}
		// Retorna o ID gerado — o front usa pra redirecionar após salvar
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int64{"id_pedido": id})
	}
}

// AtualizarPedido edita o cabeçalho do pedido e substitui todos os itens.
// O ID do pedido é obrigatório no corpo do JSON.
// Rota: PUT /api/atualizar/pedido
func AtualizarPedido(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}
		var p model.Pedido
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		if p.ID == 0 {
			http.Error(w, "ID obrigatório", http.StatusBadRequest)
			return
		}
		if err := repository.AtualizarPedido(db, &p); err != nil {
			http.Error(w, "Erro ao atualizar pedido", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// PedidoPorID retorna um pedido completo (cabeçalho + itens) pelo ID na URL.
// Usado na tela de edição pra carregar o pedido e seus itens no formulário.
// Rota: GET /api/pedido/listar/{id}
func PedidoPorID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extrai o ID da URL — ex: /api/pedido/listar/2 → id = 2
		idStr := strings.TrimPrefix(r.URL.Path, "/api/pedido/listar/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}
		p, err := repository.BuscarPedidoPorID(db, id)
		if err != nil {
			http.Error(w, "Pedido não encontrado", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}

// BuscarTodosPedidos lista todos os pedidos do mais recente pro mais antigo.
// Já vem com nome do cliente e da forma de pagamento pelo JOIN no repositório.
// Rota: GET /api/todos/pedido
func BuscarTodosPedidos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pedidos, err := repository.BuscarTodosPedidos(db)
		if err != nil {
			http.Error(w, "Erro ao buscar pedidos", http.StatusInternalServerError)
			return
		}
		// Garante [] em vez de null quando não há pedidos
		if pedidos == nil {
			pedidos = []model.Pedido{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pedidos)
	}
}

// BuscarPedidosComFiltro filtra pedidos pelo status de entrega — recebe ?status= na query.
// Ex: ?status=Pendente devolve todos os pedidos pendentes.
// Rota: GET /api/pedido/buscar?status=
func BuscarPedidosComFiltro(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		pedidos, err := repository.BuscarPedidosPorStatus(db, status)
		if err != nil {
			http.Error(w, "Erro ao buscar pedidos", http.StatusInternalServerError)
			return
		}
		if pedidos == nil {
			pedidos = []model.Pedido{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pedidos)
	}
}

// ExcluirPedido remove o pedido e todos os itens vinculados (numa transação no repositório).
// Rota: DELETE /api/pedido/excluir/{id}
func ExcluirPedido(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}
		idStr := strings.TrimPrefix(r.URL.Path, "/api/pedido/excluir/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}
		if err := repository.ExcluirPedido(db, id); err != nil {
			http.Error(w, "Erro ao excluir pedido", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent) // 204
	}
}
