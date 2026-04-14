package repository

import (
	"database/sql"
	"sistema-confeitaria/model"
)

func CriarFormaPagamento(db *sql.DB, f *model.FormaPagamento) error {
	query := `INSERT INTO FORMA_PAGAMENTO (descricao, ativo)
VALUES (?, ?)`

	_, err := db.Exec(query,
		f.Descricao,
		f.Ativo,
	)

	return err
}

func AtualizarFormaPagamento(db *sql.DB, f *model.FormaPagamento) error {
	query := `
UPDATE FORMA_PAGAMENTO
SET descricao = ?, ativo = ?
WHERE id_forma_pagamento = ?
`

	_, err := db.Exec(query,
		f.Descricao,
		f.Ativo,
		f.ID,
	)

	return err
}

func BuscarFormaPagamentoPorID(db *sql.DB, id int) (*model.FormaPagamento, error) {
	query := `
SELECT id_forma_pagamento, descricao, ativo
FROM FORMA_PAGAMENTO
WHERE id_forma_pagamento = ?
`

	row := db.QueryRow(query, id)

	var f model.FormaPagamento
	err := row.Scan(&f.ID, &f.Descricao, &f.Ativo)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func BuscarTodasFormasPagamento(db *sql.DB) ([]model.FormaPagamento, error) {
	query := `
SELECT id_forma_pagamento, descricao, ativo
FROM FORMA_PAGAMENTO
ORDER BY descricao ASC
`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var formas []model.FormaPagamento
	for rows.Next() {
		var f model.FormaPagamento
		err := rows.Scan(&f.ID, &f.Descricao, &f.Ativo)
		if err != nil {
			return nil, err
		}
		formas = append(formas, f)
	}

	return formas, nil
}

func ExcluirFormaPagamento(db *sql.DB, id int) error {
	query := `
DELETE FROM FORMA_PAGAMENTO
WHERE id_forma_pagamento = ?
`

	_, err := db.Exec(query, id)
	return err
}
