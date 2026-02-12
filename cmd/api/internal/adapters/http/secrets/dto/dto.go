package dto

import "github.com/shizumico/arcane/cmd/api/internal/core/application/secrets"

type SaveSecretRequest struct {
	Username string `json:"username"`
	Service  string `json:"service"`
	Cipher   string `json:"cipher"`
	Nonce    string `json:"nonce"`
}

func (req SaveSecretRequest) ToCommand(pubkey, signature string) secrets.SaveCommand {
	return secrets.SaveCommand{
		Pubkey:    pubkey,
		Username:  req.Username,
		Service:   req.Service,
		Cipher:    req.Cipher,
		Nonce:     req.Nonce,
		Signature: signature,
	}
}

type ListUsernamesResponse struct {
	Usernames []string `json:"usernames,omitempty"`
}

type ListServicesResponse struct {
	Services []string `json:"services,omitempty"`
}
