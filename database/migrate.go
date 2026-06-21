package database

import (
	"database/sql"
	"log"
)

// RunMigrations cria todas as tabelas do sistema caso ainda não existam.
// É chamada uma vez na inicialização do servidor (em main.go).
// O IF NOT EXISTS garante que rodar de novo não apaga nada que já existe.
func RunMigrations(db *sql.DB) {

	// Lista de tabelas na ordem certa — as que têm FK precisam vir depois das que referenciam
	stmts := []struct {
		name string
		sql  string
	}{
		// Tabela de usuários do sistema (quem faz login)
		{"USUARIO", `
		CREATE TABLE IF NOT EXISTS USUARIO (
			id_usuario   integer PRIMARY KEY,
			nome_usuario varchar(100),
			cpf          char(14),
			email_usuario varchar(100),
			senha        varchar(100)
		);`},

		// Formas de pagamento aceitas (Dinheiro, Pix, Cartão, etc.)
		// ativo: 'S' = ativa, 'N' = inativa — o CHECK impede qualquer outro valor
		{"FORMA_PAGAMENTO", `
		CREATE TABLE IF NOT EXISTS FORMA_PAGAMENTO (
			id_forma_pagamento integer PRIMARY KEY,
			descricao          varchar(100),
			ativo              char(2) CHECK (ativo IN ('S','N'))
		);`},

		// Clientes da confeitaria
		// localizador: campo livre pra endereço complementar, ponto de referência, etc.
		{"CLIENTE", `
		CREATE TABLE IF NOT EXISTS CLIENTE (
			id_cliente   integer PRIMARY KEY,
			nome         varchar(100),
			cpf          char(14),
			telefone     varchar(11),
			email        varchar(100),
			endereco     varchar(150),
			localizador  varchar(500)
		);`},

		// Produtos do catálogo (bolos, tortas, doces...)
		// ativo: mesmo padrão 'S'/'N' da FORMA_PAGAMENTO
		{"PRODUTO", `
		CREATE TABLE IF NOT EXISTS PRODUTO (
			id_produto     integer PRIMARY KEY,
			nome           varchar(100),
			preco_unitario numeric(7,2),
			descricao      varchar(100),
			categoria      varchar(100),
			ativo          char(2) CHECK (ativo IN ('S','N'))
		);`},

		// Cabeçalho do pedido — dados gerais, sem os itens (itens ficam em ITEM_PEDIDO)
		{"PEDIDO", `
		CREATE TABLE IF NOT EXISTS PEDIDO (
			id_pedido          integer PRIMARY KEY,
			status_entrega     varchar(100),
			dt_pedido          datetime,
			vlr_total_pedido   numeric(7,2),
			dt_entrega         datetime,
			observacao_entrega varchar(100),
			id_cliente         integer,
			id_forma_pagamento integer,
			id_usuario         integer,
			FOREIGN KEY(id_cliente)         REFERENCES CLIENTE(id_cliente),
			FOREIGN KEY(id_forma_pagamento) REFERENCES FORMA_PAGAMENTO(id_forma_pagamento),
			FOREIGN KEY(id_usuario)         REFERENCES USUARIO(id_usuario)
		);`},

		// Itens de cada pedido — cada linha é um produto dentro de um pedido
		// vlr_total_item = preco_unitario * qtd_item (calculado no front e salvo aqui)
		{"ITEM_PEDIDO", `
		CREATE TABLE IF NOT EXISTS ITEM_PEDIDO (
			id_item        integer PRIMARY KEY,
			qtd_item       integer,
			vlr_total_item numeric(7,2),
			id_produto     integer,
			id_pedido      integer,
			FOREIGN KEY(id_produto) REFERENCES PRODUTO(id_produto),
			FOREIGN KEY(id_pedido)  REFERENCES PEDIDO(id_pedido)
		);`},
	}

	for _, s := range stmts {
		if _, err := db.Exec(s.sql); err != nil {
			log.Fatalf("Erro ao criar tabela %s: %v", s.name, err)
		}
	}

	log.Println("Tabelas criadas/verificadas com sucesso!")
}
