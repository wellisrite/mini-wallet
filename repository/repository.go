package repository

import "julo/domain"

// WalletRepository defines the methods for wallet-related operations.
type IWalletRepository interface {
	InsertWallet(wallet *domain.Wallet) error
	WalletExists(customerID string) (bool, error)
	GetWallet(customerID string) (*domain.Wallet, error)
	UpdateWallet(wallet *domain.Wallet) error
	CalculateWalletBalance(customerID string) (int, error)
}

type ITransactionRepository interface {
	GetTransactionsByCustomerID(customerID string) ([]*domain.Transaction, error)
	AddDeposit(customerID, referenceID string, amount int) (*domain.Deposit, error)
	GenerateUniqueUUID() (string, error)
	AddWithdrawal(withdrawal *domain.Withdrawal) error
}
