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

	// rota principal (POST e GET lista)
	http.HandleFunc("/api/usuarios", handler.Usuarios(db))

	// rota com ID (GET por ID e PUT)
	http.HandleFunc("/api/usuarios/", handler.UsuarioPorID(db))

	// servir frontend
	fs := http.FileServer(http.Dir("./html"))
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil)
}
