package main

import (
	"log"
	"net/http"
	"sistema-confeitaria/database"
	"sistema-confeitaria/handler"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()

	log.Println("Aplicaçção está rodando...)")

	database.RunMigrations(db)

	http.HandleFunc("/api/auth/login", handler.Login(db))

	http.ListenAndServe(":8080", nil)
}
