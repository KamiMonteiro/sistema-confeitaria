package model

// Usuario representa quem tem acesso ao sistema (tela de login).
// A tag json define como o campo aparece no JSON da API.
type Usuario struct {
	ID    int    `json:"id_usuario"`
	Nome  string `json:"nome_usuario"`
	CPF   string `json:"cpf"`
	Email string `json:"email_usuario"`
	Senha string `json:"senha"`
}

// FormaPagamento são os meios de pagamento aceitos (Pix, Dinheiro, Cartão...).
// Ativo: 'S' = ativa, 'N' = inativa — controlado pela tela de pagamentos.
type FormaPagamento struct {
	ID        int    `json:"id_forma_pagamento"`
	Descricao string `json:"descricao"`
	Ativo     string `json:"ativo"`
}

// Cliente é quem faz pedidos na confeitaria.
// Localizador é um campo livre — ponto de referência, complemento de endereço, etc.
type Cliente struct {
	ID          int    `json:"id_cliente"`
	Nome        string `json:"nome"`
	CPF         string `json:"cpf"`
	Telefone    string `json:"telefone"`
	Email       string `json:"email"`
	Endereco    string `json:"endereco"`
	Localizador string `json:"localizador"`
}

// Produto é qualquer item vendável do catálogo (bolo, torta, doce...).
// Ativo: 'S'/'N' — só produtos ativos aparecem na seleção do pedido.
type Produto struct {
	ID            int     `json:"id_produto"`
	Nome          string  `json:"nome"`
	PrecoUnitario float64 `json:"preco_unitario"`
	Descricao     string  `json:"descricao"`
	Categoria     string  `json:"categoria"`
	Ativo         string  `json:"ativo"`
}

// ItemPedido é cada linha de produto dentro de um pedido.
// NomeProduto e PrecoUnit são preenchidos só na leitura (JOIN com PRODUTO) — omitempty evita
// que apareçam vazios no JSON quando não for necessário.
type ItemPedido struct {
	ID          int     `json:"id_item"`
	QtdItem     int     `json:"qtd_item"`
	VlrTotal    float64 `json:"vlr_total_item"`
	IDProduto   int     `json:"id_produto"`
	IDPedido    int     `json:"id_pedido"`
	NomeProduto string  `json:"nome_produto,omitempty"`
	PrecoUnit   float64 `json:"preco_unitario,omitempty"`
}

// Pedido é o cabeçalho do pedido. Os itens ficam em ItemPedido.
// NomeCliente e NomeFormaPgto só vêm preenchidos nas listagens (JOIN) — não são salvos no banco.
type Pedido struct {
	ID               int          `json:"id_pedido"`
	StatusEntrega    string       `json:"status_entrega"`
	DtPedido         string       `json:"dt_pedido"`
	VlrTotal         float64      `json:"vlr_total_pedido"`
	DtEntrega        string       `json:"dt_entrega"`
	ObsEntrega       string       `json:"observacao_entrega"`
	IDCliente        int          `json:"id_cliente"`
	IDFormaPagamento int          `json:"id_forma_pagamento"`
	IDUsuario        int          `json:"id_usuario"`
	Itens            []ItemPedido `json:"itens,omitempty"`
	NomeCliente      string       `json:"nome_cliente,omitempty"`
	NomeFormaPgto    string       `json:"nome_forma_pgto,omitempty"`
}
