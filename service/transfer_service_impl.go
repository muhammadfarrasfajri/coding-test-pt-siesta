package service

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	errs "github.com/muhammadfarrasfajri/coding-test-pt-siesta/pkg"
	"github.com/muhammadfarrasfajri/coding-test-pt-siesta/repository"
)

type transferService struct {
	db   *sql.DB
	repo repository.TransferRepository
}

func NewTransferService(db *sql.DB, repo repository.TransferRepository) *transferService {
	return &transferService{
		db:   db,
		repo: repo,
	}
}

func (s *transferService) ExecuteTransfer(ctx context.Context, traceID, idempotencyKey, fromAccount, toAccount string, amount float64) error {
	if err := ctx.Err(); err != nil {
		return errs.NewError("ERR_CONTEXT_CANCELED", traceID, errs.SeverityWarn, "Request cancelled", false, err)
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return errs.NewError("ERR_TX_START", traceID, errs.SeverityError, "Failed to start transaction", true, err)
	}
	defer tx.Rollback()

	if err := s.repo.InsertIdempotency(ctx, tx, idempotencyKey); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errs.NewError("ERR_IDEMPOTENT_CONFLICT", traceID, errs.SeverityInfo, "Transaction already processed", false, nil)
		}
		return errs.NewError("ERR_IDEMPOTENT_INSERT", traceID, errs.SeverityError, "DB error on idempotency", true, err)
	}

	accounts, err := s.repo.GetAccountsForUpdate(ctx, tx, []string{fromAccount, toAccount})
	if err != nil {
		return errs.NewError("ERR_LOCK_ACCOUNT", traceID, errs.SeverityError, "Failed to lock accounts", true, err)
	}
	if len(accounts) != 2 {
		return errs.NewError("ERR_INVALID_ACCOUNT", traceID, errs.SeverityWarn, "Accounts not found", false, nil)
	}

	var fromBalance float64
	for _, acc := range accounts {
		if acc.ID == fromAccount {
			fromBalance = acc.Balance
			break
		}
	}
	if fromBalance < amount {
		return errs.NewError("ERR_INSUFFICIENT_FUNDS", traceID, errs.SeverityWarn, "Insufficient balance", false, nil)
	}

	// Eksekusi Mutasi Saldo
	if err := s.repo.UpdateBalance(ctx, tx, fromAccount, -amount); err != nil { // Debit
		return errs.NewError("ERR_DEBIT_FAILED", traceID, errs.SeverityError, "Deduct balance failed", true, err)
	}
	if err := s.repo.UpdateBalance(ctx, tx, toAccount, amount); err != nil { // Credit
		return errs.NewError("ERR_CREDIT_FAILED", traceID, errs.SeverityError, "Add balance failed", true, err)
	}

	// Mencatat Ledger
	if err := s.repo.InsertLedgerEntries(ctx, tx, traceID, fromAccount, toAccount, amount); err != nil {
		return errs.NewError("ERR_LEDGER_INSERT", traceID, errs.SeverityError, "Failed to insert ledger", true, err)
	}

	if err := tx.Commit(); err != nil {
		return errs.NewError("ERR_COMMIT", traceID, errs.SeverityFatal, "Commit failed", true, err)
	}

	return nil
}
