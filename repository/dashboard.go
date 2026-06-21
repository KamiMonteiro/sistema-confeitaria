package repository

import (
	"database/sql"
	"sistema-confeitaria/model"
)

// DashboardStats agrupa todos os números do dashboard num único retorno.
type DashboardStats struct {
	PedidosHoje      int            `json:"pedidos_hoje"`
	FaturamentoHoje  float64        `json:"faturamento_hoje"`
	ProdutosAtivos   int            `json:"produtos_ativos"`
	TotalClientes    int            `json:"total_clientes"`
	UltimosPedidos   []model.Pedido `json:"ultimos_pedidos"`
}

// BuscarEstatisticasDashboard busca todos os KPIs em consultas separadas e retorna tudo junto.
// "Hoje" é comparado com date('now') do SQLite — assume que dt_pedido vem no formato YYYY-MM-DD.
func BuscarEstatisticasDashboard(db *sql.DB) (*DashboardStats, error) {
	stats := &DashboardStats{}

	// Pedidos registrados hoje
	err := db.QueryRow(`SELECT COUNT(*) FROM PEDIDO WHERE date(dt_pedido) = date('now')`).
		Scan(&stats.PedidosHoje)
	if err != nil {
		return nil, err
	}

	// Soma do faturamento dos pedidos de hoje — COALESCE pra não retornar NULL quando não tem nenhum
	err = db.QueryRow(`SELECT COALESCE(SUM(vlr_total_pedido), 0) FROM PEDIDO WHERE date(dt_pedido) = date('now')`).
		Scan(&stats.FaturamentoHoje)
	if err != nil {
		return nil, err
	}

	// Produtos ativos no catálogo
	err = db.QueryRow(`SELECT COUNT(*) FROM PRODUTO WHERE ativo = 'S'`).
		Scan(&stats.ProdutosAtivos)
	if err != nil {
		return nil, err
	}

	// Total de clientes cadastrados
	err = db.QueryRow(`SELECT COUNT(*) FROM CLIENTE`).
		Scan(&stats.TotalClientes)
	if err != nil {
		return nil, err
	}

	// Últimos 5 pedidos — com nome do cliente e forma de pagamento pelo JOIN
	rows, err := db.Query(`
		SELECT p.id_pedido, p.status_entrega, p.dt_pedido, p.vlr_total_pedido,
		       p.dt_entrega, p.observacao_entrega, p.id_cliente, p.id_forma_pagamento, p.id_usuario,
		       COALESCE(c.nome,''), COALESCE(f.descricao,'')
		FROM PEDIDO p
		LEFT JOIN CLIENTE c ON c.id_cliente = p.id_cliente
		LEFT JOIN FORMA_PAGAMENTO f ON f.id_forma_pagamento = p.id_forma_pagamento
		ORDER BY p.id_pedido DESC
		LIMIT 5`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Pedido
		if err := rows.Scan(&p.ID, &p.StatusEntrega, &p.DtPedido, &p.VlrTotal,
			&p.DtEntrega, &p.ObsEntrega, &p.IDCliente, &p.IDFormaPagamento, &p.IDUsuario,
			&p.NomeCliente, &p.NomeFormaPgto); err != nil {
			return nil, err
		}
		stats.UltimosPedidos = append(stats.UltimosPedidos, p)
	}

	// Garante [] em vez de null quando não há pedidos
	if stats.UltimosPedidos == nil {
		stats.UltimosPedidos = []model.Pedido{}
	}

	return stats, nil
}
