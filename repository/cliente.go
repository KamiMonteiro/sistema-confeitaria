package repository

import (
	"database/sql"
	"sistema-confeitaria/model"
)

// CriarCliente insere um novo cliente no banco.
// Campos como telefone, email e localizador são opcionais — podem vir vazios.
func CriarCliente(db *sql.DB, c *model.Cliente) error {
	_, err := db.Exec(
		`INSERT INTO CLIENTE (nome, cpf, telefone, email, endereco, localizador) VALUES (?,?,?,?,?,?)`,
		c.Nome, c.CPF, c.Telefone, c.Email, c.Endereco, c.Localizador,
	)
	return err
}

// AtualizarCliente atualiza todos os campos pelo ID.
func AtualizarCliente(db *sql.DB, c *model.Cliente) error {
	_, err := db.Exec(
		`UPDATE CLIENTE SET nome=?, cpf=?, telefone=?, email=?, endereco=?, localizador=? WHERE id_cliente=?`,
		c.Nome, c.CPF, c.Telefone, c.Email, c.Endereco, c.Localizador, c.ID,
	)
	return err
}

// BuscarClientePorID retorna um cliente específico — usado na tela de edição.
func BuscarClientePorID(db *sql.DB, id int) (*model.Cliente, error) {
	row := db.QueryRow(
		`SELECT id_cliente, nome, cpf, telefone, email, endereco, localizador
		 FROM CLIENTE WHERE id_cliente=?`, id)

	var c model.Cliente
	err := row.Scan(&c.ID, &c.Nome, &c.CPF, &c.Telefone, &c.Email, &c.Endereco, &c.Localizador)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// BuscarTodosClientes retorna todos os clientes em ordem alfabética pelo nome.
// Usada na listagem principal e para popular o select de clientes no pedido.
func BuscarTodosClientes(db *sql.DB) ([]model.Cliente, error) {
	rows, err := db.Query(
		`SELECT id_cliente, nome, cpf, telefone, email, endereco, localizador
		 FROM CLIENTE ORDER BY nome ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []model.Cliente
	for rows.Next() {
		var c model.Cliente
		if err := rows.Scan(&c.ID, &c.Nome, &c.CPF, &c.Telefone, &c.Email, &c.Endereco, &c.Localizador); err != nil {
			return nil, err
		}
		lista = append(lista, c)
	}
	return lista, nil
}

// BuscarClientesPorNome filtra por nome OU email — usado na pesquisa da listagem.
// Busca parcial e case-insensitive.
func BuscarClientesPorNome(db *sql.DB, nome string) ([]model.Cliente, error) {
	filtro := "%" + nome + "%"

	rows, err := db.Query(
		`SELECT id_cliente, nome, cpf, telefone, email, endereco, localizador
		 FROM CLIENTE
		 WHERE UPPER(nome) LIKE UPPER(?) OR UPPER(email) LIKE UPPER(?)
		 ORDER BY nome ASC`,
		filtro, filtro)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []model.Cliente
	for rows.Next() {
		var c model.Cliente
		if err := rows.Scan(&c.ID, &c.Nome, &c.CPF, &c.Telefone, &c.Email, &c.Endereco, &c.Localizador); err != nil {
			return nil, err
		}
		lista = append(lista, c)
	}
	return lista, nil
}

// ExcluirCliente remove o cliente pelo ID.
// Atenção: se o cliente tiver pedidos vinculados, o banco vai barrar por causa da FK.
func ExcluirCliente(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM CLIENTE WHERE id_cliente=?`, id)
	return err
}
