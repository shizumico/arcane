package secrets

type SecretView struct {
	Username string `json:"username"`
	Service  string `json:"service"`
	Cipher   string `json:"cipher"`
	Nonce    string `json:"nonce"`
}

func NewSecretView(username, service, cipher, nonce string) SecretView {
	return SecretView{
		Username: username,
		Service:  service,
		Cipher:   cipher,
		Nonce:    nonce,
	}
}
