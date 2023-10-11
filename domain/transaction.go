package domain

import (
	"time"
)

const (
	TYPE_DEPOSIT    = "deposit"
	TYPE_WITHDRAWAL = "withdrawal"
)

// Transaction represents a financial transaction.
type Transaction struct {
	ID           string    `json:"id"`
	CustomerID   string    `json:"customer_id"`
	Amount       int       `json:"amount"`
	TransactedAt time.Time `json:"timestamp"`
	Status       string    `json:"status"`
	Type         string    `json:"type"`
	ReferenceID  string    `json:"reference_id"`
	// Other fields specific to your transaction data.
}

// Deposit represents a deposit made to a wallet.
type Deposit struct {
	ID          string    `json:"id"`
	DepositedBy string    `json:"deposited_by"`
	Status      string    `json:"status"`
	DepositedAt time.Time `json:"deposited_at"`
	Amount      int       `json:"amount"`
	ReferenceID string    `json:"reference_id"`
}

type Withdrawal struct {
	ID          string    `json:"id"`
	WithdrawnBy string    `json:"withdrawn_by"`
	Status      string    `json:"status"`
	WithdrawnAt time.Time `json:"withdrawn_at"`
	Amount      int       `json:"amount"`
	ReferenceID string    `json:"reference_id"`
}
