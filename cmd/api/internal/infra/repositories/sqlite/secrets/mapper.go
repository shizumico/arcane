package secrets

import (
	"github.com/shizumico/arcane/cmd/api/internal/core/query/secrets"
)

func toSecretView(username, service, cipher, nonce string) secrets.SecretView {
	return secrets.SecretView{
		Service:  service,
		Username: username,
		Cipher:   cipher,
		Nonce:    nonce,
	}
}
