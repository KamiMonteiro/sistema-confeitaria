package repository

import (
	"database/sql"
	"sistema-confeitaria/model"
)

// CriarFormaPagamento insere uma nova forma de pagamento.
// O ativo já vem normalizado como 'S' ou 'N' pelo handler antes de chegar aqui.
func CriarFormaPagamento(db *sql.DB, f *model.FormaPagamento) error {
	_, err := db.Exec(
		`INSERT INTO FORMA_PAGAMENTO (descricao, ativo) VALUES (?, ?)`,
		f.Descricao, f.Ativo,
	)
	return err
}

// AtualizarFormaPagamento atualiza descrição e status pelo ID.
func AtualizarFormaPagamento(db *sql.DB, f *model.FormaPagamento) error {
	_, err := db.Exec(
		`UPDATE FORMA_PAGAMENTO SET descricao = ?, ativo = ? WHERE id_forma_pagamento = ?`,
		f.Descricao, f.Ativo, f.ID,
	)
	return err
}

// BuscarFormaPagamentoPorID retorna uma forma de pagamento pelo ID.
// Usada na tela de edição pra carregar os dados no formulário.
func BuscarFormaPagamentoPorID(db *sql.DB, id int) (*model.FormaPagamento, error) {
	row := db.QueryRow(`
		SELECT id_forma_pagamento, descricao, ativo
		FROM FORMA_PAGAMENTO
		WHERE id_forma_pagamento = ?`, id)

	var f model.FormaPagamento
	err := row.Scan(&f.ID, &f.Descricao, &f.Ativo)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// BuscarTodasFormasPagamento retorna todas as formas, ativas e inativas, em ordem alfabética.
// Usada na listagem e também pra popular o select do formulário de pedido.
func BuscarTodasFormasPagamento(db *sql.DB) ([]model.FormaPagamento, error) {
	rows, err := db.Query(`
		SELECT id_forma_pagamento, descricao, ativo
		FROM FORMA_PAGAMENTO
		ORDER BY descricao ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var formas []model.FormaPagamento
	for rows.Next() {
		var f model.FormaPagamento
		if err := rows.Scan(&f.ID, &f.Descricao, &f.Ativo); err != nil {
			return nil, err
		}
		formas = append(formas, f)
	}
	return formas, nil
}

// BuscarFormasPagamentoPorDescricao filtra pelo nome — busca parcial, ignora maiúsculas/minúsculas.
// Usada na caixa de pesquisa da tela de formas de pagamento.
func BuscarFormasPagamentoPorDescricao(db *sql.DB, descricao string) ([]model.FormaPagamento, error) {
	filtroLike := "%" + descricao + "%"

	rows, err := db.Query(`
		SELECT id_forma_pagamento, descricao, ativo
		FROM FORMA_PAGAMENTO
		WHERE UPPER(descricao) LIKE UPPER(?)
		ORDER BY descricao COLLATE NOCASE ASC`,
		filtroLike)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var formas []model.FormaPagamento
	for rows.Next() {
		var f model.FormaPagamento
		if err := rows.Scan(&f.ID, &f.Descricao, &f.Ativo); err != nil {
			return nil, err
		}
		formas = append(formas, f)
	}
	return formas, nil
}

// ExcluirFormaPagamento remove pelo ID.
func ExcluirFormaPagamento(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM FORMA_PAGAMENTO WHERE id_forma_pagamento = ?`, id)
	return err
}
