package models

type Client struct {
	ID        int    `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	CPF_CNPJ  string `db:"cpf_cnpj" json:"cpf_cnpj"`
	Blocklist bool   `db:"blocklist" json:"blocklist"`
}
