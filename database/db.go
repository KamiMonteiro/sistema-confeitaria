package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./confeitaria.db")
	if err != nil {
		log.Fatal("Erro ao abrir o banco de dados:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	log.Println("Banco de dados conectado com sucesso ✅")

	return db
}
