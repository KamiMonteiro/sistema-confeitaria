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

// CriarCliente recebe os dados do novo cliente e salva no banco.
// Só o nome é obrigatório — os demais campos podem vir em branco.
// Rota: POST /api/novo/cliente
func CriarCliente(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}
		var c model.Cliente
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(c.Nome) == "" {
			http.Error(w, "Nome obrigatório", http.StatusBadRequest)
			return
		}
		if err := repository.CriarCliente(db, &c); err != nil {
			http.Error(w, "Erro ao salvar cliente", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated) // 201
	}
}

// AtualizarCliente edita os dados de um cliente existente.
// ID e nome são obrigatórios — o resto pode ser vazio.
// Rota: PUT /api/atualizar/cliente
func AtualizarCliente(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}
		var c model.Cliente
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		if c.ID == 0 || strings.TrimSpace(c.Nome) == "" {
			http.Error(w, "ID e nome obrigatórios", http.StatusBadRequest)
			return
		}
		if err := repository.AtualizarCliente(db, &c); err != nil {
			http.Error(w, "Erro ao atualizar cliente", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// ClientePorID retorna os dados de um cliente pelo ID na URL.
// Usado na tela de edição pra preencher o formulário.
// Rota: GET /api/cliente/listar/{id}
func ClientePorID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extrai o ID da URL — ex: /api/cliente/listar/3 → id = 3
		idStr := strings.TrimPrefix(r.URL.Path, "/api/cliente/listar/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}
		c, err := repository.BuscarClientePorID(db, id)
		if err != nil {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(c)
	}
}

// BuscarTodosClientes retorna todos os clientes sem filtro, em ordem alfabética.
// Também é chamada pelo form de pedido pra popular o select de clientes.
// Rota: GET /api/todos/cliente
func BuscarTodosClientes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientes, err := repository.BuscarTodosClientes(db)
		if err != nil {
			http.Error(w, "Erro ao buscar clientes", http.StatusInternalServerError)
			return
		}
		// Garante [] em vez de null quando não há registros
		if clientes == nil {
			clientes = []model.Cliente{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clientes)
	}
}

// BuscarClientesComFiltro pesquisa clientes por nome ou email — recebe ?nome= na query.
// Rota: GET /api/cliente/buscar?nome=
func BuscarClientesComFiltro(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filtro := r.URL.Query().Get("nome")
		clientes, err := repository.BuscarClientesPorNome(db, filtro)
		if err != nil {
			http.Error(w, "Erro ao buscar clientes", http.StatusInternalServerError)
			return
		}
		if clientes == nil {
			clientes = []model.Cliente{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clientes)
	}
}

// ExcluirCliente remove o cliente pelo ID na URL.
// Vai falhar se o cliente tiver pedidos vinculados — o banco barra pela FK.
// Rota: DELETE /api/cliente/excluir/{id}
func ExcluirCliente(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}
		idStr := strings.TrimPrefix(r.URL.Path, "/api/cliente/excluir/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}
		if err := repository.ExcluirCliente(db, id); err != nil {
			http.Error(w, "Erro ao excluir cliente", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent) // 204 — sem corpo na resposta
	}
}
