-- +goose Up
-- +goose StatementBegin

CREATE TABLE accounts (
    id VARCHAR(50) PRIMARY KEY,
    balance DECIMAL(15, 2) NOT NULL CHECK (balance >= 0),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ledger_entries (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(50) NOT NULL,
    account_id VARCHAR(50) REFERENCES accounts(id),
    amount DECIMAL(15, 2) NOT NULL,
    dc_indicator VARCHAR(1) CHECK (dc_indicator IN ('D', 'C')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE idempotency_keys (
    idempotency_key VARCHAR(100) PRIMARY KEY,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ledger_transaction_id ON ledger_entries(transaction_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_ledger_transaction_id;
DROP TABLE IF EXISTS ledger_entries;
DROP TABLE IF EXISTS idempotency_keys;
DROP TABLE IF EXISTS accounts;

-- +goose StatementEnd
