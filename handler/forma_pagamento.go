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

// Rota: POST /api/novo/pagamento
func CriarPagamento(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		var f model.FormaPagamento
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(f.Descricao) == "" {
			http.Error(w, "Descrição obrigatória", http.StatusBadRequest)
			return
		}

		// Normaliza o valor de ativo — aceita SIM/S/NAO/N e converte pra 'S' ou 'N'
		var err error
		f.Ativo, err = normalizeAtivo(f.Ativo)
		if err != nil {
			http.Error(w, "Status inválido", http.StatusBadRequest)
			return
		}

		if err := repository.CriarFormaPagamento(db, &f); err != nil {
			http.Error(w, "Erro ao salvar forma de pagamento", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// AtualizarPagamento edita a descrição e o status de uma forma de pagamento.
// Rota: PUT /api/atualizar/pagamento
func AtualizarPagamento(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		var f model.FormaPagamento
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		if f.ID == 0 {
			http.Error(w, "ID obrigatório", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(f.Descricao) == "" {
			http.Error(w, "Descrição obrigatória", http.StatusBadRequest)
			return
		}

		var err error
		f.Ativo, err = normalizeAtivo(f.Ativo)
		if err != nil {
			http.Error(w, "Status inválido", http.StatusBadRequest)
			return
		}

		if err := repository.AtualizarFormaPagamento(db, &f); err != nil {
			http.Error(w, "Erro ao atualizar forma de pagamento", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// ConsultarPagamento retorna uma forma de pagamento pelo ID na URL.
// Rota: GET /api/pagamento/listar/{id}
func ConsultarPagamento(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/api/pagamento/listar/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		f, err := repository.BuscarFormaPagamentoPorID(db, id)
		if err != nil {
			http.Error(w, "Forma de pagamento não encontrada", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(f)
	}
}

// BuscarPagamentoPorDescricao pesquisa formas de pagamento pelo nome — recebe ?descricao=.
// Rota: GET /api/pagamento/buscar?descricao=
func BuscarPagamentoPorDescricao(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filtro := r.URL.Query().Get("descricao")

		formas, err := repository.BuscarFormasPagamentoPorDescricao(db, filtro)
		if err != nil {
			http.Error(w, "Erro ao buscar tipo de pagamento", http.StatusInternalServerError)
			return
		}

		if formas == nil {
			formas = []model.FormaPagamento{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(formas)
	}
}

// ExcluirPagamento remove uma forma de pagamento pelo ID na URL.
// Rota: DELETE /api/pagamento/excluir/{id}
func ExcluirPagamento(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/api/pagamento/excluir/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		if err := repository.ExcluirFormaPagamento(db, id); err != nil {
			http.Error(w, "Erro ao excluir forma de pagamento", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent) // 204 — sucesso sem corpo de resposta
	}
}

// BuscarTodasFormasPagamento lista todas as formas de pagamento (ativas e inativas).
// Rota: GET /api/todos/pagamento
func BuscarTodasFormasPagamento(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		formas, err := repository.BuscarTodasFormasPagamento(db)
		if err != nil {
			http.Error(w, "Erro ao buscar formas de pagamento", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(formas)
	}
}

// normalizeAtivo converte os valores que vêm do front (SIM/NAO/S/N) pro padrão do banco ('S'/'N').
// Chamada por CriarPagamento e AtualizarPagamento antes de salvar.
func normalizeAtivo(value string) (string, error) {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "SIM", "S":
		return "S", nil
	case "NAO", "N", "NÃO":
		return "N", nil
	default:
		return "", errors.New("valor inválido")
	}
}
