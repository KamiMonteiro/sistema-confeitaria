package model

type Usuario struct {
	ID    int    `json:"id_usuario"`
	Nome  string `json:"nome_usuario"`
	CPF   string `json:"cpf"`
	Email string `json:"email_usuario"`
	Senha string `json:"senha"`
}
