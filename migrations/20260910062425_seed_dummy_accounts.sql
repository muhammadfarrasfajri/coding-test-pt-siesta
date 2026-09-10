-- +goose Up
-- +goose StatementBegin

INSERT INTO accounts (id, balance) VALUES 
('A001', 100000), 
('B002', 50000)
ON CONFLICT (id) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM accounts WHERE id IN ('A001', 'B002');

-- +goose StatementEnd