package database

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) {
	createTable := `
	CREATE TABLE IF NOT EXISTS USUARIO (
	id_usuario integer PRIMARY KEY,
	nome_usuario varchar(100),
	cpf char(14),
	email_usuario varchar(100),
	senha varchar(100)
);
	`

	_, err := db.Exec(createTable)
	if err != nil {
		log.Fatal("Erro ao criar tabela:", err)
	}

	log.Println("Tabela criada/verificada com sucesso!")
}
