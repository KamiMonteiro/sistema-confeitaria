package main

import (
	"net/http"
	"sistema-confeitaria/database"
	"sistema-confeitaria/handler"
)

func corsHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	db := database.ConnectDB()
	defer db.Close()

	database.RunMigrations(db)

	// rota principal (POST de ciração)
	http.HandleFunc("/api/novo/usuario", corsHandler(handler.CriarUsuario(db)))

	//GET para listar usuarios geral
	http.HandleFunc("/api/todos/usuario", corsHandler(handler.BuscarTodosUsuario(db)))

	//GET para listar usuario por id
	http.HandleFunc("/api/usuarios/listar/", corsHandler(handler.UsuarioPorID(db)))

	// DELETE para excluir usuário por id
	http.HandleFunc("/api/usuarios/excluir/", corsHandler(handler.ExcluirUsuario(db)))

	// rota com ID (PUT usuário)
	http.HandleFunc("/api/atualizar/usuarios", corsHandler(handler.AtualizarUsuario(db)))

	// rota de login
	http.HandleFunc("/api/auth/login", corsHandler(handler.Login(db)))

	// servir frontend
	fs := http.FileServer(http.Dir("./html"))
	http.Handle("/", corsMiddleware(fs))

	http.ListenAndServe(":8080", nil)
}
