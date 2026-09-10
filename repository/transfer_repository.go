package repository

import (
	"context"
	"database/sql"

	"github.com/muhammadfarrasfajri/coding-test-pt-siesta/model"
)

type TransferRepository interface {
	InsertIdempotency(ctx context.Context, tx *sql.Tx, key string) error
	GetAccountsForUpdate(ctx context.Context, tx *sql.Tx, accountIDs []string) ([]model.Account, error)
	UpdateBalance(ctx context.Context, tx *sql.Tx, accountID string, amount float64) error
	InsertLedgerEntries(ctx context.Context, tx *sql.Tx, txID, fromID, toID string, amount float64) error
}
