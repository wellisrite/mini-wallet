package repository

import (
	"database/sql"
	"fmt"
	"julo/domain"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// TransactionsRepository represents the repository for wallet-related operations.
type TransactionsRepository struct {
	db *sql.DB
}

// NewTransactionsRepository creates a new TransactionsRepository with the given database connection.
func NewTransactionsRepository(db *sql.DB) ITransactionRepository {
	return &TransactionsRepository{db}
}

func (tr *TransactionsRepository) GetTransactionsByCustomerID(customerID string) ([]*domain.Transaction, error) {
	// Implement the retrieval of transactions by customer ID here.
	// You should write a SQL query or use an ORM to fetch the transactions from the database.
	// Return the retrieved transactions or an error if any.

	// Example of a SQL query (assuming you have a transactions table):
	query := "SELECT id, status, transacted_at, type, amount, reference_id FROM transactions WHERE customer_id = $1"

	rows, err := tr.db.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		err := rows.Scan(&transaction.ID, &transaction.Status, &transaction.TransactedAt, &transaction.Type, &transaction.Amount, &transaction.ReferenceID)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

// AddDeposit adds a deposit to the transaction table.
func (tr *TransactionsRepository) AddDeposit(customerID, referenceID string, amount int) (*domain.Deposit, error) {

	// Get the current timestamp for the deposit.
	depositTime := time.Now()

	tx, err := tr.db.Begin()
	if err != nil {
		// Handle the error
		return nil, err
	}

	// Generate a new unique ID for the deposit (you can use a UUID library).
	depositID, err := tr.GenerateUniqueUUID()
	if err != nil {
		return nil, err
	}

	// Insert the deposit into the transactions table.
	query := `
		INSERT INTO transactions (id, customer_id, reference_id, transacted_at, amount, type, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'success')
	`
	_, err = tx.Exec(
		query,
		depositID,
		customerID,
		referenceID,
		depositTime,
		amount,
		domain.TYPE_DEPOSIT,
	)
	if err != nil {
		tx.Rollback()
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "23505" { // 23505 is the PostgreSQL error code for unique_violation
			// Handle the unique constraint violation error here
			// For example, you can return a custom error or take appropriate action
			return nil, fmt.Errorf("Transaction with reference_id %s already exists", referenceID)
		}

		// Handle the error
		return nil, err
	}

	// Update the wallet balance
	updateBalanceQuery := `
	   UPDATE wallet
	   SET balance = balance + $1
	   WHERE owned_by = $2
   `
	_, err = tx.Exec(updateBalanceQuery, amount, customerID)
	if err != nil {
		tx.Rollback()

		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()

		return nil, err
	}

	// Create and return the deposit object.
	deposit := &domain.Deposit{
		ID:          depositID,
		DepositedBy: customerID,
		ReferenceID: referenceID,
		DepositedAt: depositTime,
		Amount:      amount,
		Status:      "success",
	}

	return deposit, nil
}

// AddWithdrawal inserts a withdrawal record into the database.
func (tr *TransactionsRepository) AddWithdrawal(withdrawal *domain.Withdrawal) error {
	tx, err := tr.db.Begin()
	if err != nil {
		tx.Rollback()
		// Handle the error
		return err
	}
	// Check if there are sufficient funds
	currentBalanceQuery := "SELECT balance FROM wallet WHERE owned_by = $1"
	var currentBalance int
	if err := tx.QueryRow(currentBalanceQuery, withdrawal.WithdrawnBy).Scan(&currentBalance); err != nil {
		return err
	}

	if currentBalance < withdrawal.Amount {
		return fmt.Errorf("insufficient fund") // Define this error accordingly
	}

	query := `
        INSERT INTO transactions (id, customer_id, status, transacted_at, type, amount, reference_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err = tx.Exec(
		query,
		withdrawal.ID,
		withdrawal.WithdrawnBy,
		withdrawal.Status,
		withdrawal.WithdrawnAt,
		domain.TYPE_WITHDRAWAL,
		withdrawal.Amount,
		withdrawal.ReferenceID,
	)

	if err != nil {
		tx.Rollback()
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "23505" { // 23505 is the PostgreSQL error code for unique_violation
			// Handle the unique constraint violation error here
			// For example, you can return a custom error or take appropriate action
			return fmt.Errorf("Transaction with reference_id %s already exists", withdrawal.ReferenceID)
		}

		return err
	}

	// Update the wallet balance
	updateBalanceQuery := `
	  UPDATE wallet
	  SET balance = balance - $1
	  WHERE owned_by = $2
  `
	_, err = tx.Exec(updateBalanceQuery, withdrawal.Amount, withdrawal.WithdrawnBy)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (tr *TransactionsRepository) GenerateUniqueUUID() (string, error) {
	var uuidStr string

	// Keep generating UUIDs until a unique one is found
	for {
		newUUID := uuid.New()
		uuidStr = newUUID.String()

		// Check if the generated UUID already exists in the database
		var count int
		err := tr.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE id = $1", uuidStr).Scan(&count)
		if err != nil {
			// Handle the database query error
			return "", err
		}

		// If the UUID is not found in the database, break the loop
		if count == 0 {
			break
		}
	}

	return uuidStr, nil
}
