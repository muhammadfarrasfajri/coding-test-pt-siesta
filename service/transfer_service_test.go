package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/muhammadfarrasfajri/coding-test-pt-siesta/model"
	errs "github.com/muhammadfarrasfajri/coding-test-pt-siesta/pkg"
)

type MockTransferRepo struct {
	mock.Mock
}

func (m *MockTransferRepo) InsertIdempotency(ctx context.Context, tx *sql.Tx, key string) error {
	args := m.Called(ctx, tx, key)
	return args.Error(0)
}

func (m *MockTransferRepo) GetAccountsForUpdate(ctx context.Context, tx *sql.Tx, accountIDs []string) ([]model.Account, error) {
	args := m.Called(ctx, tx, accountIDs)
	return args.Get(0).([]model.Account), args.Error(1)
}

func (m *MockTransferRepo) UpdateBalance(ctx context.Context, tx *sql.Tx, accountID string, amount float64) error {
	args := m.Called(ctx, tx, accountID, amount)
	return args.Error(0)
}

func (m *MockTransferRepo) InsertLedgerEntries(ctx context.Context, tx *sql.Tx, txID, fromID, toID string, amount float64) error {
	args := m.Called(ctx, tx, txID, fromID, toID, amount)
	return args.Error(0)
}

// Fungsi Test Utama (Table-Driven Test)
func TestExecuteTransfer(t *testing.T) {
	// Setup SQL Mock
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Gagal membuat sqlmock: %v", err)
	}
	defer db.Close()

	// Skenario Test
	tests := []struct {
		name            string
		amount          float64
		mockSetup       func(mockRepo *MockTransferRepo)
		expectedErrCode string
	}{
		{
			name:   "Skenario 1: Sukses Melakukan Transfer",
			amount: 50000,
			mockSetup: func(mockRepo *MockTransferRepo) {
				sqlMock.ExpectBegin() // Ekspektasi BeginTx() dipanggil

				mockRepo.On("InsertIdempotency", mock.Anything, mock.Anything, "key-123").Return(nil)

				// Simulasi Saldo Pengirim cukup (100rb)
				mockRepo.On("GetAccountsForUpdate", mock.Anything, mock.Anything, []string{"A-1", "B-2"}).
					Return([]model.Account{{ID: "A-1", Balance: 100000}, {ID: "B-2", Balance: 0}}, nil)

				mockRepo.On("UpdateBalance", mock.Anything, mock.Anything, "A-1", float64(-50000)).Return(nil)
				mockRepo.On("UpdateBalance", mock.Anything, mock.Anything, "B-2", float64(50000)).Return(nil)
				mockRepo.On("InsertLedgerEntries", mock.Anything, mock.Anything, "trace-1", "A-1", "B-2", float64(50000)).Return(nil)

				sqlMock.ExpectCommit() // Ekspektasi Commit() dipanggil
			},
			expectedErrCode: "", // Tidak ada error
		},
		{
			name:   "Skenario 2: Gagal karena Saldo Tidak Cukup (Insufficient Funds)",
			amount: 150000,
			mockSetup: func(mockRepo *MockTransferRepo) {
				sqlMock.ExpectBegin()

				mockRepo.On("InsertIdempotency", mock.Anything, mock.Anything, "key-456").Return(nil)

				// Simulasi Saldo Pengirim HANYA 100rb, tapi mau transfer 150rb
				mockRepo.On("GetAccountsForUpdate", mock.Anything, mock.Anything, []string{"A-1", "B-2"}).
					Return([]model.Account{{ID: "A-1", Balance: 100000}, {ID: "B-2", Balance: 0}}, nil)

				// Perhatikan: Karena saldo tidak cukup, UpdateBalance tidak akan dipanggil
				// dan Tx.Rollback() akan otomatis tereksekusi lewat defer
				sqlMock.ExpectRollback()
			},
			expectedErrCode: "ERR_INSUFFICIENT_FUNDS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTransferRepo)
			tt.mockSetup(mockRepo)

			// Inisialisasi Service dengan DB Mock dan Repo Mock
			svc := NewTransferService(db, mockRepo)

			// Identitas dummy berdasarkan test case
			idempotencyKey := "key-123"
			if tt.expectedErrCode == "ERR_INSUFFICIENT_FUNDS" {
				idempotencyKey = "key-456"
			}

			// Eksekusi fungsi
			err := svc.ExecuteTransfer(context.Background(), "trace-1", idempotencyKey, "A-1", "B-2", tt.amount)

			// Asersi (Pengecekan Hasil)
			if tt.expectedErrCode == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				coreErr, ok := err.(*errs.CoreBankingError)
				assert.True(t, ok, "Error harus berupa tipe CoreBankingError")
				assert.Equal(t, tt.expectedErrCode, coreErr.Code)
			}

			// Pastikan semua method di repo terpanggil sesuai ekspektasi
			mockRepo.AssertExpectations(t)
			assert.NoError(t, sqlMock.ExpectationsWereMet())
		})
	}
}
