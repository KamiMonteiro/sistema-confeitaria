package database

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) {
	createTableUsuario := `
	CREATE TABLE IF NOT EXISTS USUARIO (
		id_usuario integer PRIMARY KEY,
		nome_usuario varchar(100),
		cpf char(14),
		email_usuario varchar(100),
		senha varchar(100)
	);
	`

	createTableFormaPagamento := `
	CREATE TABLE IF NOT EXISTS FORMA_PAGAMENTO (
		id_forma_pagamento integer PRIMARY KEY,
		descricao varchar(100),
		ativo char(2) CHECK (ativo IN ('S','N'))
	);
	`

	_, err := db.Exec(createTableUsuario)
	if err != nil {
		log.Fatal("Erro ao criar tabela USUARIO:", err)
	}

	_, err = db.Exec(createTableFormaPagamento)
	if err != nil {
		log.Fatal("Erro ao criar tabela FORMA_PAGAMENTO:", err)
	}

	log.Println("Tabelas criadas/verificadas com sucesso!")
}
