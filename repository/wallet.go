// julo/repository/wallet_repository.go
package repository

import (
	"database/sql"

	"julo/domain"
)

// WalletRepository represents the repository for wallet-related operations.
type WalletRepository struct {
	db *sql.DB
}

// NewWalletRepository creates a new WalletRepository with the given database connection.
func NewWalletRepository(db *sql.DB) IWalletRepository {
	return &WalletRepository{db}
}

// InsertWallet inserts a wallet into the database.
func (wr *WalletRepository) InsertWallet(wallet *domain.Wallet) error {
	tx, err := wr.db.Begin()
	if err != nil {
		// Handle the error
		return err
	}

	// Insert the wallet data into the "wallets" table
	query := `
        INSERT INTO wallet (id, owned_by, status, enabled_at, balance)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err = tx.Exec(
		query,
		wallet.ID,
		wallet.OwnedBy,
		wallet.Status,
		wallet.EnabledAt,
		wallet.Balance,
	)
	if err != nil {
		tx.Rollback()
		// Handle the error
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// Add more repository functions as needed (e.g., UpdateWallet, GetWalletByID, etc.)
// Check if the customer ID already exists in the database
func (wr *WalletRepository) WalletExists(customerID string) (bool, error) {
	// Query the database to check if a wallet with the given customerID exists
	query := "SELECT COUNT(*) FROM wallet WHERE owned_by = $1"
	var count int
	err := wr.db.QueryRow(query, customerID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (wr *WalletRepository) GetWallet(customerID string) (*domain.Wallet, error) {
	// Query the database to retrieve the wallet data by customerID
	query := "SELECT id, owned_by, status, enabled_at, balance FROM wallet WHERE owned_by = $1"
	row := wr.db.QueryRow(query, customerID)

	var wallet domain.Wallet
	err := row.Scan(&wallet.ID, &wallet.OwnedBy, &wallet.Status, &wallet.EnabledAt, &wallet.Balance)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// UpdateWallet updates the wallet in the database
func (wr *WalletRepository) UpdateWallet(wallet *domain.Wallet) error {
	tx, err := wr.db.Begin()
	if err != nil {
		// Handle the error
		return err
	}

	query := `
        UPDATE wallet
        SET
            status = $2,
            enabled_at = $3,
            balance = $4
        WHERE id = $1
    `
	_, err = tx.Exec(
		query,
		wallet.ID,
		wallet.Status,
		wallet.EnabledAt,
		wallet.Balance,
	)
	if err != nil {
		tx.Rollback()
		// Handle the error
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// CalculateWalletBalance calculates the wallet balance by summing deposits and subtracting withdrawals.
func (wr *WalletRepository) CalculateWalletBalance(customerID string) (int, error) {
	tx, err := wr.db.Begin()
	if err != nil {
		// Handle the error
		return 0, err
	}

	depositQuery := "SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE customer_id = $1 AND type = 'deposit'"
	withdrawalQuery := "SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE customer_id = $1 AND type = 'withdrawal'"

	var depositTotal, withdrawalTotal int
	err = tx.QueryRow(depositQuery, customerID).Scan(&depositTotal)
	if err != nil {
		return 0, err
	}

	err = tx.QueryRow(withdrawalQuery, customerID).Scan(&withdrawalTotal)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	balance := depositTotal - withdrawalTotal
	return balance, nil
}
