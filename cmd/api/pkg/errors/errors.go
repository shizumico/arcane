package errors

import "errors"

var (
	ErrInvalidPubkeyFormat    = errors.New("invalid pubkey format")
	ErrInvalidSignatureFormat = errors.New("invalid signature format")
	ErrInvalidSignature       = errors.New("invalid signature")
)
