CREATE TABLE secrets(
    pubkey TEXT NOT NULL,
    username TEXT NOT NULL,
    service TEXT NOT NULL,
    cipher TEXT NOT NULL,
    nonce TEXT NOT NULL,
    PRIMARY KEY(pubkey, username, service)
);