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

/*func Usuarios(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPost:
			// 👉 criar usuário
			CriarUsuario(db, w, r)

		case http.MethodGet:
			// 👉 listar usuários (opcional por enquanto)
			BuscarUsuario(db)(w, r)

		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	}
}*/

func CriarUsuario(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		var u model.Usuario

		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		// validação simples
		if u.Nome == "" || u.Email == "" || u.Senha == "" {
			http.Error(w, "Campos obrigatórios", http.StatusBadRequest)
			return
		}

		// salvar no banco
		err = repository.CriarUsuario(db, &u)
		if err != nil {
			http.Error(w, "Erro ao salvar", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func AtualizarUsuario(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPut {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		var u model.Usuario

		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, "Erro no JSON", http.StatusBadRequest)
			return
		}

		err = repository.AtualizarUsuario(db, &u)
		if err != nil {
			http.Error(w, "Erro ao atualizar", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func BuscarTodosUsuario(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		usuarios, err := repository.BuscarTodosUsuario(db)
		if err != nil {
			http.Error(w, "Erro ao buscar usuários", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(usuarios)
	}
}

func ExcluirUsuario(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/api/usuarios/excluir/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		err = repository.ExcluirUsuario(db, id)
		if err != nil {
			http.Error(w, "Erro ao excluir usuário", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func BuscarUsuarioPorID(db *sql.DB, id int) (*model.Usuario, error) {
	query := `
	SELECT id_usuario, nome_usuario, cpf, email_usuario
	FROM USUARIO
	WHERE id_usuario = ?
	`

	row := db.QueryRow(query, id)

	var u model.Usuario
	err := row.Scan(&u.ID, &u.Nome, &u.CPF, &u.Email)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func UsuarioPorID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// extrair ID da URL
		idStr := strings.TrimPrefix(r.URL.Path, "/api/usuarios/listar/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		// chamar repository
		user, err := repository.BuscarUsuarioPorID(db, id)
		if err != nil {
			http.Error(w, "Usuário não encontrado", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds struct {
			Email string `json:"email"`
			Senha string `json:"senha"`
		}

		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		user, err := repository.AutenticarUsuario(db, creds.Email, creds.Senha)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Credenciais inválidas"})
			return
		}

		// Retornar token dummy e dados do usuário
		response := map[string]interface{}{
			"token": "dummy-token-" + strconv.Itoa(user.ID),
			"user":  user,
		}

		json.NewEncoder(w).Encode(response)
	}
}

func CriarPagamento(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		var f model.FormaPagamento

		err := json.NewDecoder(r.Body).Decode(&f)
		if err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(f.Descricao) == "" {
			http.Error(w, "Descrição obrigatória", http.StatusBadRequest)
			return
		}

		f.Ativo, err = normalizeAtivo(f.Ativo)
		if err != nil {
			http.Error(w, "Status inválido", http.StatusBadRequest)
			return
		}

		err = repository.CriarFormaPagamento(db, &f)
		if err != nil {
			http.Error(w, "Erro ao salvar forma de pagamento", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func AtualizarPagamento(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		var f model.FormaPagamento
		err := json.NewDecoder(r.Body).Decode(&f)
		if err != nil {
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

		f.Ativo, err = normalizeAtivo(f.Ativo)
		if err != nil {
			http.Error(w, "Status inválido", http.StatusBadRequest)
			return
		}

		err = repository.AtualizarFormaPagamento(db, &f)
		if err != nil {
			http.Error(w, "Erro ao atualizar forma de pagamento", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

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

		err = repository.ExcluirFormaPagamento(db, id)
		if err != nil {
			http.Error(w, "Erro ao excluir forma de pagamento", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

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
