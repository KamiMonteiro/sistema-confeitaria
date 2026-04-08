package repository

import (
	"database/sql"
	"sistema-confeitaria/model"
)

func CriarUsuario(db *sql.DB, u *model.Usuario) error {

	query := `INSERT INTO USUARIO (nome_usuario, cpf, email_usuario, senha)
	VALUES (?, ?, ?, ?)
	`

	_, err := db.Exec(query,
		u.Nome,
		u.CPF,
		u.Email,
		u.Senha,
	)

	return err
}

func AtualizarUsuario(db *sql.DB, u *model.Usuario) error {

	query := `
	UPDATE USUARIO
	SET nome_usuario = ?, email_usuario = ?, cpf = ?
	WHERE id_usuario = ?
	`

	if u.Senha != "" {
		query = `
		UPDATE USUARIO
		SET nome_usuario=?, email_usuario=?, cpf=?, senha=?
		WHERE id_usuario=?
		`
	} else {
		query = `
		UPDATE USUARIO
		SET nome_usuario=?, email_usuario=?, cpf=?
		WHERE id_usuario=?
		`
	}

	_, err := db.Exec(query,
		u.Nome,
		u.Email,
		u.CPF,
		u.ID,
	)

	return err
}

func ExcluirUsuario(db *sql.DB, id int) error {
	query := `
	DELETE FROM USUARIO
	WHERE id_usuario = ?
	`

	_, err := db.Exec(query, id)
	return err
}

func BuscarUsuarioPorID(db *sql.DB, id int) (*model.Usuario, error) {
	query := `
	SELECT id_usuario, nome_usuario, cpf, email_usuario
	FROM USUARIO
	WHERE id_usuario = ?
	`

	row := db.QueryRow(query, id)

	var u model.Usuario
	err := row.Scan(&u.ID, &u.Nome, &u.CPF, &u.Email)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func BuscarTodosUsuario(db *sql.DB) ([]model.Usuario, error) {

	query := `
	SELECT id_usuario, nome_usuario, cpf, email_usuario
	FROM USUARIO
	ORDER BY nome_usuario ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []model.Usuario

	for rows.Next() {
		var u model.Usuario

		err := rows.Scan(&u.ID, &u.Nome, &u.CPF, &u.Email)
		if err != nil {
			return nil, err
		}

		usuarios = append(usuarios, u)
	}

	return usuarios, nil
}

func AutenticarUsuario(db *sql.DB, email, senha string) (*model.Usuario, error) {
	query := `
	SELECT id_usuario, nome_usuario, cpf, email_usuario
	FROM USUARIO
	WHERE email_usuario = ? AND senha = ?
	`

	row := db.QueryRow(query, email, senha)

	var u model.Usuario
	err := row.Scan(&u.ID, &u.Nome, &u.CPF, &u.Email)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
