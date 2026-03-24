package repository

import (
	"database/sql"
	"sistema-confeitaria/model"
)

func BuscarUsuarioPorEmail(db *sql.DB, email string) (*model.Usuario, error) {

	query := `SELECT id_usuario, nome_usuario, cpf, email_usuario, senha FROM USUARIO WHERE email_usuario = ?`

	row := db.QueryRow(query, email)

	var u model.Usuario
	err := row.Scan(
		&u.ID,
		&u.Nome,
		&u.CPF,
		&u.Email,
		&u.Senha,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
