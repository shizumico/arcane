package secrets

type Secret struct {
	PubKey   string
	Username string
	Service  string
	Cipher   string
	Nonce    string
}

func NewSecret(pubkey, username, service, cipher, nonce string) *Secret {
	return &Secret{
		PubKey:   pubkey,
		Username: username,
		Service:  service,
		Cipher:   cipher,
		Nonce:    nonce,
	}
}
