package service

import "context"

type TransferService interface {
	ExecuteTransfer(ctx context.Context, traceID, idempotencyKey, fromAccount, toAccount string, amount float64) error
}
