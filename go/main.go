package main

import (
	"net/http"
	"sistema-confeitaria/database"
	"sistema-confeitaria/handler"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()

	database.RunMigrations(db)

	// rota principal (POST de ciração)
	http.HandleFunc("/api/novo/usuario", handler.CriarUsuario(db))

	//GET para listar usuarios geral
	http.HandleFunc("/api/todos/usuario", handler.BuscarTodosUsuario(db))

	//GET para listar usuario por id
	http.HandleFunc("/api/usuarios/listar/", handler.UsuarioPorID(db))

	// rota com ID (PUT usuário)
	http.HandleFunc("/api/atualizar/usuarios/{id}", handler.AtualizarUsuario(db))

	// servir frontend
	fs := http.FileServer(http.Dir("./html"))
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil)
}
