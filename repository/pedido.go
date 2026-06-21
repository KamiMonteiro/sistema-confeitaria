package repository

import (
	"database/sql"
	"sistema-confeitaria/model"
)

// CriarPedido salva o pedido e todos os seus itens numa única transação.
// Se qualquer item falhar, o pedido inteiro é desfeito (Rollback).
// Retorna o ID gerado pro pedido.
func CriarPedido(db *sql.DB, p *model.Pedido) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	// 1. Insere o cabeçalho do pedido
	res, err := tx.Exec(
		`INSERT INTO PEDIDO (status_entrega, dt_pedido, vlr_total_pedido, dt_entrega, observacao_entrega,
		  id_cliente, id_forma_pagamento, id_usuario)
		 VALUES (?,?,?,?,?,?,?,?)`,
		p.StatusEntrega, p.DtPedido, p.VlrTotal, p.DtEntrega, p.ObsEntrega,
		p.IDCliente, p.IDFormaPagamento, p.IDUsuario,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Pega o ID gerado automaticamente pelo SQLite
	idPedido, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// 2. Insere cada item vinculado ao pedido recém-criado
	for _, item := range p.Itens {
		_, err = tx.Exec(
			`INSERT INTO ITEM_PEDIDO (qtd_item, vlr_total_item, id_produto, id_pedido) VALUES (?,?,?,?)`,
			item.QtdItem, item.VlrTotal, item.IDProduto, idPedido,
		)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return idPedido, tx.Commit()
}

// AtualizarPedido edita o cabeçalho do pedido e substitui todos os itens.
// Estratégia: deleta os itens antigos e insere os novos — mais simples que fazer diff.
func AtualizarPedido(db *sql.DB, p *model.Pedido) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 1. Atualiza o cabeçalho (dt_pedido não é alterada na edição — mantém a data original)
	_, err = tx.Exec(
		`UPDATE PEDIDO SET status_entrega=?, dt_entrega=?, vlr_total_pedido=?, observacao_entrega=?,
		  id_cliente=?, id_forma_pagamento=?, id_usuario=?
		 WHERE id_pedido=?`,
		p.StatusEntrega, p.DtEntrega, p.VlrTotal, p.ObsEntrega,
		p.IDCliente, p.IDFormaPagamento, p.IDUsuario, p.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2. Remove todos os itens antigos desse pedido
	if _, err = tx.Exec(`DELETE FROM ITEM_PEDIDO WHERE id_pedido=?`, p.ID); err != nil {
		tx.Rollback()
		return err
	}

	// 3. Insere os itens novos (vindos do form)
	for _, item := range p.Itens {
		_, err = tx.Exec(
			`INSERT INTO ITEM_PEDIDO (qtd_item, vlr_total_item, id_produto, id_pedido) VALUES (?,?,?,?)`,
			item.QtdItem, item.VlrTotal, item.IDProduto, p.ID,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// BuscarPedidoPorID retorna o pedido com todos os seus itens.
// Faz JOIN com CLIENTE e FORMA_PAGAMENTO pra trazer os nomes direto (evita chamadas extras).
func BuscarPedidoPorID(db *sql.DB, id int) (*model.Pedido, error) {
	row := db.QueryRow(`
		SELECT p.id_pedido, p.status_entrega, p.dt_pedido, p.vlr_total_pedido,
		       p.dt_entrega, p.observacao_entrega, p.id_cliente, p.id_forma_pagamento, p.id_usuario,
		       COALESCE(c.nome,''), COALESCE(f.descricao,'')
		FROM PEDIDO p
		LEFT JOIN CLIENTE c ON c.id_cliente = p.id_cliente
		LEFT JOIN FORMA_PAGAMENTO f ON f.id_forma_pagamento = p.id_forma_pagamento
		WHERE p.id_pedido = ?`, id)

	var p model.Pedido
	err := row.Scan(&p.ID, &p.StatusEntrega, &p.DtPedido, &p.VlrTotal,
		&p.DtEntrega, &p.ObsEntrega, &p.IDCliente, &p.IDFormaPagamento, &p.IDUsuario,
		&p.NomeCliente, &p.NomeFormaPgto)
	if err != nil {
		return nil, err
	}

	// Busca os itens separado — query auxiliar abaixo
	itens, err := buscarItensDoPedido(db, id)
	if err != nil {
		return nil, err
	}
	p.Itens = itens
	return &p, nil
}

// buscarItensDoPedido é uma função interna (minúscula = não exportada).
// Traz os itens com nome e preço do produto pelo JOIN com PRODUTO.
func buscarItensDoPedido(db *sql.DB, idPedido int) ([]model.ItemPedido, error) {
	rows, err := db.Query(`
		SELECT i.id_item, i.qtd_item, i.vlr_total_item, i.id_produto, i.id_pedido,
		       COALESCE(p.nome,''), COALESCE(p.preco_unitario,0)
		FROM ITEM_PEDIDO i
		LEFT JOIN PRODUTO p ON p.id_produto = i.id_produto
		WHERE i.id_pedido = ?`, idPedido)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itens []model.ItemPedido
	for rows.Next() {
		var item model.ItemPedido
		if err := rows.Scan(&item.ID, &item.QtdItem, &item.VlrTotal, &item.IDProduto, &item.IDPedido,
			&item.NomeProduto, &item.PrecoUnit); err != nil {
			return nil, err
		}
		itens = append(itens, item)
	}
	return itens, nil
}

// BuscarTodosPedidos lista todos os pedidos do mais recente pro mais antigo.
// Traz nome do cliente e forma de pagamento pelo JOIN — evita N+1 queries no front.
func BuscarTodosPedidos(db *sql.DB) ([]model.Pedido, error) {
	rows, err := db.Query(`
		SELECT p.id_pedido, p.status_entrega, p.dt_pedido, p.vlr_total_pedido,
		       p.dt_entrega, p.observacao_entrega, p.id_cliente, p.id_forma_pagamento, p.id_usuario,
		       COALESCE(c.nome,''), COALESCE(f.descricao,'')
		FROM PEDIDO p
		LEFT JOIN CLIENTE c ON c.id_cliente = p.id_cliente
		LEFT JOIN FORMA_PAGAMENTO f ON f.id_forma_pagamento = p.id_forma_pagamento
		ORDER BY p.dt_pedido DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []model.Pedido
	for rows.Next() {
		var p model.Pedido
		if err := rows.Scan(&p.ID, &p.StatusEntrega, &p.DtPedido, &p.VlrTotal,
			&p.DtEntrega, &p.ObsEntrega, &p.IDCliente, &p.IDFormaPagamento, &p.IDUsuario,
			&p.NomeCliente, &p.NomeFormaPgto); err != nil {
			return nil, err
		}
		lista = append(lista, p)
	}
	return lista, nil
}

// BuscarPedidosPorStatus filtra pedidos pelo status de entrega (ex: "Pendente", "Entregue").
// Busca parcial — então "pend" já acha "Pendente".
func BuscarPedidosPorStatus(db *sql.DB, status string) ([]model.Pedido, error) {
	rows, err := db.Query(`
		SELECT p.id_pedido, p.status_entrega, p.dt_pedido, p.vlr_total_pedido,
		       p.dt_entrega, p.observacao_entrega, p.id_cliente, p.id_forma_pagamento, p.id_usuario,
		       COALESCE(c.nome,''), COALESCE(f.descricao,'')
		FROM PEDIDO p
		LEFT JOIN CLIENTE c ON c.id_cliente = p.id_cliente
		LEFT JOIN FORMA_PAGAMENTO f ON f.id_forma_pagamento = p.id_forma_pagamento
		WHERE UPPER(p.status_entrega) LIKE UPPER(?)
		ORDER BY p.dt_pedido DESC`,
		"%"+status+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []model.Pedido
	for rows.Next() {
		var p model.Pedido
		if err := rows.Scan(&p.ID, &p.StatusEntrega, &p.DtPedido, &p.VlrTotal,
			&p.DtEntrega, &p.ObsEntrega, &p.IDCliente, &p.IDFormaPagamento, &p.IDUsuario,
			&p.NomeCliente, &p.NomeFormaPgto); err != nil {
			return nil, err
		}
		lista = append(lista, p)
	}
	return lista, nil
}

// ExcluirPedido remove o pedido e todos os seus itens numa transação.
// Precisa deletar os itens primeiro por causa da FK — senão o banco rejeita.
func ExcluirPedido(db *sql.DB, id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Primeiro os itens, depois o pedido
	if _, err = tx.Exec(`DELETE FROM ITEM_PEDIDO WHERE id_pedido=?`, id); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec(`DELETE FROM PEDIDO WHERE id_pedido=?`, id); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
