package main

import (
	"log"
	"sistema-confeitaria/database"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()

	log.Println("Aplicaçção está rodando...)")

	database.RunMigrations(db)
}
