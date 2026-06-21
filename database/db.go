package database

import (
	"database/sql"
	"log"

	// driver do SQLite — precisa do import com _ pra registrar sem usar diretamente
	_ "github.com/mattn/go-sqlite3"
)

// ConnectDB abre a conexão com o banco SQLite e verifica se está acessível.
// O arquivo do banco fica em ./confeitaria.db (relativo a onde o servidor roda, ou seja, dentro de /go)
func ConnectDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./confeitaria.db")
	if err != nil {
		log.Fatal("Erro ao abrir o banco de dados:", err)
	}

	// Ping confirma que a conexão está de pé de verdade
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	log.Println("Banco de dados conectado com sucesso ✅")

	return db
}
