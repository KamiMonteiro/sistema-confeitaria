package repository

import (
	"database/sql"
	"sistema-confeitaria/model"
)

// CriarProduto insere um novo produto no catálogo.
// O ativo vem como 'S' ou 'N' — normalizado pelo handler antes.
func CriarProduto(db *sql.DB, p *model.Produto) error {
	_, err := db.Exec(
		`INSERT INTO PRODUTO (nome, preco_unitario, descricao, categoria, ativo) VALUES (?,?,?,?,?)`,
		p.Nome, p.PrecoUnitario, p.Descricao, p.Categoria, p.Ativo,
	)
	return err
}

// AtualizarProduto atualiza todos os campos pelo ID.
func AtualizarProduto(db *sql.DB, p *model.Produto) error {
	_, err := db.Exec(
		`UPDATE PRODUTO SET nome=?, preco_unitario=?, descricao=?, categoria=?, ativo=? WHERE id_produto=?`,
		p.Nome, p.PrecoUnitario, p.Descricao, p.Categoria, p.Ativo, p.ID,
	)
	return err
}

// BuscarProdutoPorID retorna um produto específico — usado na tela de edição.
func BuscarProdutoPorID(db *sql.DB, id int) (*model.Produto, error) {
	row := db.QueryRow(
		`SELECT id_produto, nome, preco_unitario, descricao, categoria, ativo
		 FROM PRODUTO WHERE id_produto=?`, id)

	var p model.Produto
	err := row.Scan(&p.ID, &p.Nome, &p.PrecoUnitario, &p.Descricao, &p.Categoria, &p.Ativo)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// BuscarTodosProdutos retorna todos os produtos em ordem alfabética.
// Inclui ativos e inativos — a filtragem de status é feita no front.
func BuscarTodosProdutos(db *sql.DB) ([]model.Produto, error) {
	rows, err := db.Query(
		`SELECT id_produto, nome, preco_unitario, descricao, categoria, ativo
		 FROM PRODUTO ORDER BY nome ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []model.Produto
	for rows.Next() {
		var p model.Produto
		if err := rows.Scan(&p.ID, &p.Nome, &p.PrecoUnitario, &p.Descricao, &p.Categoria, &p.Ativo); err != nil {
			return nil, err
		}
		lista = append(lista, p)
	}
	return lista, nil
}

// BuscarProdutosPorNome filtra por nome OU categoria — usado na pesquisa da listagem.
// Busca parcial e case-insensitive.
func BuscarProdutosPorNome(db *sql.DB, nome string) ([]model.Produto, error) {
	filtro := "%" + nome + "%"

	rows, err := db.Query(
		`SELECT id_produto, nome, preco_unitario, descricao, categoria, ativo
		 FROM PRODUTO
		 WHERE UPPER(nome) LIKE UPPER(?) OR UPPER(categoria) LIKE UPPER(?)
		 ORDER BY nome ASC`,
		filtro, filtro)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []model.Produto
	for rows.Next() {
		var p model.Produto
		if err := rows.Scan(&p.ID, &p.Nome, &p.PrecoUnitario, &p.Descricao, &p.Categoria, &p.Ativo); err != nil {
			return nil, err
		}
		lista = append(lista, p)
	}
	return lista, nil
}

// ExcluirProduto remove o produto pelo ID.
// Atenção: se o produto tiver itens em pedidos, o banco vai barrar por causa da FK.
func ExcluirProduto(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM PRODUTO WHERE id_produto=?`, id)
	return err
}
