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

package handler

import (
	"database/sql"
	"net/http"
)

func Usuarios(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPost:
			// 👉 criar usuário
			CriarUsuario(db, w, r)

		case http.MethodGet:
			// 👉 listar usuários (opcional por enquanto)
			w.Write([]byte("Listagem de usuários"))

		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	}
}


func CriarUsuario(db *sql.DB, w http.ResponseWriter, r *http.Request){
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

func BuscarUsuario(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idStr := strings.TrimPrefix(r.URL.Path, "/api/usuarios/")
		id, _ := strconv.Atoi(idStr)

		user, err := repository.BuscarUsuarioPorID(db, id)
		if err != nil {
			http.Error(w, "Usuário não encontrado", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func UsuarioPorID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			// buscar usuário por ID
			BuscarUsuario(db, w, r)

		case http.MethodPut:
			// atualizar usuário
			AtualizarUsuario(db, w, r)

		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	}
}