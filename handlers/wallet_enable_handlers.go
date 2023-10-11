// julo/handlers/wallet_enable_handlers.go
package handlers

import (
	"net/http"

	"julo/domain"
	"julo/domain/dto"
)

// EnableWalletHandler handles enabling the wallet.
func (h *Handlers) EnableWalletHandler(w http.ResponseWriter, r *http.Request) {
	// Access the customerXID from the request context
	customerXID := r.Context().Value("customerXID").(string)

	// fetch wallet from DB
	wallet, err := h.WalletRepository.GetWallet(customerXID)
	if err != nil {
		http.Error(w, "Failed to get wallet", http.StatusInternalServerError)
		return
	}

	if wallet.Status == domain.STATUS_ENABLED {
		h.ErrorResponse(w, http.StatusBadRequest, "Already enabled")
		return
	}

	domain.EnableWallet(wallet)

	// Insert the wallet into the database
	if err := h.WalletRepository.UpdateWallet(wallet); err != nil {
		http.Error(w, "Failed to insert wallet into the database", http.StatusInternalServerError)
		return
	}

	// Create a response
	response := dto.EnableWalletResponse{
		Status: "success",
		Data:   &dto.WalletResponseData{Wallet: wallet},
	}

	h.respondWithJSON(w, http.StatusCreated, response)
}
