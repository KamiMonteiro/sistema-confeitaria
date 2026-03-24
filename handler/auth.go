package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"sistema-confeitaria/repository"
)

type LoginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		user, err := repository.BuscarUsuarioPorEmail(db, req.Email)
		if err != nil {
			http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
			return
		}

		if user.Senha != req.Senha {
			http.Error(w, "Senha inválida", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login realizado com sucesso"))
	}
}
