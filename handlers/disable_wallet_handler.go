// julo/handlers/wallet_disable_handlers.go
package handlers

import (
	"net/http"

	"julo/domain"
	"julo/domain/dto"
)

// DisableWalletHandler handles disabling the wallet.
func (h *Handlers) DisableWalletHandler(w http.ResponseWriter, r *http.Request) {
	// Access the customerXID from the request context
	customerXID := r.Context().Value("customerXID").(string)

	// Fetch wallet from DB
	wallet, err := h.WalletRepository.GetWallet(customerXID)
	if err != nil {
		http.Error(w, "Failed to get wallet", http.StatusInternalServerError)
		return
	}

	if wallet.Status == domain.STATUS_DISABLED {
		h.ErrorResponse(w, http.StatusBadRequest, "Wallet is already disabled")
		return
	}

	// Disable the wallet
	wallet.Status = domain.STATUS_DISABLED

	// Update the wallet in the database
	if err := h.WalletRepository.UpdateWallet(wallet); err != nil {
		http.Error(w, "Failed to update wallet in the database", http.StatusInternalServerError)
		return
	}

	// Create a response
	response := dto.ViewWalletBalanceResponse{
		Status: "success",
		Data:   &dto.WalletResponseData{Wallet: wallet},
	}

	h.respondWithJSON(w, http.StatusOK, response)
}
