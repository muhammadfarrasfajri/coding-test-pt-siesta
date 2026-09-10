package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/muhammadfarrasfajri/coding-test-pt-siesta/model"
)

type transferRepository struct {
}

func NewTransferRepository() *transferRepository {
	return &transferRepository{}
}

func (r *transferRepository) InsertIdempotency(ctx context.Context, tx *sql.Tx, key string) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO idempotency_keys (idempotency_key, status) VALUES ($1, 'COMPLETED')`, key)
	return err
}

func (r *transferRepository) GetAccountsForUpdate(ctx context.Context, tx *sql.Tx, accountIDs []string) ([]model.Account, error) {

	rows, err := tx.QueryContext(ctx, `
		SELECT id, balance FROM accounts 
		WHERE id = ANY($1) 
		ORDER BY id FOR UPDATE
	`, pq.Array(accountIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []model.Account
	for rows.Next() {
		var acc model.Account
		if err := rows.Scan(&acc.ID, &acc.Balance); err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (r *transferRepository) UpdateBalance(ctx context.Context, tx *sql.Tx, accountID string, amount float64) error {
	_, err := tx.ExecContext(ctx, `UPDATE accounts SET balance = balance + $1, updated_at = $2 WHERE id = $3`, amount, time.Now(), accountID)
	return err
}

func (r *transferRepository) InsertLedgerEntries(ctx context.Context, tx *sql.Tx, txID, fromID, toID string, amount float64) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO ledger_entries (transaction_id, account_id, amount, dc_indicator) 
		VALUES ($1, $2, $3, 'D'), ($1, $4, $5, 'C')
	`, txID, fromID, amount, toID, amount)
	return err
}
