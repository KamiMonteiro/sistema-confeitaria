package model

type Usuario struct {
	ID    int    `json:"id_usuario"`
	Nome  string `json:"nome_usuario"`
	CPF   string `json:"cpf"`
	Email string `json:"email_usuario"`
	Senha string `json:"senha"`
}

type FormaPagamento struct {
	ID        int    `json:"id_forma_pagamento"`
	Descricao string `json:"descricao"`
	Ativo     string `json:"ativo"`
}
