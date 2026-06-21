package repository

import (
	"database/sql"
	"sistema-confeitaria/model"
)

// CriarUsuario insere um novo usuário no banco.
// A senha é salva como texto simples — sem hash por enquanto.
func CriarUsuario(db *sql.DB, u *model.Usuario) error {
	query := `INSERT INTO USUARIO (nome_usuario, cpf, email_usuario, senha)
	VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, u.Nome, u.CPF, u.Email, u.Senha)
	return err
}

// AtualizarUsuario atualiza os dados do usuário.
// Se a senha vier em branco no payload, não atualiza ela — mantém a que já está no banco.
func AtualizarUsuario(db *sql.DB, u *model.Usuario) error {
	var query string

	if u.Senha != "" {
		// Senha foi informada — atualiza junto
		query = `UPDATE USUARIO SET nome_usuario=?, email_usuario=?, cpf=?, senha=? WHERE id_usuario=?`
		_, err := db.Exec(query, u.Nome, u.Email, u.CPF, u.Senha, u.ID)
		return err
	}

	// Senha em branco — só atualiza nome, email e cpf
	query = `UPDATE USUARIO SET nome_usuario=?, email_usuario=?, cpf=? WHERE id_usuario=?`
	_, err := db.Exec(query, u.Nome, u.Email, u.CPF, u.ID)
	return err
}

// ExcluirUsuario remove o usuário pelo ID.
func ExcluirUsuario(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM USUARIO WHERE id_usuario = ?`, id)
	return err
}

// BuscarUsuarioPorID retorna um usuário específico pelo ID.
// A senha não é retornada por segurança — só os dados de exibição.
func BuscarUsuarioPorID(db *sql.DB, id int) (*model.Usuario, error) {
	row := db.QueryRow(`
		SELECT id_usuario, nome_usuario, cpf, email_usuario
		FROM USUARIO
		WHERE id_usuario = ?`, id)

	var u model.Usuario
	err := row.Scan(&u.ID, &u.Nome, &u.CPF, &u.Email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// BuscarTodosUsuario retorna todos os usuários em ordem alfabética.
// Usada na tela de listagem — sem filtro.
func BuscarTodosUsuario(db *sql.DB) ([]model.Usuario, error) {
	rows, err := db.Query(`
		SELECT id_usuario, nome_usuario, cpf, email_usuario
		FROM USUARIO
		ORDER BY nome_usuario COLLATE NOCASE ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []model.Usuario
	for rows.Next() {
		var u model.Usuario
		if err := rows.Scan(&u.ID, &u.Nome, &u.CPF, &u.Email); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}
	return usuarios, nil
}

// BuscarUsuariosComFiltro busca pelo nome OU email — usado na caixa de pesquisa da lista.
// O % nos dois lados do filtro faz a busca conter (LIKE %termo%).
func BuscarUsuariosComFiltro(db *sql.DB, filtro string) ([]model.Usuario, error) {
	filtroLike := "%" + filtro + "%"

	rows, err := db.Query(`
		SELECT id_usuario, nome_usuario, cpf, email_usuario
		FROM USUARIO
		WHERE nome_usuario LIKE ? OR email_usuario LIKE ?
		ORDER BY nome_usuario COLLATE NOCASE ASC`,
		filtroLike, filtroLike)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []model.Usuario
	for rows.Next() {
		var u model.Usuario
		if err := rows.Scan(&u.ID, &u.Nome, &u.CPF, &u.Email); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}
	return usuarios, nil
}

// AutenticarUsuario verifica email e senha — usada no login.
// Compara direto no banco (texto simples). Se não encontrar, retorna erro.
func AutenticarUsuario(db *sql.DB, email, senha string) (*model.Usuario, error) {
	row := db.QueryRow(`
		SELECT id_usuario, nome_usuario, cpf, email_usuario
		FROM USUARIO
		WHERE email_usuario = ? AND senha = ?`,
		email, senha)

	var u model.Usuario
	err := row.Scan(&u.ID, &u.Nome, &u.CPF, &u.Email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
