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
