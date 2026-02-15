package secrets

// -- Read-side --

const (
	ListUsernamesSQL = `SELECT DISTINCT username FROM secrets WHERE pubkey = ?;`
	ListServicesSQL  = `SELECT service FROM secrets WHERE pubkey = ? AND username = ?;`
	GetSecretSQL     = `SELECT cipher, nonce FROM secrets WHERE pubkey = ? AND username = ? AND service = ?;`
)

// -- Write-side --

const UpsertSecretSQL = `
	INSERT INTO secrets (pubkey, username, service, cipher, nonce)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(pubkey, username, service)
	DO UPDATE SET
		cipher = excluded.cipher,
		nonce = excluded.nonce;
`
