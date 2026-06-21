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

// cpfNumeros retira tudo que não for dígito do CPF — útil pra validar independente de formatação.
func cpfNumeros(cpf string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, cpf)
}

// validarCPF verifica se o CPF tem exatamente 11 dígitos numéricos.
// Não valida os dígitos verificadores — só checa o tamanho.
func validarCPF(cpf string) bool {
	digits := cpfNumeros(cpf)
	return len(digits) == 11
}

// CriarUsuario recebe o JSON do novo usuário, valida e salva no banco.
// Rota: POST /api/novo/usuario
func CriarUsuario(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		var u model.Usuario
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		// Campos obrigatórios
		if u.Nome == "" || u.Email == "" || u.Senha == "" {
			http.Error(w, "Campos obrigatórios", http.StatusBadRequest)
			return
		}

		if !validarCPF(u.CPF) {
			http.Error(w, "CPF deve conter 11 números.", http.StatusBadRequest)
			return
		}

		if err := repository.CriarUsuario(db, &u); err != nil {
			http.Error(w, "Erro ao salvar", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated) // 201
	}
}

// AtualizarUsuario edita os dados do usuário. Senha é opcional — se vier em branco, não altera.
// Rota: PUT /api/atualizar/usuarios
func AtualizarUsuario(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		var u model.Usuario
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Erro no JSON", http.StatusBadRequest)
			return
		}

		if u.Nome == "" || u.Email == "" || strings.TrimSpace(u.CPF) == "" {
			http.Error(w, "Campos obrigatórios", http.StatusBadRequest)
			return
		}

		if !validarCPF(u.CPF) {
			http.Error(w, "CPF deve conter 11 números.", http.StatusBadRequest)
			return
		}

		if err := repository.AtualizarUsuario(db, &u); err != nil {
			http.Error(w, "Erro ao atualizar", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// BuscarTodosUsuario lista todos os usuários sem filtro.
// Rota: GET /api/todos/usuario
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

// BuscarUsuariosComFiltro pesquisa por nome ou email — recebe ?nome= na query string.
// Rota: GET /api/usuarios/buscar?nome=
func BuscarUsuariosComFiltro(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filtro := r.URL.Query().Get("nome")

		usuarios, err := repository.BuscarUsuariosComFiltro(db, filtro)
		if err != nil {
			http.Error(w, "Erro ao buscar usuários", http.StatusInternalServerError)
			return
		}

		// Garante que retorna [] em vez de null quando não há resultados
		if usuarios == nil {
			usuarios = []model.Usuario{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(usuarios)
	}
}

// ExcluirUsuario remove o usuário pelo ID passado na URL.
// Rota: DELETE /api/usuarios/excluir/{id}
func ExcluirUsuario(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		// Extrai o ID da URL — ex: /api/usuarios/excluir/5 → id = 5
		idStr := strings.TrimPrefix(r.URL.Path, "/api/usuarios/excluir/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		if err := repository.ExcluirUsuario(db, id); err != nil {
			http.Error(w, "Erro ao excluir usuário", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// UsuarioPorID retorna os dados de um usuário pelo ID na URL.
// Rota: GET /api/usuarios/listar/{id}
func UsuarioPorID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/usuarios/listar/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		user, err := repository.BuscarUsuarioPorID(db, id)
		if err != nil {
			http.Error(w, "Usuário não encontrado", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

// Login valida email e senha e retorna os dados do usuário + um token simples.
// O token é apenas "dummy-token-{id}" — não tem segurança real, é só pra identificar a sessão.
// Rota: POST /api/auth/login
func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds struct {
			Email string `json:"email"`
			Senha string `json:"senha"`
		}

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		user, err := repository.AutenticarUsuario(db, creds.Email, creds.Senha)
		if err != nil {
			// Credenciais erradas — retorna 401
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Credenciais inválidas"})
			return
		}

		response := map[string]interface{}{
			"token": "dummy-token-" + strconv.Itoa(user.ID),
			"user":  user,
		}
		json.NewEncoder(w).Encode(response)
	}
}
