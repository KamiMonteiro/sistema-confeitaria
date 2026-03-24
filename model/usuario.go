package model

type Usuario struct {
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	CPF   string `json:"cpf"`
	Email string `json:"email"`
	Senha string `json:"senha"`
}
