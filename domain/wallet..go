// julo/domain/wallet_domain.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	STATUS_ENABLED  = "enabled"
	STATUS_DISABLED = "disabled"
)

// Wallet represents a user's wallet.
type Wallet struct {
	ID        string    `json:"id"`
	OwnedBy   string    `json:"owned_by"`
	Status    string    `json:"status"`
	EnabledAt time.Time `json:"enabled_at"`
	Balance   int       `json:"balance"`
}

// InitializeWallet initializes a new wallet and returns a Wallet instance.
func InitializeWallet(customerXID string) *Wallet {
	// Generate a unique wallet ID (for demonstration purposes, we use a random UUID here)
	walletID := uuid.New().String()

	// Create a new wallet with default values
	wallet := Wallet{
		ID:      walletID,
		OwnedBy: customerXID,
		Status:  STATUS_DISABLED, // Default status is disabled
		Balance: 0,               // Initial balance is 0
	}

	return &wallet
}

// EnableWallet enables the wallet and updates the status and enabled timestamp.
func EnableWallet(wallet *Wallet) {
	wallet.Status = STATUS_ENABLED
	wallet.EnabledAt = time.Now()
}
